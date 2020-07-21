// Copyright 1999-2020. Plesk International GmbH. All rights reserved.

package actions

import (
	"os"
	"path/filepath"
	"strings"

	"git.plesk.ru/~abashurov/pleskapp/api/factory"
	"git.plesk.ru/~abashurov/pleskapp/locales"
	"git.plesk.ru/~abashurov/pleskapp/types"
	"git.plesk.ru/~abashurov/pleskapp/upload"
	"git.plesk.ru/~abashurov/pleskapp/utils"
)

func getPrereq(host types.Server, domain types.Domain) (*types.FtpUser, *string, *string, error) {
	var fullpath string
	var docroot string
	// TODO: Ideally, there should be no need to do this
	{
		api := factory.GetDomainManagement(host.GetServerAuth())
		i, err := api.GetDomain(domain.Name)
		if err != nil {
			return nil, nil, nil, err
		}

		fullpath = i.WWWRoot
		parts := strings.Split(i.WWWRoot, domain.Name)
		docroot = parts[len(parts)-1]
	}

	ftp := FindCachedFtpUser(domain)
	if ftp == nil {
		var err error
		ftp, err = FtpUserCreate(host, domain, nil)
		if err != nil {
			return nil, nil, nil, err
		}
	}

	return ftp, &docroot, &fullpath, nil
}

func UploadFileToRoot(host types.Server, domain types.Domain, ovw bool, file string) (string, error) {
	ftp, docroot, path, err := getPrereq(host, domain)
	if err != nil {
		return "", err
	}

	connection, err := upload.Connect(*ftp, domain.Name, "/")
	if err != nil {
		return "", err
	}

	cr, f := utils.GetClientRootName(file)
	return strings.Split(*path, *docroot)[0], connection.UploadFile(cr, "/", f, true)
}

func UploadFile(host types.Server, domain types.Domain, ovw bool, file string) error {
	ftp, docroot, _, err := getPrereq(host, domain)
	if err != nil {
		return err
	}

	connection, err := upload.Connect(*ftp, domain.Name, *docroot)
	if err != nil {
		return err
	}

	cr, f := utils.GetClientRootName(file)
	return connection.UploadFile(cr, *docroot, f, ovw)
}

func UploadDirectory(host types.Server, domain types.Domain, ovw bool, dry bool, dir string, root *string) error {
	ftp, docroot, _, err := getPrereq(host, domain)
	if err != nil {
		return err
	}

	connection, err := upload.Connect(*ftp, domain.Name, *docroot)
	if err != nil {
		return err
	}

	serverPath := root
	if serverPath == nil {
		serverPath = docroot
	}

	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		pathPart := strings.TrimPrefix(path, dir)
		if pathPart == "." {
			return nil
		}

		if dry {
			utils.Log.Print(locales.L.Get("upload.dry.run.upload", dir+"/"+pathPart, *serverPath+"/"+pathPart))
			return nil
		}

		return connection.UploadFile(dir+"/", *serverPath+"/", pathPart, ovw)
	})
}
