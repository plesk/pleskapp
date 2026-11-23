// Copyright 1999-2026. WebPros International GmbH.

package json

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/plesk/pleskapp/plesk/internal/locales"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const settingsPath = "/api/v2/cli/settings/call"

func TestCallCliGivesTheReplyOfTheUtility(t *testing.T) {
	s, c := startServer(t, map[string]reply{
		settingsPath: jsonReply(http.StatusOK, `{"code":0,"stdout":"min_password_strength=strong\n","stderr":"notice\n"}`),
	})

	r, err := callCli(c, "settings", cliGateRequest{
		Params: []string{"--get", "min_password_strength"},
		Env:    map[string]string{"PSA_PASSWORD": "the-password"},
	})

	require.NoError(t, err)
	assert.Equal(t, cliGateResponce{Code: 0, Stdout: "min_password_strength=strong\n", Stderr: "notice\n"}, r)
	checkCliGateRequest(t, s, settingsPath, cliGateCall{
		Params: []string{"--get", "min_password_strength"},
		Env:    map[string]string{"PSA_PASSWORD": "the-password"},
	})
}

func TestCallCliGivesTheErrorOfTheConnection(t *testing.T) {
	s, c := startServer(t, map[string]reply{})
	s.Close()

	r, err := callCli(c, "settings", cliGateRequest{Params: []string{"--list"}, Env: map[string]string{}})

	var urlErr *url.Error
	assert.ErrorAs(t, err, &urlErr)
	assert.Equal(t, cliGateResponce{}, r)
}

func TestCallCliErrors(t *testing.T) {
	checkErrors(t, settingsPath, func(c *resty.Client) error {
		_, err := callCli(c, "settings", cliGateRequest{Params: []string{"--list"}, Env: map[string]string{}})
		return err
	}, errorCases{
		needReauth:    true,
		statusError:   locales.L.Get("api.errors.failed.request", errorReplyError.Code, errorReplyError.Message, errorReplyError.Errors),
		cliGateFailed: locales.L.Get("api.errors.cligate.error.responce", cliGateFailure.Code, cliGateFailure.Stdout, cliGateFailure.Stderr),
	})
}
