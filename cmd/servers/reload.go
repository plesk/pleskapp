// Copyright 1999-2021. Plesk International GmbH.

package cmd

import (
	"github.com/plesk/pleskapp/plesk/actions"
	"github.com/plesk/pleskapp/plesk/locales"
	"github.com/plesk/pleskapp/plesk/utils"
	"github.com/spf13/cobra"
)

var reloadCmd = &cobra.Command{
	Use:   "reload [IP ADDRESS|HOSTNAME]",
	Short: locales.L.Get("server.reload.description"),
	RunE: func(cmd *cobra.Command, args []string) error {
		server, err := getServer(args)
		if err != nil {
			return err
		}

		cmd.SilenceUsage = true
		err = actions.ServerUpdate(*server)

		if err == nil {
			utils.Log.PrintL("server.reload.success")
		}

		return err
	},
	Args: cobra.MaximumNArgs(1),
}

func init() {
	ServersCmd.AddCommand(reloadCmd)
}
