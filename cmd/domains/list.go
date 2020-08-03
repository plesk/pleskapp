// Copyright 1999-2020. Plesk International GmbH.

package cmd

import (
	"git.plesk.ru/projects/SBX/repos/pleskapp/actions"
	"git.plesk.ru/projects/SBX/repos/pleskapp/config"
	"git.plesk.ru/projects/SBX/repos/pleskapp/locales"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   locales.L.Get("domain.list.cmd"),
	Short: locales.L.Get("domain.list.description"),
	RunE: func(cmd *cobra.Command, args []string) error {
		server, err := config.GetServer(args[0])
		if err != nil {
			return err
		}

		return actions.DomainList(*server)
	},
	Args: cobra.ExactArgs(1),
}

func init() {
	DomainsCmd.AddCommand(listCmd)
}
