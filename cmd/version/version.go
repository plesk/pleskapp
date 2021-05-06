// Copyright 1999-2021. Plesk International GmbH.

package cmd

import (
	"fmt"

	"github.com/plesk/pleskapp/plesk/locales"
	"github.com/spf13/cobra"
)

var (
	Revision  string
	BuildTime string
)

var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: locales.L.Get("version.description"),
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("Revision:\t%s\nBuild time:\t%s\n", Revision, BuildTime)

		return nil
	},
}
