// Copyright 1999-2020. Plesk International GmbH.

package json

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httputil"
	"reflect"

	"github.com/plesk/pleskapp/locales"
	"github.com/plesk/pleskapp/utils"
)

func getClient(ignoreSsl bool) *http.Client {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.TLSClientConfig = &tls.Config{
		InsecureSkipVerify: ignoreSsl,
	}

	client := &http.Client{
		Transport: transport,
	}

	return client
}

func addBasicHeaders(request *http.Request, apiKey string) {
	request.Header["X-API-Key"] = []string{apiKey}
	request.Header["Content-Type"] = []string{"application/json"}
	request.Header["Accept"] = []string{"application/json"}
}

func tryParseResponceOrParseError(data []byte, something interface{}) (*jsonError, error) {
	var err = json.Unmarshal(data, something)
	if err != nil {
		return nil, err
	}

	if reflect.ValueOf(something).Elem().IsZero() {
		var e jsonError
		_ = json.Unmarshal(data, &e)
		return &e, nil
	}

	return nil, nil
}

func doAndThenCheckAuthFailure(c *http.Client, req *http.Request, a string, flags ...bool) (*http.Response, error) {
	if utils.Log.HasDebug() {
		data, _ := httputil.DumpRequestOut(req, true)
		utils.Log.Debug(fmt.Sprintf("Sending a request %s", data))
	}

	res, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode == 403 || res.StatusCode == 401 {
		if len(flags) > 0 && flags[0] {
			err = errors.New(locales.L.Get("api.errors.auth.wrong.pass", a))
		} else {
			err = errors.New(locales.L.Get("api.errors.auth.failed.reauth", a))
		}
	}

	if utils.Log.HasDebug() {
		data, _ := httputil.DumpResponse(res, true)
		utils.Log.Debug(fmt.Sprintf("Got response %s", data))
	}
	return res, err
}
