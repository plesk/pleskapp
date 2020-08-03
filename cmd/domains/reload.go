// Copyright 1999-2020. Plesk International GmbH.

package cmd

import (
	"github.com/plesk/pleskapp/actions"
	"github.com/plesk/pleskapp/config"
	"github.com/plesk/pleskapp/locales"
	"github.com/plesk/pleskapp/utils"
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

		return utils.Log.PrintSuccessOrError("domain.reload.success", nil, actions.DomainReload(*server))
	},
	Args: cobra.ExactArgs(1),
}

func init() {
	DomainsCmd.AddCommand(reloadCmd)
}
