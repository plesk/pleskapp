// Copyright 1999-2021. Plesk International GmbH.

package upload

import (
	"crypto/tls"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/jlaffaye/ftp"
	"github.com/plesk/pleskapp/plesk/internal/locales"
	"github.com/plesk/pleskapp/plesk/internal/types"
	"github.com/plesk/pleskapp/plesk/internal/utils"
)

type ftpConnection struct {
	auth    types.FtpUser
	inner   *ftp.ServerConn
	active  sync.Mutex
	docRoot string
	cwd     string
}

func Connect(creds types.FtpUser, server string, docRoot string) (*ftpConnection, error) {
	var connection ftpConnection
	var config = &tls.Config{
		InsecureSkipVerify:    true,
		VerifyPeerCertificate: nil,
	}

	inner, err := ftp.Dial(server+":21", ftp.DialWithTimeout(15*time.Second), ftp.DialWithExplicitTLS(config))
	if err != nil {
		return nil, err
	}

	err = inner.Login(creds.Login, creds.Password)
	if err != nil {
		return nil, err
	}
	cwd, err := inner.CurrentDir()
	if err != nil {
		return nil, err
	}

	connection = ftpConnection{
		auth:    creds,
		inner:   inner,
		active:  sync.Mutex{},
		docRoot: docRoot,
		cwd:     cwd,
	}

	// Keepalive
	go func(connection *ftpConnection) {
		for {
			if connection.inner == nil {
				return
			}

			connection.active.Lock()
			_ = connection.inner.NoOp()
			connection.active.Unlock()

			time.Sleep(15 * time.Second)
		}
	}(&connection)

	return &connection, nil
}

func (c *ftpConnection) Cwd(targetPath string) error {
	if c.cwd != targetPath {
		utils.Log.Debug(locales.L.Get("debug.cwd", targetPath))
		err := c.inner.ChangeDir(targetPath)
		if err != nil {
			return err
		}

		c.cwd = targetPath
	}
	return nil
}

func (c *ftpConnection) findFile(path string, fileName string) (*ftp.Entry, error) {
	stat, err := c.inner.List(path)
	if len(stat) > 0 && err == nil {
		for _, s := range stat {
			if s.Name != fileName {
				continue
			}

			return s, nil
		}
	}
	return nil, err
}

func (c *ftpConnection) UploadFile(
	clientRoot string,
	serverRoot string,
	fileName string,
	overwrite bool,
	isWindows bool,
) error {
	c.active.Lock()
	defer c.active.Unlock()

	var baseName string
	var basePath string

	split := utils.StrSplitRN(fileName, "/", 2)
	if len(split) > 1 {
		basePath = strings.ReplaceAll("/"+split[0]+"/", "//", "/")
		baseName = split[1]
	} else {
		basePath = "/"
		baseName = split[0]
	}

	if baseName == "" {
		// It's root, we skip it
		return nil
	}

	err := c.Cwd(serverRoot + basePath)
	if err != nil {
		return err
	}

	file, err := os.Stat(clientRoot + basePath + baseName)
	if err != nil {
		return err
	}

	entry, err := c.findFile(serverRoot+basePath, baseName)
	if entry != nil && err == nil {
		if file.IsDir() && entry.Type == ftp.EntryTypeFolder {
			utils.Log.Debug(locales.L.Get("debug.dir.skip", fileName))
			return nil
		}

		if !file.IsDir() && entry.Type == ftp.EntryTypeFile && !overwrite {
			utils.Log.Debug(locales.L.Get("debug.file.skip", fileName))
			return nil
		}
	}

	if file.IsDir() {
		err = c.inner.MakeDir(baseName)

		if err != nil {
			utils.Log.Error(locales.L.Get("errors.mkdir.failed", clientRoot+baseName, err.Error()))
		} else {
			utils.Log.Debug(locales.L.Get("debug.mkdir.success", clientRoot+baseName))
		}
	} else {
		file, err := os.Open(clientRoot + basePath + baseName)
		if err != nil {
			return err
		}
		err = c.inner.Stor(baseName, file)

		if err != nil {
			utils.Log.Error(locales.L.Get("errors.stor.failed", clientRoot+baseName, err.Error()))
		} else {
			utils.Log.Print(locales.L.Get("debug.stor.success", clientRoot+baseName))
		}
	}

	return nil
}

func (c *ftpConnection) Disconnect() {
	c.active.Lock()
	_ = c.inner.Quit()
	c.inner = nil

	c.active.Unlock()
}
