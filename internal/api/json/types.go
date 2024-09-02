// Copyright 1999-2024. Plesk International GmbH.

package json

import (
	"github.com/plesk/pleskapp/plesk/internal/locales"
	"net/url"
)

type jsonError struct {
	Code    int              `json:"code"`
	Message string           `json:"message"`
	Errors  []jsonErrorInner `json:"errors"`
}

type jsonErrorInner struct {
	Property string `json:"property"`
	Message  string `json:"message"`
}

type cliGateRequest struct {
	Params []string          `json:"params"`
	Env    map[string]string `json:"env"`
}

type cliGateResponce struct {
	Code   int    `json:"code"`
	Stdout string `json:"stdout"`
	Stderr string `json:"stderr"`
}

type statusResponse struct {
	Status string `json:"string"`
}

type domainReference struct {
	Name string `json:"name"`
}

type domainInfo struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	ASCIIName    string `json:"ascii_name"`
	HostingType  string `json:"hosting_type"`
	BaseDomainID int    `json:"base_domain_id"`
	GUID         string `json:"guid"`
	Created      string `json:"created"`
	WWWRoot      string `json:"www_root"`
}

type ftpUserInfo struct {
	Name        string `json:"name"`
	Home        string `json:"home"`
	Quota       int    `json:"quota"`
	Permissions struct {
		Write string `json:"write"`
		Read  string `json:"read"`
	} `json:"permissions"`
	ParentDomain int `json:"parent_domain"`
}

type databaseUserInfo struct {
	ID         int    `json:"id"`
	Login      string `json:"login"`
	DatabaseID int    `json:"database_id"`
}

type databaseServerInfo struct {
	ID        int    `json:"id"`
	Host      string `json:"host"`
	Port      int    `json:"port"`
	Type      string `json:"type"`
	Status    string `json:"status"`
	DBCount   int    `json:"db_count"`
	IsDefault bool   `json:"is_default"`
	IsLocal   bool   `json:"is_local"`
}

type databaseInfo struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	Type           string `json:"type"`
	ParentDomainID int    `json:"parent_domain"`
	ServerID       int    `json:"server_id"`
	DefaultUserID  int    `json:"default_user_id"`
}

type authError struct {
	server     string
	needReauth bool
}

func (e authError) Error() string {
	serverUrl, err := url.Parse(e.server)
	host := e.server
	if err == nil {
		host = serverUrl.Hostname()
	}

	if e.needReauth {
		return locales.L.Get("api.errors.auth.failed.reauth", host)
	}
	return locales.L.Get("api.errors.auth.wrong.pass", host)
}
