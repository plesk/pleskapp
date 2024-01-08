// Copyright 1999-2024. Plesk International GmbH.

package cmd

import (
	"github.com/plesk/pleskapp/plesk/internal/actions"
	"github.com/plesk/pleskapp/plesk/internal/config"
	"github.com/plesk/pleskapp/plesk/internal/locales"
	"github.com/spf13/cobra"
)

var ServerLoginCmd = &cobra.Command{
	Use:   "login [SERVER]",
	Short: locales.L.Get("server.login.description"),
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true

		generateOnly, _ := cmd.Flags().GetBool("generate")

		server, err := config.GetServerByArgs(args)
		if err != nil {
			return err
		}

		return actions.ServerLogin(*server, generateOnly)
	},
	Args: cobra.MaximumNArgs(1),
}

func init() {
	ServerLoginCmd.Flags().BoolP("generate", "g", false, locales.L.Get("server.login.generate.flag"))

	ServersCmd.AddCommand(ServerLoginCmd)
}
