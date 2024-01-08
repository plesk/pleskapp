// Copyright 1999-2024. Plesk International GmbH.

package actions

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/plesk/pleskapp/plesk/internal/api/factory"
	"github.com/plesk/pleskapp/plesk/internal/locales"
	"github.com/plesk/pleskapp/plesk/internal/types"
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
		_, err := os.Stat(path + "/.plesk")
		if err == nil {
			return errors.New(locales.L.Get("errors.path.already.exists", path+"/.plesk"))
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

	os.WriteFile(path+"/.plesk", aj, 0600)

	return nil
}

func AppDeploy(host types.Server, app types.App, path string, domain types.Domain) error {
	api := factory.GetDomainManagement(host.GetServerAuth())
	err := api.AddDomainFeatures(domain.Name, app.Features, host.Info.IsWindows)
	if err != nil {
		return err
	}

	if app.TargetPath != "" {
		return UploadDirectory(host, domain, true, false, path, &app.TargetPath)
	}

	return UploadDirectory(host, domain, true, false, path, nil)
}
