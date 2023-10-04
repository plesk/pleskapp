// Copyright 1999-2023. Plesk International GmbH.

package cmd

import (
	"github.com/plesk/pleskapp/plesk/internal/actions"
	"github.com/plesk/pleskapp/plesk/internal/config"
	"github.com/plesk/pleskapp/plesk/internal/locales"
	"github.com/spf13/cobra"
)

var dbCmd = &cobra.Command{
	Use:   "db [SERVER]",
	Short: locales.L.Get("server.db.description"),
	RunE: func(cmd *cobra.Command, args []string) error {
		server, err := config.GetServerByArgs(args)
		if err != nil {
			return err
		}

		return actions.ServerSSH(*server, "plesk db", true)
	},
}
