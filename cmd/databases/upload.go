// Copyright 1999-2021. Plesk International GmbH.

package cmd

import (
	"errors"
	"os"
	"strings"

	"github.com/plesk/pleskapp/plesk/actions"
	"github.com/plesk/pleskapp/plesk/api/factory"
	"github.com/plesk/pleskapp/plesk/config"
	"github.com/plesk/pleskapp/plesk/locales"
	"github.com/plesk/pleskapp/plesk/utils"
	"github.com/spf13/cobra"
)

var uploadCmd = &cobra.Command{
	Use:   "upload [SERVER] [DOMAIN] [NAME] [FILE]",
	Short: locales.L.Get("database.deploy.description"),
	RunE: func(cmd *cobra.Command, args []string) error {
		server, err := config.GetServer(args[0])
		if err != nil {
			return err
		}

		domain, err := config.GetDomain(*server, args[1])
		if err != nil {
			return err
		}

		// TODO: Restrict to one domain
		db, err := actions.DatabaseFindNonLocal(factory.GetDatabaseManagement(server.GetServerAuth()), *server, args[2])
		if err != nil {
			return err
		}

		s, err := os.Stat(args[3])
		if err != nil {
			return err
		}

		if s.IsDir() {
			return errors.New(locales.L.Get("errors.path.is.directory", args[3]))
		}

		cmd.SilenceUsage = true
		path, err := actions.UploadFileToRoot(*server, *domain, true, args[3])
		if err != nil {
			return err
		}

		fp := strings.Split(args[3], "/")
		err = actions.DatabaseDeploy(*server, *domain, *db, path+"/"+fp[len(fp)-1])

		if err == nil {
			utils.Log.PrintL("database.deploy.success", db.Name, args[3])
		}

		return err
	},
	Args: cobra.ExactArgs(4),
}

func init() {
	DatabasesCmd.AddCommand(uploadCmd)
}
