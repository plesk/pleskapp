// Copyright 1999-2020. Plesk International GmbH. All rights reserved.

package cmd

import (
	"strings"

	"git.plesk.ru/~abashurov/pleskapp/actions"
	"git.plesk.ru/~abashurov/pleskapp/config"
	"git.plesk.ru/~abashurov/pleskapp/locales"
	"git.plesk.ru/~abashurov/pleskapp/types"
	"git.plesk.ru/~abashurov/pleskapp/utils"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   locales.L.Get("domain.create.cmd"),
	Short: locales.L.Get("domain.create.description"),
	RunE: func(cmd *cobra.Command, args []string) error {
		addr4 := []string{}
		addr6 := []string{}

		for i, addr := range args {
			if i > 1 {
				if strings.Contains(addr, ".") {
					addr4 = append(addr4, addr)
				} else if strings.Contains(addr, ":") {
					addr6 = append(addr6, addr)
				}
			}
		}

		ips := types.ServerIPAddresses{
			IPv4: addr4,
			IPv6: addr6,
		}

		server, err := config.GetServer(args[0])
		if err != nil {
			return err
		}

		cmd.SilenceUsage = true
		return utils.Log.PrintSuccessOrError("domain.create.success", nil, actions.DomainAdd(*server, args[1], ips))
	},
	Args: cobra.RangeArgs(3, 4),
}

func init() {
	DomainsCmd.AddCommand(createCmd)
}
