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

const domainPath = "/api/v2/cli/domain/call"

func TestGetDomainSysUserGivesTheFtpLogin(t *testing.T) {
	s, c := startServer(t, map[string]reply{domainPath: readReply(t, "cli-domain-info.json")})

	user, err := NewDomains(c).getDomainSysUser("example.com")

	require.NoError(t, err)
	assert.Equal(t, "example_com", user)
	checkCliGateRequest(t, s, domainPath, cliGateCall{
		Params: []string{"--info", "example.com"},
		Env:    map[string]string{},
	})
}

func TestGetDomainSysUserGivesAnErrorWhenTheReplyHasNoFtpLogin(t *testing.T) {
	_, c := startServer(t, map[string]reply{
		domainPath: jsonReply(http.StatusOK, `{"code":0,"stdout":"Hosting Settings\nHosting type: virtual\n","stderr":""}`),
	})

	_, err := NewDomains(c).getDomainSysUser("example.com")

	assert.EqualError(t, err, locales.L.Get("api.errors.domain.info.not.found"))
}

func TestGetDomainSysUserErrors(t *testing.T) {
	checkErrors(t, domainPath, func(c *resty.Client) error {
		_, err := NewDomains(c).getDomainSysUser("example.com")
		return err
	}, errorCases{
		needReauth:    true,
		statusError:   locales.L.Get("api.errors.domain.info.failed", errorReplyError.Code, errorReplyError.Message),
		cliGateFailed: locales.L.Get("api.errors.cligate.error.responce", cliGateFailure.Code, cliGateFailure.Stdout, cliGateFailure.Stderr),
	})
}

func TestAddDomainFeaturesWithOneFeature(t *testing.T) {
	s, c := startServer(t, map[string]reply{
		domainPath: jsonReply(http.StatusOK, `{"code":0,"stdout":"SUCCESS: Update of domain 'example.com' completed.\n","stderr":""}`),
	})

	err := NewDomains(c).AddDomainFeatures("example.com", []string{"php"}, false)

	require.NoError(t, err)
	checkCliGateRequest(t, s, domainPath, cliGateCall{
		Params: []string{"--update", "example.com", "-php", "true", "-php_handler_id", "plesk-php74-fpm"},
		Env:    map[string]string{},
	})
}

func TestAddDomainFeaturesSendsTheParamsOfEachPlatform(t *testing.T) {
	tests := []struct {
		name      string
		feature   string
		isWindows bool
		params    []string
	}{
		{"php73 on Linux", "php73", false, []string{"--update", "example.com", "-php", "true", "-php_handler_id", "plesk-php73-fpm"}},
		{"php73 on Windows", "php73", true, []string{"--update", "example.com", "-php", "true", "-php_handler_id", "fastcgi-7.3"}},
		{"nginx on Linux", "nginx", false, []string{"--update", "example.com", "-nginx-serve-php", "true"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, c := startServer(t, map[string]reply{domainPath: valueReply(t, http.StatusOK, cliGateFailure)})

			_ = NewDomains(c).AddDomainFeatures("example.com", []string{tt.feature}, tt.isWindows)

			checkCliGateRequest(t, s, domainPath, cliGateCall{Params: tt.params, Env: map[string]string{}})
		})
	}
}

func TestAddDomainFeaturesSendsNoRequest(t *testing.T) {
	tests := []struct {
		name      string
		feature   string
		isWindows bool
		wantError string
	}{
		{"unknown feature", "unknown", false, ""},
		{"nginx on Windows", "nginx", true, locales.L.Get("errors.feature.not.supported", "nginx")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, c := startServer(t, map[string]reply{domainPath: valueReply(t, http.StatusOK, cliGateFailure)})

			err := NewDomains(c).AddDomainFeatures("example.com", []string{tt.feature}, tt.isWindows)

			if tt.wantError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.wantError)
			}
			assert.Empty(t, s.Requests())
		})
	}
}

func TestAddDomainFeaturesErrors(t *testing.T) {
	checkErrors(t, domainPath, func(c *resty.Client) error {
		return NewDomains(c).AddDomainFeatures("example.com", []string{"php"}, false)
	}, errorCases{
		needReauth:    true,
		statusError:   locales.L.Get("api.errors.failed.request", errorReplyError.Code, errorReplyError.Message, errorReplyError.Errors),
		cliGateFailed: locales.L.Get("api.errors.cligate.error.responce", cliGateFailure.Code, cliGateFailure.Stdout, cliGateFailure.Stderr),
	})
}
