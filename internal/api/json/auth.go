// Copyright 1999-2023. Plesk International GmbH.

package json

import (
	"encoding/json"
	"errors"

	"github.com/go-resty/resty/v2"
	"github.com/plesk/pleskapp/plesk/internal/api"
	"github.com/plesk/pleskapp/plesk/internal/locales"
)

type createAPIKeyRequest struct {
	IP          string `json:"ip"`
	Login       string `json:"login"`
	Description string `json:"description"`
}

type createAPIKeyResponse struct {
	Key string `json:"key"`
}

// Auth gets API keys via REST CLI Gate
type Auth struct {
	client *resty.Client
}

func NewAuth(c *resty.Client) Auth {
	return Auth{
		client: c,
	}
}

// GetAPIKey gets API keys via REST CLI Gate
func (j Auth) GetAPIKey(a api.Auth) (string, error) {
	// FIXME: Enable direct call when 18.0.25-18.0.28 are no longer used:
	//  https://jira.plesk.ru/browse/PPP-49425

	if false {
		p := createAPIKeyRequest{
			IP:          "",
			Login:       "admin",
			Description: "PleskApp API key",
		}
		req, err := json.Marshal(p)
		if err != nil {
			return "", err
		}

		res, err := j.client.R().
			SetBody(req).
			SetResult(&createAPIKeyResponse{}).
			SetError(&jsonError{}).
			Post("/api/v2/auth/keys")

		if err != nil {
			return "", err
		}

		if res.IsSuccess() {
			var r = res.Result().(*createAPIKeyResponse)
			return r.Key, nil
		}

		if res.StatusCode() == 403 {
			return "", authError{server: j.client.HostURL, needReauth: false}
		}

		var r = res.Error().(*jsonError)
		return "", errors.New(locales.L.Get("api.errors.auth.cli.failed", r.Code, r.Message))
	}

	p := cliGateRequest{
		Params: []string{
			"--create",
			"-description",
			"PleskApp API key",
		},
		Env: map[string]string{},
	}
	req, err := json.Marshal(p)
	if err != nil {
		return "", err
	}

	res, err := j.client.R().
		SetBody(req).
		SetResult(&cliGateResponce{}).
		SetError(&jsonError{}).
		Post("/api/v2/cli/secret_key/call")

	if err != nil {
		return "", err
	}

	if res.IsSuccess() {
		var r = res.Result().(*cliGateResponce)
		if r.Code == 0 {
			return r.Stdout, nil
		}

		return "", errors.New(locales.L.Get("api.errors.auth.failed", r.Stderr))
	}

	if res.StatusCode() == 403 {
		return "", authError{server: j.client.HostURL, needReauth: false}
	}

	var r = res.Error().(*jsonError)
	return "", jsonErrorToError(*r)
}

func (j Auth) GetLoginLink(auth api.Auth) (string, error) {
	p := cliGateRequest{
		Params: []string{
			"--get-login-link",
		},
		Env: map[string]string{},
	}
	req, err := json.Marshal(p)
	if err != nil {
		return "", err
	}

	res, err := j.client.R().
		SetBody(req).
		SetResult(&cliGateResponce{}).
		SetError(&jsonError{}).
		Post("/api/v2/cli/admin/call")

	if err != nil {
		return "", err
	}

	if res.IsSuccess() {
		var r = res.Result().(*cliGateResponce)
		if r.Code == 0 {
			return r.Stdout, nil
		}

		return "", errors.New(locales.L.Get("api.errors.auth.failed", r.Stderr))
	}

	if res.StatusCode() == 403 {
		return "", authError{server: j.client.HostURL, needReauth: true}
	}

	var r = res.Error().(*jsonError)
	return "", jsonErrorToError(*r)
}

// RemoveAPIKey removes API key via REST CLI Gate
func (j Auth) RemoveAPIKey(auth api.Auth) (string, error) {
	key := auth.GetAPIKey()
	res, err := j.client.R().
		SetError(&jsonError{}).
		Delete("/api/v2/auth/keys/" + *key)

	if err != nil {
		return "", err
	}

	if res.IsError() {
		var r = res.Error().(*jsonError)
		return "", jsonErrorToError(*r)
	}

	return string(res.Body()), nil
}
