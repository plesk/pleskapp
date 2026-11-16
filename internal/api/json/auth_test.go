// Copyright 1999-2026. WebPros International GmbH.

package json

import (
	"net/http"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/plesk/pleskapp/plesk/internal/locales"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	secretKeyPath = "/api/v2/cli/secret_key/call"
	adminPath     = "/api/v2/cli/admin/call"
)

func TestGetAPIKeyGivesTheStdoutOfTheCliGate(t *testing.T) {
	s, c := startServer(t, map[string]reply{
		secretKeyPath: jsonReply(http.StatusOK, `{"code":0,"stdout":"the-new-key","stderr":""}`),
	})

	key, err := NewAuth(c).GetAPIKey(nil)

	require.NoError(t, err)
	assert.Equal(t, "the-new-key", key)
	checkCliGateRequest(t, s, secretKeyPath, cliGateCall{
		Params: []string{"--create", "-description", "PleskApp API key"},
		Env:    map[string]string{},
	})
}

func TestGetAPIKeyErrors(t *testing.T) {
	checkErrors(t, secretKeyPath, func(c *resty.Client) error {
		_, err := NewAuth(c).GetAPIKey(nil)
		return err
	}, errorCases{
		needReauth:    false,
		statusError:   locales.L.Get("api.errors.failed.request", errorReplyError.Code, errorReplyError.Message, errorReplyError.Errors),
		cliGateFailed: locales.L.Get("api.errors.auth.failed", cliGateFailure.Stderr),
	})
}

func TestGetLoginLinkGivesTheStdoutOfTheCliGate(t *testing.T) {
	link := "https://plesk.example.com:8443/login?secret=abc\n"
	s, c := startServer(t, map[string]reply{
		adminPath: jsonReply(http.StatusOK, `{"code":0,"stdout":"https://plesk.example.com:8443/login?secret=abc\n","stderr":""}`),
	})

	got, err := NewAuth(c).GetLoginLink(nil)

	require.NoError(t, err)
	assert.Equal(t, link, got)
	checkCliGateRequest(t, s, adminPath, cliGateCall{
		Params: []string{"--get-login-link"},
		Env:    map[string]string{},
	})
}

func TestGetLoginLinkErrors(t *testing.T) {
	checkErrors(t, adminPath, func(c *resty.Client) error {
		_, err := NewAuth(c).GetLoginLink(nil)
		return err
	}, errorCases{
		needReauth:    true,
		statusError:   locales.L.Get("api.errors.failed.request", errorReplyError.Code, errorReplyError.Message, errorReplyError.Errors),
		cliGateFailed: locales.L.Get("api.errors.auth.failed", cliGateFailure.Stderr),
	})
}
