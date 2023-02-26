// Copyright 1999-2023. Plesk International GmbH.

package json

import (
	"encoding/json"
	"errors"
	"strconv"

	"github.com/go-resty/resty/v2"
	"github.com/plesk/pleskapp/plesk/internal/api"
	"github.com/plesk/pleskapp/plesk/internal/locales"
	"github.com/plesk/pleskapp/plesk/internal/types"
)

type jsonDatabases struct {
	client *resty.Client
}

type createDatabaseRequest struct {
	Name         string          `json:"name"`
	Type         string          `json:"type"`
	ParentDomain domainReference `json:"parent_domain"`
	ServerID     int             `json:"server_id"`
}

type createDatabaseUserRequest struct {
	Login      string `json:"login"`
	Password   string `json:"password"`
	DatabaseID int    `json:"database_id"`
}

func NewDatabases(c *resty.Client) jsonDatabases {
	return jsonDatabases{
		client: c,
	}
}

func (j jsonDatabases) ListDatabases() ([]api.DatabaseInfo, error) {
	res, err := j.client.R().
		SetResult([]databaseInfo{}).
		SetError(&jsonError{}).
		Get("/api/v2/databases")

	if err != nil {
		return jsonDatabaseInfoToInfo([]databaseInfo{}), err
	}

	if res.IsSuccess() {
		var r *[]databaseInfo = res.Result().(*[]databaseInfo)
		return jsonDatabaseInfoToInfo(*r), err
	}

	if res.StatusCode() == 403 {
		return jsonDatabaseInfoToInfo([]databaseInfo{}), authError{server: j.client.HostURL, needReauth: true}
	}

	var r *jsonError = res.Error().(*jsonError)
	return jsonDatabaseInfoToInfo([]databaseInfo{}), errors.New(locales.L.Get("api.errors.failed.request", r.Code, r.Message, r.Errors))
}

func (j jsonDatabases) ListDomainDatabases(domain string) ([]api.DatabaseInfo, error) {
	res, err := j.client.R().
		SetResult([]databaseInfo{}).
		SetError(&jsonError{}).
		SetQueryParam("domain", domain).
		Get("/api/v2/databases")

	if err != nil {
		return jsonDatabaseInfoToInfo([]databaseInfo{}), err
	}

	if res.IsSuccess() {
		var r *[]databaseInfo = res.Result().(*[]databaseInfo)
		return jsonDatabaseInfoToInfo(*r), err
	}

	if res.StatusCode() == 403 {
		return jsonDatabaseInfoToInfo([]databaseInfo{}), authError{server: j.client.HostURL, needReauth: true}
	}

	var r *jsonError = res.Error().(*jsonError)
	return jsonDatabaseInfoToInfo([]databaseInfo{}), jsonErrorToError(*r)
}

func (j jsonDatabases) CreateDatabase(domain types.Domain, db types.NewDatabase, dbs types.DatabaseServer) (*api.DatabaseInfo, error) {
	p := createDatabaseRequest{
		Name: db.Name,
		Type: db.Type,
		ParentDomain: domainReference{
			Name: domain.Name,
		},
		ServerID: dbs.ID,
	}
	req, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}

	res, err := j.client.R().
		SetBody(req).
		SetResult(&databaseInfo{}).
		SetError(&jsonError{}).
		Post("/api/v2/databases")

	if err != nil {
		return nil, err
	}

	if res.IsSuccess() {
		var r *databaseInfo = res.Result().(*databaseInfo)
		return &api.DatabaseInfo{
			ID:               r.ID,
			Name:             r.Name,
			Type:             r.Type,
			ParentDomainID:   r.ParentDomainID,
			DatabaseServerID: r.ServerID,
		}, nil
	}

	if res.StatusCode() == 403 {
		return nil, authError{server: j.client.HostURL, needReauth: true}
	}

	var r *jsonError = res.Error().(*jsonError)
	return nil, jsonErrorToError(*r)
}

func (j jsonDatabases) RemoveDatabase(db types.Database) error {
	res, err := j.client.R().
		SetError(&jsonError{}).
		Post("/api/v2/databases/" + strconv.Itoa(db.ID))

	if err != nil {
		return err
	}

	if res.StatusCode() == 403 {
		return authError{server: j.client.HostURL, needReauth: true}
	}

	if res.IsError() {
		var r *jsonError = res.Error().(*jsonError)
		return jsonErrorToError(*r)
	}

	return nil
}

func (j jsonDatabases) DeployDatabase(
	db types.Database,
	dbu types.DatabaseUser,
	dbs types.DatabaseServer,
	file string,
	isWindows bool,
	sysuser *string,
) error {
	var p cliGateRequest
	if !isWindows {
		s := "root"
		if sysuser != nil {
			s = *sysuser
		}
		p = cliGateRequest{
			Params: []string{
				"--restore",
				"--server=" + dbs.Host,
				"--server-type=" + dbs.Type,
				"--server-login=" + dbu.Login,
				"--server-port=" + strconv.Itoa(dbs.Port),
				"--database=" + db.Name,
				"--backup-path=" + file,
				"--sysuser=" + s,
			},
			Env: map[string]string{
				"PSA_PASSWORD": dbu.Password,
			},
		}
	} else {
		// FIXME: REST API seems to be unable to run dbbackup via pm_ApiCli::callSbin on Windows
		p = cliGateRequest{
			Params: []string{
				"--restore",
				"-server=" + dbs.Host,
				"-server-type=" + dbs.Type,
				"-server-login=" + dbu.Login,
				"-server-pwd=" + dbu.Password,
				"-port=" + strconv.Itoa(dbs.Port),
				"-database=" + db.Name,
				"-backup-path=" + file,
			},
			Env: map[string]string{},
		}
	}

	req, err := json.Marshal(p)
	if err != nil {
		return err
	}

	res, err := j.client.R().
		SetBody(req).
		SetResult(&cliGateResponce{}).
		SetError(&jsonError{}).
		Post("/api/v2/cli/dbbackup/call")

	if err != nil {
		return err
	}

	if res.IsSuccess() {
		var r *cliGateResponce = res.Result().(*cliGateResponce)
		if r.Code != 0 || len(r.Stderr) != 0 {
			return jsonCliGateResponceToError(*r)
		}
	}

	if res.StatusCode() == 403 {
		return authError{server: j.client.HostURL, needReauth: true}
	}

	var r *jsonError = res.Error().(*jsonError)
	return jsonErrorToError(*r)
}

func (j jsonDatabases) ListDatabaseServers() ([]api.DatabaseServerInfo, error) {
	res, err := j.client.R().
		SetResult([]databaseServerInfo{}).
		SetError(&jsonError{}).
		Get("/api/v2/dbservers")

	if err != nil {
		return []api.DatabaseServerInfo{}, err
	}

	if res.IsSuccess() {
		var r *[]databaseServerInfo = res.Result().(*[]databaseServerInfo)
		return jsonDatabaseServerInfoToInfo(*r), err
	}

	if res.StatusCode() == 403 {
		return []api.DatabaseServerInfo{}, authError{server: j.client.HostURL, needReauth: true}
	}

	var r *jsonError = res.Error().(*jsonError)
	return []api.DatabaseServerInfo{}, jsonErrorToError(*r)
}

func (j jsonDatabases) CreateDatabaseUser(db types.Database, dbuser types.NewDatabaseUser) (*api.DatabaseUserInfo, error) {
	p := createDatabaseUserRequest{
		Login:      dbuser.Login,
		Password:   dbuser.Password,
		DatabaseID: db.ID,
	}
	req, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}

	res, err := j.client.R().
		SetBody(req).
		SetResult(&databaseUserInfo{}).
		SetError(&jsonError{}).
		Post("/api/v2/dbusers")

	if err != nil {
		return nil, err
	}

	if res.IsSuccess() {
		var r *databaseUserInfo = res.Result().(*databaseUserInfo)
		return &api.DatabaseUserInfo{
			ID:         r.ID,
			Login:      r.Login,
			DatabaseID: r.DatabaseID,
		}, nil
	}

	if res.StatusCode() == 403 {
		return nil, authError{server: j.client.HostURL, needReauth: true}
	}

	var r *jsonError = res.Error().(*jsonError)
	return nil, jsonErrorToError(*r)
}

func (j jsonDatabases) RemoveDatabaseUser(dbu types.DatabaseUser) error {
	res, err := j.client.R().
		SetError(&jsonError{}).
		Delete("/api/v2/dbusers/" + strconv.Itoa(dbu.ID))

	if err != nil {
		return err
	}

	if res.StatusCode() == 403 {
		return authError{server: j.client.HostURL, needReauth: true}
	}

	if res.IsError() {
		var r *jsonError = res.Error().(*jsonError)
		return jsonErrorToError(*r)
	}

	return nil
}

func (j jsonDatabases) ListDatabaseUsers(db types.Database) ([]api.DatabaseUserInfo, error) {
	res, err := j.client.R().
		SetResult([]databaseUserInfo{}).
		SetError(&jsonError{}).
		SetQueryParam("dbId", strconv.Itoa(db.ID)).
		Get("/api/v2/dbusers")

	if err != nil {
		return []api.DatabaseUserInfo{}, err
	}

	if res.IsSuccess() {
		var r *[]databaseUserInfo = res.Result().(*[]databaseUserInfo)
		return jsonDatabaseUserInfoToInfo(*r), err
	}

	if res.StatusCode() == 403 {
		return []api.DatabaseUserInfo{}, authError{server: j.client.HostURL, needReauth: true}
	}

	var r *jsonError = res.Error().(*jsonError)
	return []api.DatabaseUserInfo{}, jsonErrorToError(*r)
}
