// Copyright 1999-2022. Plesk International GmbH.

package cmd

import (
	"fmt"
	"github.com/plesk/pleskapp/plesk/internal/actions"
	"github.com/plesk/pleskapp/plesk/internal/config"
	"github.com/plesk/pleskapp/plesk/internal/locales"
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
			fmt.Println(locales.L.Get("database.delete.success", args[1]))
		}

		return err
	},
	Args: cobra.ExactArgs(2),
}

func init() {
	DatabasesCmd.AddCommand(deleteCmd)
}
