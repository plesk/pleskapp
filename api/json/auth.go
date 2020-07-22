// Copyright 1999-2020. Plesk International GmbH. All rights reserved.

package json

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"git.plesk.ru/~abashurov/pleskapp/api"
	"git.plesk.ru/~abashurov/pleskapp/locales"
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
	client *http.Client
}

func NewAuth(a api.AuthClient) JsonAuth {
	return JsonAuth{
		client: getClient(a.GetIgnoreSsl()),
	}
}

//GetAPIKey gets API keys via REST CLI Gate
func (j JsonAuth) GetAPIKey(preAuth api.PreAuth) (string, error) {
	// FIXME: Enable direct call when 18.0.25-18.0.28 are no longer used:
	//  https://jira.plesk.ru/browse/PPP-49425

	if false {
		p := createApiKeyRequest{
			IP:          "",
			Login:       "admin",
			Description: "PleskApp API key",
		}
		jd, err := json.Marshal(p)
		if err != nil {
			return "", err
		}

		req, err := http.NewRequest("POST", api.GetApiUrl(preAuth, "/api/v2/auth/keys"), bytes.NewBuffer(jd))
		if err != nil {
			return "", err
		}

		req.Header["Content-Type"] = []string{"application/json"}
		req.Header["Accept"] = []string{"application/json"}
		req.SetBasicAuth(preAuth.GetLogin(), preAuth.GetPassword())

		res, err := j.client.Do(req)
		if err != nil {
			return "", err
		}

		data, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return "", err
		}

		var jRes createApiKeyResponse
		e, err := tryParseResponceOrParseError(data, jRes)
		if err != nil {
			return "", err
		}

		if res.StatusCode == 201 && e == nil {
			return jRes.Key, nil
		}

		return "", fmt.Errorf("Failed to acquire an API key using provided password: [%d: %s]", e.Code, e.Message)
	}

	p := cliGateRequest{
		Params: []string{
			"--create",
			"-description",
			"PleskApp API key",
		},
		Env: map[string]string{},
	}

	jd, _ := json.Marshal(p)

	req, _ := http.NewRequest("POST", api.GetApiUrl(preAuth, "/api/v2/cli/secret_key/call"), bytes.NewBuffer(jd))
	req.Header["Content-Type"] = []string{"application/json"}
	req.Header["Accept"] = []string{"application/json"}
	req.SetBasicAuth(preAuth.GetLogin(), preAuth.GetPassword())

	res, err := doAndThenCheckAuthFailure(j.client, req, preAuth.GetAddress(), true)
	if err != nil {
		return "", err
	}
	var data, _ = ioutil.ReadAll(res.Body)
	var jRes cliGateResponce
	e, err := tryParseResponceOrParseError(data, &jRes)
	if err != nil {
		return "", err
	}

	if jRes.Code == 0 && e == nil {
		return jRes.Stdout, nil
	}

	if jRes.Code != 0 && e == nil {
		return "", errors.New(locales.L.Get("api.errors.auth.failed", jRes.Stderr))
	}

	return "", errors.New(locales.L.Get("api.errors.auth.cli.failed", e.Code, e.Message))
}

func (j JsonAuth) GetLoginLink(auth api.Auth) (string, error) {
	p := cliGateRequest{
		Params: []string{
			"--get-login-link",
		},
		Env: map[string]string{},
	}

	jd, err := json.Marshal(p)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", api.GetApiUrl(auth, "/api/v2/cli/admin/call"), bytes.NewBuffer(jd))
	if err != nil {
		return "", err
	}

	addBasicHeaders(req, auth.GetApiKey())
	res, err := doAndThenCheckAuthFailure(j.client, req, auth.GetAddress())
	if err != nil {
		return "", err
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	var jRes cliGateResponce
	e, err := tryParseResponceOrParseError(data, &jRes)
	if err != nil {
		return "", err
	}

	if e == nil {
		return jRes.Stdout, nil
	}

	return "", fmt.Errorf("Failed to acquire an API key using provided password: [%d: %s]", e.Code, e.Message)
}

type removeAPIKey struct {
	Params []string `json:"params"`
}

//RemoveAPIKey removes API key via REST CLI Gate
func (j JsonAuth) RemoveAPIKey(auth api.Auth) (string, error) {
	var req, _ = http.NewRequest("DELETE", api.GetApiUrl(auth, "/api/v2/auth/keys/"+auth.GetApiKey()), bytes.NewBuffer([]byte{}))
	addBasicHeaders(req, auth.GetApiKey())

	res, err := j.client.Do(req)
	if err != nil {
		return "", err
	}

	var data, _ = ioutil.ReadAll(res.Body)
	return string(data), err
}
