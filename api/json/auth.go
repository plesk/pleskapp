// Copyright 1999-2020. Plesk International GmbH.

package json

import (
	"encoding/json"
	"errors"

	"github.com/go-resty/resty/v2"
	"github.com/plesk/pleskapp/plesk/api"
	"github.com/plesk/pleskapp/plesk/locales"
)

type createApiKeyRequest struct {
	IP          string `json:"ip"`
	Login       string `json:"login"`
	Description string `json:"description"`
}

type createApiKeyResponse struct {
	Key string `json:"key"`
}

//JsonAuth gets API keys via REST CLI Gate
type JsonAuth struct {
	client *resty.Client
}

func NewAuth(c *resty.Client) JsonAuth {
	return JsonAuth{
		client: c,
	}
}

//GetAPIKey gets API keys via REST CLI Gate
func (j JsonAuth) GetAPIKey(a api.Auth) (string, error) {
	// FIXME: Enable direct call when 18.0.25-18.0.28 are no longer used:
	//  https://jira.plesk.ru/browse/PPP-49425

	if false {
		p := createApiKeyRequest{
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
			SetResult(&createApiKeyResponse{}).
			SetError(&jsonError{}).
			Post("/api/v2/auth/keys")

		if err != nil {
			return "", err
		}

		if res.IsSuccess() {
			var r *createApiKeyResponse = res.Result().(*createApiKeyResponse)
			return r.Key, nil
		}

		if res.StatusCode() == 403 {
			return "", authError{server: j.client.HostURL, needReauth: false}
		}

		var r *jsonError = res.Error().(*jsonError)
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
		var r *cliGateResponce = res.Result().(*cliGateResponce)
		if r.Code == 0 {
			return r.Stdout, nil
		}

		return "", errors.New(locales.L.Get("api.errors.auth.failed", r.Stderr))
	}

	if res.StatusCode() == 403 {
		return "", authError{server: j.client.HostURL, needReauth: false}
	}

	var r *jsonError = res.Error().(*jsonError)
	return "", jsonErrorToError(*r)
}

func (j JsonAuth) GetLoginLink(auth api.Auth) (string, error) {
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
		var r *cliGateResponce = res.Result().(*cliGateResponce)
		if r.Code == 0 {
			return r.Stdout, nil
		}

		return "", errors.New(locales.L.Get("api.errors.auth.failed", r.Stderr))
	}

	if res.StatusCode() == 403 {
		return "", authError{server: j.client.HostURL, needReauth: true}
	}

	var r *jsonError = res.Error().(*jsonError)
	return "", jsonErrorToError(*r)
}

type removeAPIKey struct {
	Params []string `json:"params"`
}

//RemoveAPIKey removes API key via REST CLI Gate
func (j JsonAuth) RemoveAPIKey(auth api.Auth) (string, error) {
	key := auth.GetApiKey()
	res, err := j.client.R().
		SetError(&jsonError{}).
		Delete("/api/v2/auth/keys/" + *key)

	if err != nil {
		return "", err
	}

	if res.IsError() {
		var r *jsonError = res.Error().(*jsonError)
		return "", jsonErrorToError(*r)
	}

	return string(res.Body()), nil
}
