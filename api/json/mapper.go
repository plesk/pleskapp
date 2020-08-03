// Copyright 1999-2020. Plesk International GmbH.

package json

import (
	"errors"

	"github.com/plesk/pleskapp/plesk/api"
	"github.com/plesk/pleskapp/plesk/locales"
)

func jsonDomainInfoToInfo(domains []domainInfo) []api.DomainInfo {
	var ds []api.DomainInfo

	for _, i := range domains {
		ds = append(ds, api.DomainInfo{
			ID:             i.ID,
			Name:           i.Name,
			HostingType:    i.HostingType,
			ParentDomainID: i.BaseDomainID,
			GUID:           i.GUID,
			WWWRoot:        i.WWWRoot,
		})
	}
	return ds
}

func jsonFTPUserInfoToInfo(users []ftpUserInfo) []api.FTPUserInfo {
	var us []api.FTPUserInfo

	for _, i := range users {
		us = append(us, api.FTPUserInfo{
			Name:           i.Home,
			Home:           i.Name,
			Quota:          i.Quota,
			ParentDomainID: i.ParentDomain,
		})
	}
	return us
}

func jsonDatabaseInfoToInfo(dbs []databaseInfo) []api.DatabaseInfo {
	var db []api.DatabaseInfo

	for _, i := range dbs {
		db = append(db, api.DatabaseInfo{
			ID:               i.ID,
			Name:             i.Name,
			Type:             i.Type,
			ParentDomainID:   i.ParentDomainID,
			DatabaseServerID: i.ServerID,
		})
	}
	return db
}

func jsonDatabaseUserInfoToInfo(dbus []databaseUserInfo) []api.DatabaseUserInfo {
	var db []api.DatabaseUserInfo

	for _, i := range dbus {
		db = append(db, api.DatabaseUserInfo{
			ID:         i.ID,
			Login:      i.Login,
			DatabaseID: i.DatabaseID,
		})
	}
	return db
}

func jsonDatabaseServerInfoToInfo(dbss []databaseServerInfo) []api.DatabaseServerInfo {
	var db []api.DatabaseServerInfo

	for _, i := range dbss {
		db = append(db, api.DatabaseServerInfo{
			ID:        i.ID,
			Host:      i.Host,
			Port:      i.Port,
			Type:      i.Type,
			Status:    i.Status,
			IsDefault: i.IsDefault,
			IsLocal:   i.IsLocal,
		})
	}
	return db
}

func jsonCliGateResponceToError(r cliGateResponce) error {
	return errors.New(locales.L.Get("api.errors.cligate.error.responce", r.Code, r.Stdout, r.Stderr))
}

func jsonErrorToError(e jsonError) error {
	return errors.New(locales.L.Get("api.errors.failed.request", e.Code, e.Message, e.Errors))
}
