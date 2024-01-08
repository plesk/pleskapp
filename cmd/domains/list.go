// Copyright 1999-2024. Plesk International GmbH.

package cmd

import (
	"github.com/plesk/pleskapp/plesk/internal/actions"
	"github.com/plesk/pleskapp/plesk/internal/config"
	"github.com/plesk/pleskapp/plesk/internal/locales"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list [SERVER]",
	Short: locales.L.Get("domain.list.description"),
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true

		server, err := config.GetServerByArgs(args)
		if err != nil {
			return err
		}

		return actions.DomainList(*server)
	},
	Args: cobra.MaximumNArgs(1),
}

func init() {
	DomainsCmd.AddCommand(listCmd)
}
