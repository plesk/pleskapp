// Copyright 1999-2023. Plesk International GmbH.

package cmd

import (
	"errors"
	"fmt"
	"github.com/plesk/pleskapp/plesk/internal/actions"
	"github.com/plesk/pleskapp/plesk/internal/config"
	"github.com/plesk/pleskapp/plesk/internal/locales"
	"github.com/plesk/pleskapp/plesk/internal/types"
	"github.com/spf13/cobra"
	"os"
	"path"
	"strings"
)

var deployCmd = &cobra.Command{
	Use:   "deploy [SERVER]",
	Short: locales.L.Get("deploy.description"),
	RunE: func(cmd *cobra.Command, args []string) error {
		workDir, _ := os.Getwd()
		appName := path.Base(workDir)
		fmt.Println("App name:", appName)

		server, err := config.GetServerByArgs(args)
		if err != nil {
			return err
		}

		if len(server.Info.IP.IPv4) == 0 {
			return errors.New("IPv4 address is required")
		}
		defaultIp := server.Info.IP.IPv4[0]
		fmt.Printf("IP: %v\n", defaultIp)

		domainName := fmt.Sprintf("%s.%s.plesk.page", appName, strings.ReplaceAll(defaultIp, ".", "-"))
		fmt.Println("Domain name:", domainName)

		var domain *types.Domain
		domain, err = config.GetDomain(*server, domainName)
		if err != nil {
			fmt.Printf("Creating the domain %s...\n", domainName)
			err = actions.DomainAdd(server, domainName, types.ServerIPAddresses{
				IPv4: []string{defaultIp},
			})
			if err != nil {
				return err
			}
			fmt.Println("Domain has been created.")
			domain, err = config.GetDomain(*server, domainName)
			if err != nil {
				return err
			}
		} else {
			fmt.Printf("Domain %s has been found.\n", domainName)
		}

		fmt.Println("Uploading the content...")
		err = actions.UploadDirectory(*server, *domain, true, false, workDir, nil)
		if err != nil {
			return err
		}

		return nil
	},
}
