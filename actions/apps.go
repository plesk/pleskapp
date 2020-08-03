// Copyright 1999-2020. Plesk International GmbH.

package actions

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"

	"github.com/plesk/pleskapp/api/factory"
	"github.com/plesk/pleskapp/locales"
	"github.com/plesk/pleskapp/types"
)

func AppAdd(
	host types.Server,
	domain types.Domain,
	features []string,
	subdir string,
	path string,
	overwrite bool,
) error {
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
		Server:     host.Host,
		Domain:     domain.Name,
	}
	aj, err := json.Marshal(a)
	if err != nil {
		return err
	}

	ioutil.WriteFile(path+"/.pleskapp", aj, 0600)

	return nil
}

func AppDeploy(host types.Server, app types.App, path string, domain types.Domain) error {
	api := factory.GetDomainManagement(host.GetServerAuth())
	err := api.AddDomainFeatures(domain.Name, app.Features)
	if err != nil {
		return err
	}

	if app.TargetPath != "" {
		return UploadDirectory(host, domain, true, false, path, &app.TargetPath)
	}

	return UploadDirectory(host, domain, true, false, path, nil)
}
