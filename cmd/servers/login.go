// Copyright 1999-2021. Plesk International GmbH.

package cmd

import (
	"github.com/plesk/pleskapp/plesk/actions"
	"github.com/plesk/pleskapp/plesk/config"
	"github.com/plesk/pleskapp/plesk/locales"
	"github.com/spf13/cobra"
)

var loginCmd = &cobra.Command{
	Use:   "login [SERVER]",
	Short: locales.L.Get("server.login.description"),
	RunE: func(cmd *cobra.Command, args []string) error {
		generateOnly, _ := cmd.Flags().GetBool("generate")

		server, err := config.GetServer(args[0])
		if err != nil {
			return err
		}

		cmd.SilenceUsage = true
		return actions.ServerLogin(*server, generateOnly)
	},
	Args: cobra.ExactArgs(1),
}

func init() {
	loginCmd.Flags().BoolP("generate", "g", false, locales.L.Get("server.login.generate.flag"))

	ServersCmd.AddCommand(loginCmd)
}
