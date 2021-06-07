// Copyright 1999-2021. Plesk International GmbH.

package cmd

import (
	appsCmd "github.com/plesk/pleskapp/plesk/cmd/apps"
	databasesCmd "github.com/plesk/pleskapp/plesk/cmd/databases"
	domainsCmd "github.com/plesk/pleskapp/plesk/cmd/domains"
	serversCmd "github.com/plesk/pleskapp/plesk/cmd/servers"
	syncCmd "github.com/plesk/pleskapp/plesk/cmd/sync"
	"github.com/plesk/pleskapp/plesk/utils"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:           "plesk",
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
		versionCmd,
		contextCmd,
		loginCmd,
		sshCmd,
		webCmd,
		completionCmd,
		pleskCmd,
		deployCmd,
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
