// Copyright 1999-2020. Plesk International GmbH.

package cmd

import (
	"git.plesk.ru/projects/SBX/repos/pleskapp/actions"
	"git.plesk.ru/projects/SBX/repos/pleskapp/config"
	"git.plesk.ru/projects/SBX/repos/pleskapp/locales"
	"git.plesk.ru/projects/SBX/repos/pleskapp/utils"
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
				utils.Log.Print(locales.L.Get("server.delete.success"))
			}
		}
	},
	Args: cobra.MinimumNArgs(1),
}

func init() {
	ServersCmd.AddCommand(deleteCmd)
}
