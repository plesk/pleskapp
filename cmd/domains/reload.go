// Copyright 1999-2023. Plesk International GmbH.

package cmd

import (
	"fmt"
	"github.com/plesk/pleskapp/plesk/internal/actions"
	"github.com/plesk/pleskapp/plesk/internal/config"
	"github.com/plesk/pleskapp/plesk/internal/locales"
	"github.com/spf13/cobra"
)

var reloadCmd = &cobra.Command{
	Use:   "reload [SERVER]",
	Short: locales.L.Get("domain.reload.description"),
	RunE: func(cmd *cobra.Command, args []string) error {
		server, err := config.GetServer(args[0])
		if err != nil {
			return err
		}

		cmd.SilenceUsage = true
		err = actions.DomainReload(server)

		if err == nil {
			fmt.Println(locales.L.Get("domain.reload.success"))
		}

		return err
	},
	Args: cobra.ExactArgs(1),
}

func init() {
	DomainsCmd.AddCommand(reloadCmd)
}
