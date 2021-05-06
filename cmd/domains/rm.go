// Copyright 1999-2021. Plesk International GmbH.

package cmd

import (
	"github.com/plesk/pleskapp/plesk/actions"
	"github.com/plesk/pleskapp/plesk/config"
	"github.com/plesk/pleskapp/plesk/locales"
	"github.com/plesk/pleskapp/plesk/utils"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   locales.L.Get("domain.delete.cmd"),
	Short: locales.L.Get("domain.delete.description"),
	RunE: func(cmd *cobra.Command, args []string) error {
		var lastErr error = nil
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
					lastErr = err
					utils.Log.Error(err.Error())
				}
			}
		}

		cmd.SilenceUsage = true
		if lastErr == nil {
			utils.Log.PrintL("domain.delete.success")
		}

		return lastErr
	},
	Args: cobra.MaximumNArgs(2),
}

func init() {
	DomainsCmd.AddCommand(deleteCmd)
}
