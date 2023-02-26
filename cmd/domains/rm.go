// Copyright 1999-2023. Plesk International GmbH.

package cmd

import (
	"fmt"
	"github.com/plesk/pleskapp/plesk/internal/actions"
	"github.com/plesk/pleskapp/plesk/internal/config"
	"github.com/plesk/pleskapp/plesk/internal/locales"
	"github.com/spf13/cobra"
	"log"
)

var deleteCmd = &cobra.Command{
	Use:   "delete [SERVER] [DOMAIN...]",
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
					log.Println(err.Error())
					continue
				}

				err = actions.DomainDelete(*server, *domain)
				if err != nil {
					lastErr = err
					log.Println(err.Error())
				}
			}
		}

		cmd.SilenceUsage = true
		if lastErr == nil {
			fmt.Println(locales.L.Get("domain.delete.success"))
		}

		return lastErr
	},
	Args: cobra.MaximumNArgs(2),
}

func init() {
	DomainsCmd.AddCommand(deleteCmd)
}
