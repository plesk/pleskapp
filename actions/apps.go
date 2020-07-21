// Copyright 1999-2020. Plesk International GmbH. All rights reserved.

package actions

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"

	"git.plesk.ru/~abashurov/pleskapp/api/factory"
	"git.plesk.ru/~abashurov/pleskapp/locales"
	"git.plesk.ru/~abashurov/pleskapp/types"
	"git.plesk.ru/~abashurov/pleskapp/utils"
)

func AppAdd(features []string, subdir string, path string, overwrite bool) error {
	d, err := os.Stat(path)
	if err != nil {
		return err
	}

	if !d.IsDir() {
		return errors.New(locales.L.Get("errors.path.is.not.directory", path))
	}

	if !overwrite {
		_, err := os.Stat(path + "/.pleskapp")
		if err == nil {
			return errors.New(locales.L.Get("errors.path.already.exists", path+"/.pleskapp"))
		}
	}

	a := types.App{
		TargetPath: subdir,
		Features:   features,
	}
	aj, err := json.Marshal(a)
	if err != nil {
		return err
	}

	ioutil.WriteFile(path+"/.pleskapp", aj, 0600)

	return nil
}

func AppDeploy(host types.Server, path string, domain types.Domain) error {
	d, err := os.Stat(path)
	if err != nil {
		return err
	}

	if !d.IsDir() {
		return errors.New(locales.L.Get("errors.path.is.not.directory", path))
	}

	var c types.App
	f, err := ioutil.ReadFile(path + "/.pleskapp")
	if err == nil {
		err = json.Unmarshal(f, &c)
		if err != nil {
			utils.Log.Error(locales.L.Get("errors.cannot.parse.config", path+"/.pleskapp"))
			c = types.App{}
		}
	}

	api := factory.GetDomainManagement(host.GetServerAuth())
	err = api.AddDomainFeatures(domain.Name, c.Features)
	if err != nil {
		return err
	}

	if c.TargetPath != "" {
		return UploadDirectory(host, domain, true, false, path, &c.TargetPath)
	}

	return UploadDirectory(host, domain, true, false, path, nil)
}
