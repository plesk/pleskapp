// Copyright 1999-2021. Plesk International GmbH.

package cmd

import (
	"strings"

	"github.com/plesk/pleskapp/plesk/internal/actions"
	"github.com/plesk/pleskapp/plesk/internal/config"
	"github.com/plesk/pleskapp/plesk/internal/locales"
	"github.com/plesk/pleskapp/plesk/internal/types"
	"github.com/plesk/pleskapp/plesk/internal/utils"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create [SERVER] [DOMAIN] [IPv4] [IPv6]",
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
		err = actions.DomainAdd(*server, args[1], ips)

		if err == nil {
			utils.Log.PrintL("domain.create.success", args[1])
		}

		return err
	},
	Args: cobra.RangeArgs(3, 4),
}

func init() {
	DomainsCmd.AddCommand(createCmd)
}
