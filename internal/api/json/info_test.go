// Copyright 1999-2026. WebPros International GmbH.

package json

import (
	"net/http"
	"testing"

	"github.com/plesk/pleskapp/plesk/internal/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const serverPath = "/api/v2/server"

func TestGetInfoGivesTheServerInfo(t *testing.T) {
	s, c := startServer(t, map[string]reply{serverPath: readReply(t, "server-info.json")})

	info, err := NewInfo(c).GetInfo()

	require.NoError(t, err)
	assert.Equal(t, api.ServerInfo{
		IsWindows: false,
		Version:   "18.0.81.0",
		Revision:  "c90b57d24b478e9f350eff5f4bd2a3378656e2cd",
		BuildDate: "2026-09-08",
	}, info)
	r := s.onlyRequest(t)
	assert.Equal(t, http.MethodGet, r.Method)
	assert.Equal(t, serverPath, r.Path)
}

func TestGetInfoGivesTheWindowsPlatform(t *testing.T) {
	_, c := startServer(t, map[string]reply{
		serverPath: jsonReply(http.StatusOK, `{"platform":"Windows","panel_version":"18.0.81","panel_update_version":"2"}`),
	})

	info, err := NewInfo(c).GetInfo()

	require.NoError(t, err)
	assert.True(t, info.IsWindows)
	assert.Equal(t, "18.0.81.2", info.Version)
}

func TestGetInfoGivesAnAuthErrorForHTTP403(t *testing.T) {
	_, c := startServer(t, map[string]reply{serverPath: readReply(t, "forbidden.json")})

	_, err := NewInfo(c).GetInfo()

	assert.Equal(t, authError{server: c.BaseURL, needReauth: true}, err)
}
