// Copyright 1999-2023. Plesk International GmbH.

package factory

import (
	"crypto/tls"
	"io"
	"log"
	"net/http"

	"github.com/go-resty/resty/v2"
	"github.com/plesk/pleskapp/plesk/internal/api"
	"github.com/plesk/pleskapp/plesk/internal/api/json"
)

func buildClient(a api.Auth) *resty.Client {
	tr := http.DefaultTransport.(*http.Transport).Clone()
	tr.TLSClientConfig = &tls.Config{
		InsecureSkipVerify: a.GetIgnoreSsl(),
	}

	c := &http.Client{
		Transport: tr,
	}

	r := resty.NewWithClient(c).
		SetBaseURL("https://"+a.GetAddress()+":"+a.GetPort()).
		SetHeader("Content-Type", "application/json").
		SetDebug(log.Writer() != io.Discard)

	login := a.GetLogin()
	pass := a.GetPassword()
	key := a.GetAPIKey()

	if login != nil && pass != nil {
		r.SetBasicAuth(*login, *pass)
	} else {
		r.SetHeader("X-API-Key", *key)
	}

	return r
}

func GetDomainManagement(a api.Auth) api.DomainManagement {
	return json.NewDomains(buildClient(a))
}

func GetFTPUserManagement(a api.Auth) api.FTPManagement {
	return json.NewFTP(buildClient(a))
}

func GetDatabaseManagement(a api.Auth) api.DatabaseManagement {
	return json.NewDatabases(buildClient(a))
}

func GetDatabaseUserManagement(a api.Auth) api.DatabaseUserManagement {
	return json.NewDatabases(buildClient(a))
}

func GetAuthentication(a api.Auth) api.Authentication {
	return json.NewAuth(buildClient(a))
}

func GetServerInfo(a api.Auth) api.Server {
	return json.NewInfo(buildClient(a))
}
