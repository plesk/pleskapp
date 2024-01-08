// Copyright 1999-2024. Plesk International GmbH.

package actions

import (
	"fmt"
	"github.com/plesk/pleskapp/plesk/internal/locales"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/plesk/pleskapp/plesk/internal/api/factory"
	"github.com/plesk/pleskapp/plesk/internal/types"
	"github.com/plesk/pleskapp/plesk/internal/upload"
	"github.com/plesk/pleskapp/plesk/internal/utils"
)

func getPrereq(host types.Server, domain types.Domain) (*types.FtpUser, *string, *string, error) {
	var fullPath string
	var docRoot string
	// TODO: Ideally, there should be no need to do this
	{
		api := factory.GetDomainManagement(host.GetServerAuth())
		i, err := api.GetDomain(domain.Name)
		if err != nil {
			return nil, nil, nil, err
		}

		fullPath = i.WWWRoot
		parts := strings.Split(i.WWWRoot, domain.Name)
		docRoot = parts[len(parts)-1]
	}

	ftp := FindCachedFtpUser(domain)
	if ftp == nil {
		var err error
		ftp, err = FtpUserCreate(host, domain, nil)
		if err != nil {
			return nil, nil, nil, err
		}
	}

	// assume the parent directory as an app root directory
	appRoot := filepath.Dir(docRoot)

	return ftp, &appRoot, &fullPath, nil
}

func UploadFileToRoot(host types.Server, domain types.Domain, ovw bool, file string) (string, error) {
	ftp, docroot, path, err := getPrereq(host, domain)
	if err != nil {
		return "", err
	}

	connection, err := upload.Connect(*ftp, host.Host, "/")
	if err != nil {
		return "", err
	}

	cr, f := utils.GetClientRootName(file)
	return strings.Split(*path, *docroot)[0], connection.UploadFile(cr, "/", f, ovw, host.Info.IsWindows)
}

func UploadFile(host types.Server, domain types.Domain, ovw bool, file string) error {
	ftp, docroot, _, err := getPrereq(host, domain)
	if err != nil {
		return err
	}

	connection, err := upload.Connect(*ftp, host.Host, *docroot)
	if err != nil {
		return err
	}

	cr, f := utils.GetClientRootName(file)
	return connection.UploadFile(cr, *docroot, f, ovw, host.Info.IsWindows)
}

func UploadDirectory(host types.Server, domain types.Domain, ovw bool, dry bool, dir string, root *string) error {
	ftp, docroot, _, err := getPrereq(host, domain)
	if err != nil {
		return err
	}

	connection, err := upload.Connect(*ftp, host.Host, *docroot)
	if err != nil {
		return err
	}

	serverPath := root
	if serverPath == nil {
		serverPath = docroot
	}

	return filepath.Walk(dir, func(p string, info os.FileInfo, err error) error {
		if info.IsDir() && strings.HasPrefix(info.Name(), ".") {
			fmt.Println("Skip hidden directory:", info.Name())
			return filepath.SkipDir
		}

		pathPart := strings.TrimPrefix(p, dir)
		if pathPart == "." {
			return nil
		}

		if dry {
			fmt.Println(locales.L.Get("upload.dry.run.upload", path.Join(dir, pathPart), path.Join(*serverPath, pathPart)))
			return nil
		}

		return connection.UploadFile(dir+"/", *serverPath+"/", pathPart, ovw, host.Info.IsWindows)
	})
}
