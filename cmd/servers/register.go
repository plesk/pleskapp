// Copyright 1999-2020. Plesk International GmbH.

package cmd

import (
	"git.plesk.ru/projects/SBX/repos/pleskapp/actions"
	"git.plesk.ru/projects/SBX/repos/pleskapp/locales"
	"git.plesk.ru/projects/SBX/repos/pleskapp/utils"
	"github.com/spf13/cobra"
)

var registerCmd = &cobra.Command{
	Use:   locales.L.Get("server.register.cmd"),
	Short: locales.L.Get("server.register.description"),
	RunE: func(cmd *cobra.Command, args []string) error {
		ignoreSsl, _ := cmd.Flags().GetBool("ignoreSsl")

		cmd.SilenceUsage = true
		return utils.Log.PrintSuccessOrError("server.register.success", nil, actions.ServerAdd(args[0], ignoreSsl))
	},
	Args: cobra.ExactArgs(1),
}

func init() {
	registerCmd.Flags().BoolP("ignoreSsl", "s", false, locales.L.Get("server.register.ignore.ssl.flag"))
	ServersCmd.AddCommand(registerCmd)
}
