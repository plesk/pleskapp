// Copyright 1999-2020. Plesk International GmbH.

package cmd

import (
	"github.com/plesk/pleskapp/actions"
	"github.com/plesk/pleskapp/config"
	"github.com/plesk/pleskapp/locales"
	"github.com/plesk/pleskapp/utils"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   locales.L.Get("domain.delete.cmd"),
	Short: locales.L.Get("domain.delete.description"),
	RunE: func(cmd *cobra.Command, args []string) error {
		server, err := config.GetServer(args[0])
		if err != nil {
			return err
		}

		for i, d := range args {
			if i > 0 {
				domain, err := config.GetDomain(*server, d)
				if err != nil {
					utils.Log.Error(err.Error())
					continue
				}

				err = actions.DomainDelete(*server, *domain)
				if err != nil {
					utils.Log.Error(err.Error())
				}
			}
		}

		cmd.SilenceUsage = true
		return utils.Log.PrintSuccessOrError("domain.delete.success", nil, err)
	},
	Args: cobra.MaximumNArgs(2),
}

func init() {
	DomainsCmd.AddCommand(deleteCmd)
}
