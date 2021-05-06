// Copyright 1999-2021. Plesk International GmbH.

package actions

import (
	"fmt"

	"github.com/plesk/pleskapp/plesk/api"
	"github.com/plesk/pleskapp/plesk/api/factory"
	"github.com/plesk/pleskapp/plesk/config"
	"github.com/plesk/pleskapp/plesk/types"
	"github.com/plesk/pleskapp/plesk/utils"
)

func DatabaseFindNonLocal(api api.DatabaseManagement, host types.Server, dbn string) (*types.Database, error) {
	db, err := config.GetDatabase(host, dbn)
	if err != nil {
		dbs, err := api.ListDatabases()
		if err != nil {
			return db, err
		}

		for _, d := range dbs {
			if d.Name == dbn {
				db = &types.Database{
					ID:               d.ID,
					Name:             d.Name,
					Type:             d.Type,
					DatabaseServerID: d.DatabaseServerID,
				}
			}
		}
	}

	return db, nil
}

func DatabaseAdd(host types.Server, domain types.Domain, dbs types.DatabaseServer, db types.NewDatabase) error {
	api := factory.GetDatabaseManagement(host.GetServerAuth())
	newdb, err := api.CreateDatabase(domain, db, dbs)
	if err != nil {
		return err
	}

	domain.Databases = append(domain.Databases, types.Database{
		ID:               newdb.ID,
		Name:             newdb.Name,
		Type:             newdb.Type,
		DatabaseServerID: newdb.DatabaseServerID,
	})

	config.SetDomain(host, domain)
	return nil
}

func DatabaseList(host types.Server, domain types.Domain) error {
	api := factory.GetDatabaseManagement(host.GetServerAuth())
	dbs, err := api.ListDomainDatabases(domain.Name)
	if err != nil {
		return err
	}

	for _, i := range dbs {
		fmt.Printf("ID: %d\nName: %s\nType: %s\nParent Domain: %d\nDatabase Server: %d\n\n", i.ID, i.Name, i.Type, i.ParentDomainID, i.DatabaseServerID)
	}
	return nil
}

func DatabaseDeploy(host types.Server, domain types.Domain, db types.Database, file string) error {
	api := factory.GetDatabaseManagement(host.GetServerAuth())

	var dbu *types.DatabaseUser
	for _, u := range domain.DatabaseUsers {
		if u.DatabaseID == 0 || u.DatabaseID == db.ID {
			dbu = &u
		}
	}
	if dbu == nil {
		u := &types.NewDatabaseUser{
			Login:    utils.GenUsername(16),
			Password: utils.GenPassword(24),
		}

		apiU := factory.GetDatabaseUserManagement(host.GetServerAuth())
		newU, err := apiU.CreateDatabaseUser(db, *u)
		if err != nil {
			return err
		}

		dbu = &types.DatabaseUser{
			ID:         newU.ID,
			Login:      newU.Login,
			Password:   u.Password,
			DatabaseID: newU.DatabaseID,
		}

		domain.DatabaseUsers = append(domain.DatabaseUsers, *dbu)
		config.SetDomain(host, domain)
	}

	var s *string
	apiD := factory.GetDomainManagement(host.GetServerAuth())
	i, err := apiD.GetDomain(domain.Name)
	if err == nil {
		s = &i.Sysuser
	}

	return api.DeployDatabase(db, *dbu, *host.GetDatabaseServer(db.DatabaseServerID), file, host.Info.IsWindows, s)
}

func DatabaseDelete(host types.Server, dbn string) error {
	api := factory.GetDatabaseManagement(host.GetServerAuth())

	db, err := config.GetDatabase(host, dbn)
	if err != nil {
		db, err = DatabaseFindNonLocal(api, host, dbn)
		if err != nil {
			return err
		}
	}

	config.DeleteDatabase(host, dbn)

	return api.RemoveDatabase(*db)
}
