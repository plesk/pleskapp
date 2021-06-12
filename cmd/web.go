// Copyright 1999-2021. Plesk International GmbH.

package cmd

import (
	"fmt"
	"github.com/plesk/pleskapp/plesk/internal/actions"
	"github.com/plesk/pleskapp/plesk/internal/locales"
	"github.com/spf13/cobra"
)

var webCmd = &cobra.Command{
	Use:   "web",
	Short: locales.L.Get("web.description"),
	RunE: func(cmd *cobra.Command, args []string) error {
		stack := actions.DetectStack()
		fmt.Println("Detected stack:", stack)
		return actions.RunServer(stack)
	},
}
