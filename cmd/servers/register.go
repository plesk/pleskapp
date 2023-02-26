// Copyright 1999-2023. Plesk International GmbH.

package cmd

import (
	"fmt"
	"github.com/plesk/pleskapp/plesk/internal/actions"
	"github.com/plesk/pleskapp/plesk/internal/locales"
	"github.com/plesk/pleskapp/plesk/internal/types"
	"github.com/spf13/cobra"
)

var apiKey string

var registerCmd = &cobra.Command{
	Use:   "register [IP ADDRESS|HOSTNAME]",
	Short: locales.L.Get("server.register.description"),
	RunE: func(cmd *cobra.Command, args []string) error {
		ignoreSsl, _ := cmd.Flags().GetBool("ignore-ssl")
		apiKey, _ := cmd.Flags().GetString("key")
		host := args[0]

		var err error
		if apiKey != "" {
			err = actions.ServerUpdate(types.Server{
				Host:      host,
				IgnoreSsl: ignoreSsl,
				APIKey:    apiKey,
			})
		} else {
			err = actions.ServerAdd(host, ignoreSsl)
		}

		cmd.SilenceUsage = true

		if err == nil {
			fmt.Println(locales.L.Get("server.register.success"))
		}

		return err
	},
	Args: cobra.ExactArgs(1),
}

func init() {
	registerCmd.Flags().BoolP("ignore-ssl", "s", true, locales.L.Get("server.register.ignore.ssl.flag"))
	registerCmd.Flags().StringVarP(&apiKey, "key", "k", "", locales.L.Get("server.register.api-key"))
	ServersCmd.AddCommand(registerCmd)
}
