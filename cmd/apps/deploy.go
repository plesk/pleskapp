// Copyright 1999-2020. Plesk International GmbH.

package cmd

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	"git.plesk.ru/projects/SBX/repos/pleskapp/actions"
	"git.plesk.ru/projects/SBX/repos/pleskapp/config"
	"git.plesk.ru/projects/SBX/repos/pleskapp/locales"
	"git.plesk.ru/projects/SBX/repos/pleskapp/types"
	"git.plesk.ru/projects/SBX/repos/pleskapp/utils"
	"github.com/spf13/cobra"
)

var deployCmd = &cobra.Command{
	Use:   locales.L.Get("app.deploy.cmd"),
	Short: locales.L.Get("app.deploy.description"),
	RunE: func(cmd *cobra.Command, args []string) error {
		path, err := filepath.Abs(".")
		if err != nil {
			return err
		}

		if len(args) > 0 {
			path, err = filepath.Abs(args[0])
			if err != nil {
				return err
			}
		}

		d, err := os.Stat(path)
		if err != nil {
			return err
		}

		if !d.IsDir() {
			return errors.New(locales.L.Get("errors.path.is.not.directory", path))
		}

		var c types.App
		f, err := ioutil.ReadFile(path + "/.pleskapp")
		if err != nil {
			return errors.New(locales.L.Get("errors.cannot.parse.config", path+"/.pleskapp"))
		}

		err = json.Unmarshal(f, &c)
		if err != nil {
			return errors.New(locales.L.Get("errors.cannot.parse.config", path+"/.pleskapp"))
		}

		server, err := config.GetServer(c.Server)
		if err != nil {
			return err
		}

		domain, err := config.GetDomain(*server, c.Domain)
		if err != nil {
			return err
		}

		cmd.SilenceUsage = true
		return utils.Log.PrintSuccessOrError("app.deploy.success", nil, actions.AppDeploy(*server, c, path, *domain))
	},
	Args: cobra.MaximumNArgs(1),
}

func init() {
	AppsCmd.AddCommand(deployCmd)
}
