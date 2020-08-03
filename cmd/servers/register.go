// Copyright 1999-2020. Plesk International GmbH.

package cmd

import (
	"github.com/plesk/pleskapp/plesk/actions"
	"github.com/plesk/pleskapp/plesk/locales"
	"github.com/plesk/pleskapp/plesk/utils"
	"github.com/spf13/cobra"
)

var registerCmd = &cobra.Command{
	Use:   locales.L.Get("server.register.cmd"),
	Short: locales.L.Get("server.register.description"),
	RunE: func(cmd *cobra.Command, args []string) error {
		ignoreSsl, _ := cmd.Flags().GetBool("ignore-ssl")

		cmd.SilenceUsage = true
		err := actions.ServerAdd(args[0], ignoreSsl)

		if err == nil {
			utils.Log.PrintL("server.register.success")
		}

		return err
	},
	Args: cobra.ExactArgs(1),
}

func init() {
	registerCmd.Flags().BoolP("ignore-ssl", "s", true, locales.L.Get("server.register.ignore.ssl.flag"))
	ServersCmd.AddCommand(registerCmd)
}
