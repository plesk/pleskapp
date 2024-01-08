// Copyright 1999-2024. Plesk International GmbH.

package cmd

import (
	"fmt"
	"github.com/plesk/pleskapp/plesk/internal/actions"
	"github.com/plesk/pleskapp/plesk/internal/config"
	"github.com/plesk/pleskapp/plesk/internal/locales"
	"github.com/spf13/cobra"
)

var reauthCmd = &cobra.Command{
	Use:   "reauth [IP ADDRESS|HOSTNAME]",
	Short: locales.L.Get("server.reauth.description"),
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true

		server, err := config.GetServerByArgs(args)
		if err != nil {
			return err
		}

		err = actions.ServerReauth(*server)

		if err == nil {
			fmt.Println(locales.L.Get("server.reauth.success"))
		}

		return err
	},
	Args: cobra.MaximumNArgs(1),
}

func init() {
	ServersCmd.AddCommand(reauthCmd)
}
