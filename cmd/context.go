// Copyright 1999-2021. Plesk International GmbH.

package cmd

import (
	"fmt"
	"github.com/plesk/pleskapp/plesk/config"
	"github.com/plesk/pleskapp/plesk/locales"
	"github.com/plesk/pleskapp/plesk/types"
	"github.com/plesk/pleskapp/plesk/utils"
	"github.com/spf13/cobra"
)

var contextCmd = &cobra.Command{
	Use:   "context [SERVER]",
	Short: locales.L.Get("context.description"),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 1 {
			server, err := config.GetServer(args[0])
			if err != nil {
				return err
			}
			fmt.Println("Updating the context...")
			servers, _ := utils.FilterServers(config.GetServers(), server.Host)
			config.SetServers(append([]types.Server{*server}, servers...))
		}

		defaultServer, err := config.DefaultServer()
		if err != nil {
			return err
		}

		fmt.Printf("Default context (Plesk server): %s\n", defaultServer)

		return nil
	},
}
