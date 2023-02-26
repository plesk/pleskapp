// Copyright 1999-2023. Plesk International GmbH.

package types

import "github.com/plesk/pleskapp/plesk/internal/locales"

type ObjectNotFound struct {
	Object     string
	Filter     string
	Suggestion *string
}

func (e ObjectNotFound) Error() string {
	s := e.Object + " " + e.Filter + " does not exist"
	if e.Suggestion != nil {
		return s + " use " + *e.Suggestion + " to add it"
	}
	return s
}

type ServerNotFound struct {
	Server string
}

func (e ServerNotFound) Error() string {
	return locales.L.Get("errors.server.not.found", e.Server)
}

type DomainNotFound struct {
	Domain string
	Server string
}

func (e DomainNotFound) Error() string {
	return locales.L.Get("errors.domain.not.cached", e.Domain, e.Server, e.Server)
}

type DatabaseNotFound struct {
	DbName string
	Server string
}

func (e DatabaseNotFound) Error() string {
	return locales.L.Get("errors.database.not.found", e.DbName, e.Server)
}

type DatabaseServerNotFound struct {
	DbType string
	Server string
}

func (e DatabaseServerNotFound) Error() string {
	return locales.L.Get("errors.database.server.not.found", e.DbType, e.Server)
}
