// Copyright 1999-2020. Plesk International GmbH.

package cmd

import (
	"errors"

	"git.plesk.ru/projects/SBX/repos/pleskapp/actions"
	"git.plesk.ru/projects/SBX/repos/pleskapp/config"
	"git.plesk.ru/projects/SBX/repos/pleskapp/locales"
	"git.plesk.ru/projects/SBX/repos/pleskapp/types"
	"git.plesk.ru/projects/SBX/repos/pleskapp/utils"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   locales.L.Get("database.create.cmd"),
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
		return utils.Log.PrintSuccessOrError("database.create.success", nil, actions.DatabaseAdd(*server, *domain, *dbs, db))
	},
	Args: cobra.ExactArgs(3),
}

func init() {
	createCmd.Flags().String("type", "mysql", locales.L.Get("database.create.flag.type"))
	DatabasesCmd.AddCommand(createCmd)
}
