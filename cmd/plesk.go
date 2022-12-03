// Copyright 1999-2022. Plesk International GmbH.

package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var pleskCmd = &cobra.Command{
	Use:    "plesk",
	Hidden: true,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("  ____    _                _      _\n" +
			" |  _ \\  | |   ___   ___  | | __ | |\n" +
			" | |_) | | |  / _ \\ / __| | |/ / | |\n" +
			" |  __/  | | |  __/ \\__ \\ |   <  |_|\n" +
			" |_|     |_|  \\___| |___/ |_|\\_\\ (_)")
	},
}
