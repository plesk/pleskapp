// Copyright 1999-2021. Plesk International GmbH.

package cmd

import (
	serversCmd "github.com/plesk/pleskapp/plesk/cmd/servers"
	"github.com/plesk/pleskapp/plesk/locales"
	"github.com/spf13/cobra"
)

var sshCmd = &cobra.Command{
	Use:   "ssh [SERVER]",
	Short: locales.L.Get("server.ssh.description"),
	RunE: func(cmd *cobra.Command, args []string) error {
		return serversCmd.SSHCmd.RunE(cmd, args)
	},
	Args: cobra.MaximumNArgs(1),
}
