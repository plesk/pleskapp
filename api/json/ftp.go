// Copyright 1999-2020. Plesk International GmbH. All rights reserved.

package json

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"git.plesk.ru/projects/SBX/repos/pleskapp/api"
	"git.plesk.ru/projects/SBX/repos/pleskapp/types"
)

type jsonFTPUsers struct {
	auth   api.Auth
	client *http.Client
}

type createFtpUserRequest struct {
	Name         string          `json:"name"`
	Password     string          `json:"password"`
	Home         string          `json:"home"`
	Quota        int             `json:"quota"`
	ParentDomain domainReference `json:"parent_domain"`
	Permissions  permissions     `json:"permissions"`
}

type permissions struct {
	Read  string `json:"read"`
	Write string `json:"write"`
}

type updateFtpUserRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	Home     string `json:"home"`
	Quota    int    `json:"quota"`
}

func NewFTP(a api.Auth) jsonFTPUsers {
	return jsonFTPUsers{
		auth:   a,
		client: getClient(a.GetIgnoreSsl()),
	}
}

func (j jsonFTPUsers) ListDomainFtpUsers(domain string, user types.FtpUser) ([]api.FTPUserInfo, error) {
	var req, _ = http.NewRequest("GET", api.GetApiUrl(j.auth, "/api/v2/ftpusers?domain="+domain), bytes.NewBuffer([]byte{}))
	addBasicHeaders(req, j.auth.GetApiKey())

	var u []ftpUserInfo
	res, err := doAndThenCheckAuthFailure(j.client, req, j.auth.GetAddress())
	if err != nil {
		return jsonFTPUserInfoToInfo(u), err
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return jsonFTPUserInfoToInfo(u), err
	}

	err = json.Unmarshal(data, &u)
	if err != nil {
		return jsonFTPUserInfoToInfo(u), err
	}

	return jsonFTPUserInfoToInfo(u), nil
}

func (j jsonFTPUsers) CreateFtpUser(domain string, user types.FtpUser) (*api.FTPUserInfo, error) {
	p := createFtpUserRequest{
		Name:     user.Login,
		Password: user.Password,
		Home:     "/",
		Quota:    -1,
		Permissions: permissions{
			Write: "true",
			Read:  "true",
		},
		ParentDomain: domainReference{
			Name: domain,
		},
	}
	jd, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", api.GetApiUrl(j.auth, "/api/v2/ftpusers/"), bytes.NewBuffer(jd))
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

	var s ftpUserInfo
	e, err := tryParseResponceOrParseError(d, &s)
	if err != nil {
		return nil, err
	}

	if e != nil {
		if e.Code != 0 || len(e.Errors) != 0 {
			return nil, jsonErrorToError(*e)
		}
	}

	return &api.FTPUserInfo{
		Name:           s.Home,
		Home:           s.Name,
		Quota:          s.Quota,
		ParentDomainID: s.ParentDomain,
	}, nil
}

func (j jsonFTPUsers) UpdateFtpUser(domain string, user string, userNew types.FtpUser) error {
	p := updateFtpUserRequest{
		Name:     userNew.Login,
		Password: userNew.Password,
		Home:     "/",
		Quota:    -1,
	}
	jd, err := json.Marshal(p)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("UPDATE", api.GetApiUrl(j.auth, "/api/v2/ftpusers/"+user), bytes.NewBuffer(jd))
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

func (j jsonFTPUsers) DeleteFtpUser(domain string, user types.FtpUser) error {
	req, err := http.NewRequest("DELETE", api.GetApiUrl(j.auth, "/api/v2/ftpusers/"+user.Login), bytes.NewBuffer([]byte{}))
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
