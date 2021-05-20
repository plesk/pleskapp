// Copyright 1999-2021. Plesk International GmbH.

package cmd

import (
	serversCmd "github.com/plesk/pleskapp/plesk/cmd/servers"
	"github.com/plesk/pleskapp/plesk/locales"
	"github.com/spf13/cobra"
)

var loginCmd = &cobra.Command{
	Use:   "login [SERVER]",
	Short: locales.L.Get("server.login.description"),
	RunE: func(cmd *cobra.Command, args []string) error {
		return serversCmd.ServerLoginCmd.RunE(cmd, args)
	},
	Args: cobra.MaximumNArgs(1),
}

func init() {
	loginCmd.Flags().BoolP("generate", "g", false, locales.L.Get("server.login.generate.flag"))
}
