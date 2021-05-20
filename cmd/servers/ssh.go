// Copyright 1999-2021. Plesk International GmbH.

package cmd

import (
	"github.com/plesk/pleskapp/plesk/actions"
	"github.com/plesk/pleskapp/plesk/locales"
	"github.com/spf13/cobra"
)

var SshCmd = &cobra.Command{
	Use:   "ssh [SERVER]",
	Short: locales.L.Get("server.ssh.description"),
	RunE: func(cmd *cobra.Command, args []string) error {
		server, err := getServer(args)
		if err != nil {
			return err
		}

		cmd.SilenceUsage = true

		return actions.ServerSsh(*server)
	},
	Args: cobra.MaximumNArgs(1),
}

func init() {
	ServersCmd.AddCommand(SshCmd)
}