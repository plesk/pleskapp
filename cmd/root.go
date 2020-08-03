// Copyright 1999-2020. Plesk International GmbH.

package cmd

import (
	appsCmd "git.plesk.ru/projects/SBX/repos/pleskapp/cmd/apps"
	databasesCmd "git.plesk.ru/projects/SBX/repos/pleskapp/cmd/databases"
	domainsCmd "git.plesk.ru/projects/SBX/repos/pleskapp/cmd/domains"
	serversCmd "git.plesk.ru/projects/SBX/repos/pleskapp/cmd/servers"
	syncCmd "git.plesk.ru/projects/SBX/repos/pleskapp/cmd/sync"
	"git.plesk.ru/projects/SBX/repos/pleskapp/utils"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:           "pleskapp",
	Short:         "Manage Plesk servers from the local system",
	SilenceErrors: true,
}

func Execute() error {
	rootCmd.AddCommand(
		appsCmd.AppsCmd,
		databasesCmd.DatabasesCmd,
		domainsCmd.DomainsCmd,
		serversCmd.ServersCmd,
		syncCmd.SyncCmd,
	)

	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initLogger)
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Enable verbose output")
}

func initLogger() {
	v, _ := rootCmd.PersistentFlags().GetBool("verbose")
	if v {
		utils.Log.SetLevel(2)
	}
}
