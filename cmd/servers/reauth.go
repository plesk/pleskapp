// Copyright 1999-2021. Plesk International GmbH.

package cmd

import (
	"github.com/plesk/pleskapp/plesk/actions"
	"github.com/plesk/pleskapp/plesk/config"
	"github.com/plesk/pleskapp/plesk/locales"
	"github.com/plesk/pleskapp/plesk/utils"
	"github.com/spf13/cobra"
)

var reauthCmd = &cobra.Command{
	Use:   "reauth [IP ADDRESS|HOSTNAME]",
	Short: locales.L.Get("server.reauth.description"),
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true

		server, err := config.GetServerByArgs(args)
		if err != nil {
			return err
		}

		err = actions.ServerReauth(*server)

		if err == nil {
			utils.Log.PrintL("server.reauth.success")
		}

		return err
	},
	Args: cobra.MaximumNArgs(1),
}

func init() {
	ServersCmd.AddCommand(reauthCmd)
}
