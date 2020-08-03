// Copyright 1999-2020. Plesk International GmbH.

package cmd

import (
	"github.com/plesk/pleskapp/plesk/actions"
	"github.com/plesk/pleskapp/plesk/config"
	"github.com/plesk/pleskapp/plesk/locales"
	"github.com/plesk/pleskapp/plesk/utils"
	"github.com/spf13/cobra"
)

var loginCmd = &cobra.Command{
	Use:   locales.L.Get("server.login.cmd"),
	Short: locales.L.Get("server.login.description"),
	RunE: func(cmd *cobra.Command, args []string) error {
		server, err := config.GetServer(args[0])
		if err != nil {
			return err
		}

		cmd.SilenceUsage = true
		err = actions.ServerLogin(*server)

		if err != nil {
			utils.Log.Print(locales.L.Get("errors.execution.failed.generic", err.Error()))
		}
		return err
	},
	Args: cobra.ExactArgs(1),
}

func init() {
	ServersCmd.AddCommand(loginCmd)
}
