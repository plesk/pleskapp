// Copyright 1999-2020. Plesk International GmbH. All rights reserved.

package cmd

import (
	"git.plesk.ru/~abashurov/pleskapp/actions"
	"git.plesk.ru/~abashurov/pleskapp/locales"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   locales.L.Get("server.list.cmd"),
	Short: locales.L.Get("server.list.description"),
	Run: func(cmd *cobra.Command, args []string) {
		actions.ServerList()
	},
	Args: cobra.ExactArgs(0),
}

func init() {
	ServersCmd.AddCommand(listCmd)
}
