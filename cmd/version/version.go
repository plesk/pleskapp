// Copyright 1999-2021. Plesk International GmbH.

package cmd

import (
	"fmt"
	"github.com/plesk/pleskapp/plesk/actions"
	"github.com/plesk/pleskapp/plesk/config"

	"github.com/plesk/pleskapp/plesk/locales"
	"github.com/spf13/cobra"
)

var (
	Revision  string
	BuildTime string
	Version   string
)

var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: locales.L.Get("version.description"),
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Client information")
		fmt.Printf("Version:\t%s\nRevision:\t%s\nBuild time:\t%s\n", Version, Revision, BuildTime)

		defaultServerName, err := actions.DefaultServer()
		if err == nil {
			server, _ := config.GetServer(defaultServerName)
			fmt.Println()
			fmt.Println("Server information")
			fmt.Printf("Host:   \t%s\n", server.Host)
			fmt.Printf("Version:\t%s\n", server.Info.Version)
		}

		return nil
	},
}
