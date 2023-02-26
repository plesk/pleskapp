// Copyright 1999-2023. Plesk International GmbH.

package cmd

import (
	"fmt"
	"github.com/plesk/pleskapp/plesk/internal/actions"
	"github.com/plesk/pleskapp/plesk/internal/config"
	"github.com/plesk/pleskapp/plesk/internal/locales"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete [IP ADDRESS|HOSTNAME ...]",
	Short: locales.L.Get("server.delete.description"),
	Run: func(cmd *cobra.Command, args []string) {
		for _, host := range args {
			server, err := config.GetServer(host)
			if err != nil {
				fmt.Println(err.Error())
				continue
			}

			err = actions.ServerRemove(*server)
			if err != nil {
				fmt.Println(locales.L.Get("errors.server.remove.failure", host, err.Error()))
			} else {
				fmt.Println("server.delete.success")
			}
		}
	},
	Args: cobra.MinimumNArgs(1),
}

func init() {
	ServersCmd.AddCommand(deleteCmd)
}
