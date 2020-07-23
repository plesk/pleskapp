// Copyright 1999-2020. Plesk International GmbH. All rights reserved.

package cmd

import (
	"path/filepath"

	"git.plesk.ru/projects/SBX/repos/pleskapp/actions"
	"git.plesk.ru/projects/SBX/repos/pleskapp/config"
	"git.plesk.ru/projects/SBX/repos/pleskapp/locales"
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

		if len(args) > 2 {
			path, err = filepath.Abs(args[2])
			if err != nil {
				return err
			}
		}

		server, err := config.GetServer(args[0])
		if err != nil {
			return err
		}

		domain, err := config.GetDomain(*server, args[1])
		if err != nil {
			return err
		}

		cmd.SilenceUsage = true
		return utils.Log.PrintSuccessOrError("app.deploy.success", nil, actions.AppDeploy(*server, path, *domain))
	},
	Args: cobra.MinimumNArgs(2),
}

func init() {
	AppsCmd.AddCommand(deployCmd)
}
