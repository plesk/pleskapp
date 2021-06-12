// Copyright 1999-2021. Plesk International GmbH.

package cmd

import (
	"github.com/plesk/pleskapp/plesk/internal/actions"
	"github.com/plesk/pleskapp/plesk/internal/config"
	"github.com/plesk/pleskapp/plesk/internal/locales"
	"github.com/plesk/pleskapp/plesk/internal/utils"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete [SERVER] [NAME]",
	Short: locales.L.Get("database.delete.description"),
	RunE: func(cmd *cobra.Command, args []string) error {
		server, err := config.GetServer(args[0])
		if err != nil {
			return err
		}

		cmd.SilenceUsage = true
		err = actions.DatabaseDelete(*server, args[1])

		if err == nil {
			utils.Log.PrintL("database.delete.success", args[1])
		}

		return err
	},
	Args: cobra.ExactArgs(2),
}

func init() {
	DatabasesCmd.AddCommand(deleteCmd)
}
