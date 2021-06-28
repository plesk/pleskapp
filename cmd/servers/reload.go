// Copyright 1999-2021. Plesk International GmbH.

package cmd

import (
	"fmt"
	"github.com/plesk/pleskapp/plesk/internal/actions"
	"github.com/plesk/pleskapp/plesk/internal/config"
	"github.com/plesk/pleskapp/plesk/internal/locales"
	"github.com/spf13/cobra"
)

var reloadCmd = &cobra.Command{
	Use:   "reload [IP ADDRESS|HOSTNAME]",
	Short: locales.L.Get("server.reload.description"),
	RunE: func(cmd *cobra.Command, args []string) error {
		server, err := config.GetServerByArgs(args)
		if err != nil {
			return err
		}

		cmd.SilenceUsage = true
		err = actions.ServerUpdate(*server)

		if err == nil {
			fmt.Println(locales.L.Get("server.reload.success"))
		}

		return err
	},
	Args: cobra.MaximumNArgs(1),
}

func init() {
	ServersCmd.AddCommand(reloadCmd)
}
