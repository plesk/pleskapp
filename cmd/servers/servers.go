// Copyright 1999-2021. Plesk International GmbH.

package cmd

import (
	"github.com/plesk/pleskapp/plesk/actions"
	"github.com/plesk/pleskapp/plesk/config"
	"github.com/plesk/pleskapp/plesk/locales"
	"github.com/plesk/pleskapp/plesk/types"
	"github.com/spf13/cobra"
)

var ServersCmd = &cobra.Command{
	Use:   "servers",
	Short: locales.L.Get("server.description"),
	RunE: func(cmd *cobra.Command, args []string) error {
		return listCmd.RunE(cmd, args)
	},
}

func getServer(args []string) (*types.Server, error) {
	var serverName string
	if len(args) == 0 {
		serverName, _ = actions.DefaultServer()
	} else {
		serverName = args[0]
	}
	server, err := config.GetServer(serverName)
	if err != nil {
		return nil, err
	}
	return server, nil
}