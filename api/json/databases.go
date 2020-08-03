// Copyright 1999-2020. Plesk International GmbH.

package json

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/plesk/pleskapp/api"
	"github.com/plesk/pleskapp/types"
)

type jsonDatabases struct {
	auth   api.Auth
	client *http.Client
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

func NewDatabases(a api.Auth) jsonDatabases {
	return jsonDatabases{
		auth:   a,
		client: getClient(a.GetIgnoreSsl()),
	}
}

func (j jsonDatabases) ListDatabases() ([]api.DatabaseInfo, error) {
	var req, _ = http.NewRequest("GET", api.GetApiUrl(j.auth, "/api/v2/databases"), bytes.NewBuffer([]byte{}))
	addBasicHeaders(req, j.auth.GetApiKey())

	var d []databaseInfo
	res, err := doAndThenCheckAuthFailure(j.client, req, j.auth.GetAddress())
	if err != nil {
		return jsonDatabaseInfoToInfo(d), err
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return jsonDatabaseInfoToInfo(d), err
	}

	err = json.Unmarshal(data, &d)
	if err != nil {
		return jsonDatabaseInfoToInfo(d), err
	}

	return jsonDatabaseInfoToInfo(d), err
}

func (j jsonDatabases) ListDomainDatabases(domain string) ([]api.DatabaseInfo, error) {
	var req, _ = http.NewRequest("GET", api.GetApiUrl(j.auth, "/api/v2/databases?domain="+domain), bytes.NewBuffer([]byte{}))
	addBasicHeaders(req, j.auth.GetApiKey())

	var d []databaseInfo
	res, err := doAndThenCheckAuthFailure(j.client, req, j.auth.GetAddress())
	if err != nil {
		return jsonDatabaseInfoToInfo(d), err
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return jsonDatabaseInfoToInfo(d), err
	}

	err = json.Unmarshal(data, &d)
	if err != nil {
		return jsonDatabaseInfoToInfo(d), err
	}

	return jsonDatabaseInfoToInfo(d), err
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
	jd, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", api.GetApiUrl(j.auth, "/api/v2/databases/"), bytes.NewBuffer(jd))
	if err != nil {
		return nil, err
	}

	addBasicHeaders(req, j.auth.GetApiKey())

	res, err := doAndThenCheckAuthFailure(j.client, req, j.auth.GetAddress())
	if err != nil {
		return nil, err
	}

	d, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var s databaseInfo
	e, err := tryParseResponceOrParseError(d, &s)
	if err != nil {
		return nil, err
	}

	if e != nil {
		if e.Code != 0 || len(e.Errors) != 0 {
			return nil, jsonErrorToError(*e)
		}
	}

	return &api.DatabaseInfo{
		ID:               s.ID,
		Name:             s.Name,
		Type:             s.Type,
		ParentDomainID:   s.ParentDomainID,
		DatabaseServerID: s.ServerID,
	}, nil
}

func (j jsonDatabases) RemoveDatabase(db types.Database) error {
	req, err := http.NewRequest("DELETE", api.GetApiUrl(j.auth, "/api/v2/databases/"+strconv.Itoa(db.ID)), bytes.NewBuffer([]byte{}))
	if err != nil {
		return err
	}

	addBasicHeaders(req, j.auth.GetApiKey())

	res, err := doAndThenCheckAuthFailure(j.client, req, j.auth.GetAddress())
	if err != nil {
		return err
	}

	d, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	var s statusResponse
	e, err := tryParseResponceOrParseError(d, &s)
	if err != nil {
		return err
	}

	if e != nil {
		if e.Code != 0 || len(e.Errors) != 0 {
			return jsonErrorToError(*e)
		}
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
	jd, err := json.Marshal(&p)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", api.GetApiUrl(j.auth, "/api/v2/cli/dbbackup/call"), bytes.NewBuffer(jd))
	if err != nil {
		return err
	}
	addBasicHeaders(req, j.auth.GetApiKey())

	res, err := doAndThenCheckAuthFailure(j.client, req, j.auth.GetAddress())
	if err != nil {
		return err
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	var r cliGateResponce

	e, err := tryParseResponceOrParseError(data, &r)
	if err != nil {
		return err
	}

	if e != nil {
		if e.Code != 0 || len(e.Errors) != 0 {
			return jsonErrorToError(*e)
		}
	}

	if r.Code != 0 || len(r.Stderr) != 0 {
		return jsonCliGateResponceToError(r)
	}

	return nil
}

func (j jsonDatabases) ListDatabaseServers() ([]api.DatabaseServerInfo, error) {
	var d []databaseServerInfo
	req, err := http.NewRequest("GET", api.GetApiUrl(j.auth, "/api/v2/dbservers"), bytes.NewBuffer([]byte{}))
	if err != nil {
		return jsonDatabaseServerInfoToInfo(d), err
	}

	addBasicHeaders(req, j.auth.GetApiKey())

	res, err := doAndThenCheckAuthFailure(j.client, req, j.auth.GetAddress())
	if err != nil {
		return jsonDatabaseServerInfoToInfo(d), err
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return jsonDatabaseServerInfoToInfo(d), err
	}

	err = json.Unmarshal(data, &d)
	if err != nil {
		return jsonDatabaseServerInfoToInfo(d), err
	}

	return jsonDatabaseServerInfoToInfo(d), err
}

func (j jsonDatabases) CreateDatabaseUser(db types.Database, dbuser types.NewDatabaseUser) (*api.DatabaseUserInfo, error) {
	p := createDatabaseUserRequest{
		Login:      dbuser.Login,
		Password:   dbuser.Password,
		DatabaseID: db.ID,
	}
	jd, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", api.GetApiUrl(j.auth, "/api/v2/dbusers/"), bytes.NewBuffer(jd))
	if err != nil {
		return nil, err
	}

	addBasicHeaders(req, j.auth.GetApiKey())

	res, err := doAndThenCheckAuthFailure(j.client, req, j.auth.GetAddress())
	if err != nil {
		return nil, err
	}

	d, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var s databaseUserInfo
	e, err := tryParseResponceOrParseError(d, &s)
	if err != nil {
		return nil, err
	}

	if e != nil {
		if e.Code != 0 || len(e.Errors) != 0 {
			return nil, jsonErrorToError(*e)
		}
	}

	return &api.DatabaseUserInfo{
		ID:         s.ID,
		Login:      s.Login,
		DatabaseID: s.DatabaseID,
	}, nil
}

func (j jsonDatabases) RemoveDatabaseUser(dbu types.DatabaseUser) error {
	req, err := http.NewRequest("DELETE", api.GetApiUrl(j.auth, "/api/v2/dbusers/"+strconv.Itoa(dbu.ID)), bytes.NewBuffer([]byte{}))
	if err != nil {
		return err
	}

	addBasicHeaders(req, j.auth.GetApiKey())

	res, err := doAndThenCheckAuthFailure(j.client, req, j.auth.GetAddress())
	if err != nil {
		return err
	}

	d, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	var s statusResponse
	e, err := tryParseResponceOrParseError(d, &s)
	if err != nil {
		return err
	}

	if e != nil {
		if e.Code != 0 || len(e.Errors) != 0 {
			return jsonErrorToError(*e)
		}
	}

	return nil
}

func (j jsonDatabases) ListDatabaseUsers(db types.Database) ([]api.DatabaseUserInfo, error) {
	var d []databaseUserInfo

	req, err := http.NewRequest("GET", api.GetApiUrl(j.auth, "/api/v2/dbusers?dbId="+strconv.Itoa(db.ID)), bytes.NewBuffer([]byte{}))
	if err != nil {
		return jsonDatabaseUserInfoToInfo(d), err
	}

	addBasicHeaders(req, j.auth.GetApiKey())

	res, err := doAndThenCheckAuthFailure(j.client, req, j.auth.GetAddress())
	if err != nil {
		return jsonDatabaseUserInfoToInfo(d), err
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return jsonDatabaseUserInfoToInfo(d), err
	}

	err = json.Unmarshal(data, &d)
	if err != nil {
		return jsonDatabaseUserInfoToInfo(d), err
	}

	return jsonDatabaseUserInfoToInfo(d), err
}
