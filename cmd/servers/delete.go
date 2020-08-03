// Copyright 1999-2020. Plesk International GmbH.

package cmd

import (
	"github.com/plesk/pleskapp/plesk/actions"
	"github.com/plesk/pleskapp/plesk/config"
	"github.com/plesk/pleskapp/plesk/locales"
	"github.com/plesk/pleskapp/plesk/utils"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   locales.L.Get("server.delete.cmd"),
	Short: locales.L.Get("server.delete.description"),
	Run: func(cmd *cobra.Command, args []string) {
		for _, host := range args {
			server, err := config.GetServer(host)
			if err != nil {
				utils.Log.Error(err.Error())
				continue
			}

			err = actions.ServerRemove(*server)
			if err != nil {
				utils.Log.Error(locales.L.Get("errors.server.remove.failure", host, err.Error()))
			} else {
				utils.Log.PrintL("server.delete.success")
			}
		}
	},
	Args: cobra.MinimumNArgs(1),
}

func init() {
	ServersCmd.AddCommand(deleteCmd)
}
