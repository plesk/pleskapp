// Copyright 1999-2021. Plesk International GmbH.

package cmd

import (
	"github.com/plesk/pleskapp/plesk/locales"
	"github.com/spf13/cobra"
)

var DomainsCmd = &cobra.Command{
	Use:   "domains",
	Short: locales.L.Get("domain.description"),
	RunE: func(cmd *cobra.Command, args []string) error {
		return listCmd.RunE(cmd, args)
	},
}
