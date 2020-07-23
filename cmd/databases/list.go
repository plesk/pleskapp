// Copyright 1999-2020. Plesk International GmbH. All rights reserved.

package cmd

import (
	"git.plesk.ru/projects/SBX/repos/pleskapp/actions"
	"git.plesk.ru/projects/SBX/repos/pleskapp/config"
	"git.plesk.ru/projects/SBX/repos/pleskapp/locales"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   locales.L.Get("database.list.cmd"),
	Short: locales.L.Get("database.list.description"),
	RunE: func(cmd *cobra.Command, args []string) error {
		server, err := config.GetServer(args[0])
		if err != nil {
			return err
		}

		domain, err := config.GetDomain(*server, args[1])
		if err != nil {
			return err
		}

		return actions.DatabaseList(*server, *domain)
	},
	Args: cobra.ExactArgs(2),
}

func init() {
	DatabasesCmd.AddCommand(listCmd)
}
