// Copyright 1999-2020. Plesk International GmbH.

package cmd

import (
	"github.com/plesk/pleskapp/actions"
	"github.com/plesk/pleskapp/config"
	"github.com/plesk/pleskapp/locales"
	"github.com/plesk/pleskapp/utils"
	"github.com/spf13/cobra"
)

var reloadCmd = &cobra.Command{
	Use:   locales.L.Get("server.reload.cmd"),
	Short: locales.L.Get("server.reload.description"),
	RunE: func(cmd *cobra.Command, args []string) error {
		server, err := config.GetServer(args[0])
		if err != nil {
			return err
		}

		cmd.SilenceUsage = true
		return utils.Log.PrintSuccessOrError("server.reload.success", nil, actions.ServerUpdate(*server))
	},
	Args: cobra.ExactArgs(1),
}

func init() {
	ServersCmd.AddCommand(reloadCmd)
}
