// Copyright 1999-2022. Plesk International GmbH.

package cmd

import (
	"github.com/plesk/pleskapp/plesk/internal/locales"
	"github.com/spf13/cobra"
)

var ServersCmd = &cobra.Command{
	Use:   "servers",
	Short: locales.L.Get("server.description"),
	RunE: func(cmd *cobra.Command, args []string) error {
		return listCmd.RunE(cmd, args)
	},
}
