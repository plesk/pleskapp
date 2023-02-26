// Copyright 1999-2023. Plesk International GmbH.

package cmd

import (
	"github.com/plesk/pleskapp/plesk/internal/actions"
	"github.com/plesk/pleskapp/plesk/internal/config"
	"github.com/plesk/pleskapp/plesk/internal/locales"
	"github.com/spf13/cobra"
)

var SSHCmd = &cobra.Command{
	Use:   "ssh [SERVER]",
	Short: locales.L.Get("server.ssh.description"),
	RunE: func(cmd *cobra.Command, args []string) error {
		additionalCommand, _ := cmd.Flags().GetString("command")
		cmd.SilenceUsage = true

		server, err := config.GetServerByArgs(args)
		if err != nil {
			return err
		}

		return actions.ServerSSH(*server, additionalCommand)
	},
}

func init() {
	SSHCmd.Flags().StringP("command", "c", "", locales.L.Get("server.ssh.flag.command"))

	ServersCmd.AddCommand(SSHCmd)
}
