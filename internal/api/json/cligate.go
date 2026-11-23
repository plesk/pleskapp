// Copyright 1999-2026. WebPros International GmbH.

package json

import (
	"net/http"

	"github.com/go-resty/resty/v2"
	"github.com/plesk/pleskapp/plesk/internal/locales"
)

type cliGateError struct {
	reply cliGateResponce
}

func (e cliGateError) Error() string {
	return locales.L.Get("api.errors.cligate.error.responce", e.reply.Code, e.reply.Stdout, e.reply.Stderr)
}

type requestError struct {
	reply jsonError
}

func (e requestError) Error() string {
	return jsonErrorToError(e.reply).Error()
}

// callCli runs one Plesk utility through the CLI gate of the REST API.
// HTTP 403 gives authError that asks to log in again, another error status gives requestError,
// and an exit code that is not 0 gives cliGateError.
func callCli(c *resty.Client, utility string, req cliGateRequest) (cliGateResponce, error) {
	res, err := c.R().
		SetBody(req).
		SetResult(&cliGateResponce{}).
		SetError(&jsonError{}).
		Post("/api/v2/cli/" + utility + "/call")
	if err != nil {
		return cliGateResponce{}, err
	}

	if res.StatusCode() == http.StatusForbidden {
		return cliGateResponce{}, authError{server: c.BaseURL, needReauth: true}
	}
	if !res.IsSuccess() {
		return cliGateResponce{}, requestError{reply: *res.Error().(*jsonError)}
	}

	r := *res.Result().(*cliGateResponce)
	if r.Code != 0 {
		return r, cliGateError{reply: r}
	}
	return r, nil
}
