// Copyright 1999-2021. Plesk International GmbH.

package cmd

import (
	"github.com/plesk/pleskapp/plesk/actions"
	"github.com/plesk/pleskapp/plesk/locales"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   locales.L.Get("server.list.cmd"),
	Short: locales.L.Get("server.list.description"),
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true
		return actions.ServerList()
	},
	Args: cobra.ExactArgs(0),
}

func init() {
	ServersCmd.AddCommand(listCmd)
}
