// Copyright 1999-2021. Plesk International GmbH.

package cmd

import (
	"errors"

	"github.com/plesk/pleskapp/plesk/internal/actions"
	"github.com/plesk/pleskapp/plesk/internal/config"
	"github.com/plesk/pleskapp/plesk/internal/locales"
	"github.com/plesk/pleskapp/plesk/internal/types"
	"github.com/plesk/pleskapp/plesk/internal/utils"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create [SERVER] [DOMAIN] [NAME]",
	Short: locales.L.Get("database.create.description"),
	RunE: func(cmd *cobra.Command, args []string) error {
		var realdbt string = ""
		dbt, _ := cmd.Flags().GetString("type")
		for _, i := range []string{"mysql", "mssql", "postgresql"} {
			if dbt == i {
				realdbt = dbt
			}
		}

		if realdbt == "" {
			return errors.New(locales.L.Get("errors.unknown.database.type", dbt))
		}

		server, err := config.GetServer(args[0])
		if err != nil {
			return err
		}

		dbs := server.GetDatabaseServerByType(realdbt)
		if dbs == nil {
			return types.DatabaseServerNotFound{
				DbType: realdbt,
				Server: server.Host,
			}
		}

		domain, err := config.GetDomain(*server, args[1])
		if err != nil {
			return err
		}

		db := types.NewDatabase{
			Name:             args[2],
			Type:             realdbt,
			ParentDomain:     domain.Name,
			DatabaseServerID: dbs.ID,
		}

		cmd.SilenceUsage = true
		err = actions.DatabaseAdd(*server, *domain, *dbs, db)

		if err == nil {
			utils.Log.PrintL("database.create.success", db.Name)
		}

		return err
	},
	Args: cobra.ExactArgs(3),
}

func init() {
	createCmd.Flags().String("type", "mysql", locales.L.Get("database.create.flag.type"))
	DatabasesCmd.AddCommand(createCmd)
}
