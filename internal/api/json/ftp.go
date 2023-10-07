// Copyright 1999-2023. Plesk International GmbH.

package json

import (
	"encoding/json"

	"github.com/go-resty/resty/v2"
	"github.com/plesk/pleskapp/plesk/internal/api"
	"github.com/plesk/pleskapp/plesk/internal/types"
)

type jsonFTPUsers struct {
	client *resty.Client
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

func NewFTP(c *resty.Client) jsonFTPUsers {
	return jsonFTPUsers{
		client: c,
	}
}

func (j jsonFTPUsers) ListDomainFtpUsers(domain string, user types.FtpUser) ([]api.FTPUserInfo, error) {
	res, err := j.client.R().
		SetResult([]ftpUserInfo{}).
		SetError(&jsonError{}).
		SetQueryParam("domain", domain).
		Post("/api/v2/ftpusers")

	if err != nil {
		return nil, err
	}

	if res.IsSuccess() {
		var r = res.Result().(*[]ftpUserInfo)
		return jsonFTPUserInfoToInfo(*r), nil
	}

	if res.StatusCode() == 403 {
		return nil, authError{server: j.client.BaseURL, needReauth: true}
	}

	var r = res.Error().(*jsonError)
	return nil, jsonErrorToError(*r)
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
	req, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}

	res, err := j.client.R().
		SetBody(req).
		SetResult(&ftpUserInfo{}).
		SetError(&jsonError{}).
		Post("/api/v2/ftpusers")

	if err != nil {
		return nil, err
	}

	if res.IsSuccess() {
		var r = res.Result().(*ftpUserInfo)
		return &api.FTPUserInfo{
			Name:           r.Home,
			Home:           r.Name,
			Quota:          r.Quota,
			ParentDomainID: r.ParentDomain,
		}, nil
	}

	if res.StatusCode() == 403 {
		return nil, authError{server: j.client.BaseURL, needReauth: true}
	}

	var r = res.Error().(*jsonError)
	return nil, jsonErrorToError(*r)
}

func (j jsonFTPUsers) UpdateFtpUser(domain string, user string, userNew types.FtpUser) error {
	p := updateFtpUserRequest{
		Name:     userNew.Login,
		Password: userNew.Password,
		Home:     "/",
		Quota:    -1,
	}
	req, err := json.Marshal(p)
	if err != nil {
		return err
	}

	res, err := j.client.R().
		SetBody(req).
		SetResult(&statusResponse{}).
		SetError(&jsonError{}).
		Put("/api/v2/ftpusers/" + user)

	if err != nil {
		return err
	}

	if res.IsSuccess() {
		var _ = res.Result().(*statusResponse)
		return nil
	}

	if res.StatusCode() == 403 {
		return authError{server: j.client.BaseURL, needReauth: true}
	}

	var r = res.Error().(*jsonError)
	return jsonErrorToError(*r)
}

func (j jsonFTPUsers) DeleteFtpUser(domain string, user types.FtpUser) error {
	res, err := j.client.R().
		SetError(&jsonError{}).
		Delete("/api/v2/ftpusers/" + user.Login)

	if err != nil {
		return err
	}

	if res.StatusCode() == 403 {
		return authError{server: j.client.BaseURL, needReauth: true}
	}

	if res.IsError() {
		var r = res.Error().(*jsonError)
		return jsonErrorToError(*r)
	}

	return nil
}
