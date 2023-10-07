// Copyright 1999-2023. Plesk International GmbH.

package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/plesk/pleskapp/plesk/internal/actions"
	"github.com/plesk/pleskapp/plesk/internal/config"
	"github.com/plesk/pleskapp/plesk/internal/locales"
	"github.com/plesk/pleskapp/plesk/internal/types"
	"github.com/spf13/cobra"
)

var deployCmd = &cobra.Command{
	Use:   "deploy [PATH]",
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
		f, err := os.ReadFile(path + "/.plesk")
		if err != nil {
			return errors.New(locales.L.Get("errors.cannot.parse.config", path+"/.plesk"))
		}

		err = json.Unmarshal(f, &c)
		if err != nil {
			return errors.New(locales.L.Get("errors.cannot.parse.config", path+"/.plesk"))
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
		err = actions.AppDeploy(*server, c, path, *domain)
		if err == nil {
			fmt.Println(locales.L.Get("app.deploy.success", path))
		}

		return err
	},
	Args: cobra.MaximumNArgs(1),
}

func init() {
	AppsCmd.AddCommand(deployCmd)
}
