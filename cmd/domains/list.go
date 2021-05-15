// Copyright 1999-2021. Plesk International GmbH.

package cmd

import (
	"github.com/plesk/pleskapp/plesk/actions"
	"github.com/plesk/pleskapp/plesk/config"
	"github.com/plesk/pleskapp/plesk/locales"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list [SERVER]",
	Short: locales.L.Get("domain.list.description"),
	RunE: func(cmd *cobra.Command, args []string) error {
		server, err := config.GetServer(args[0])
		if err != nil {
			return err
		}

		cmd.SilenceUsage = true
		return actions.DomainList(*server)
	},
	Args: cobra.ExactArgs(1),
}

func init() {
	DomainsCmd.AddCommand(listCmd)
}
