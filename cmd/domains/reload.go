// Copyright 1999-2021. Plesk International GmbH.

package cmd

import (
	"github.com/plesk/pleskapp/plesk/actions"
	"github.com/plesk/pleskapp/plesk/config"
	"github.com/plesk/pleskapp/plesk/locales"
	"github.com/plesk/pleskapp/plesk/utils"
	"github.com/spf13/cobra"
)

var reloadCmd = &cobra.Command{
	Use:   locales.L.Get("domain.reload.cmd"),
	Short: locales.L.Get("domain.reload.description"),
	RunE: func(cmd *cobra.Command, args []string) error {
		server, err := config.GetServer(args[0])
		if err != nil {
			return err
		}

		cmd.SilenceUsage = true
		err = actions.DomainReload(*server)

		if err == nil {
			utils.Log.PrintL("domain.reload.success")
		}

		return err
	},
	Args: cobra.ExactArgs(1),
}

func init() {
	DomainsCmd.AddCommand(reloadCmd)
}
