// Copyright 1999-2026. WebPros International GmbH.

package json

import (
	"net/http"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/plesk/pleskapp/plesk/internal/locales"
	"github.com/plesk/pleskapp/plesk/internal/types"
	"github.com/stretchr/testify/assert"
)

const dbbackupPath = "/api/v2/cli/dbbackup/call"

var (
	deployDB       = types.Database{ID: 3, Name: "shop", Type: "mysql", DatabaseServerID: 1}
	deployDBUser   = types.DatabaseUser{ID: 4, Login: "shop_user", Password: "the-db-password", DatabaseID: 3}
	deployDBServer = types.DatabaseServer{ID: 1, Host: "localhost", Port: 3306, Type: "mysql"}
)

func TestDeployDatabaseSendsTheParamsOfEachPlatform(t *testing.T) {
	sysuser := "example_com"
	tests := []struct {
		name      string
		isWindows bool
		sysuser   *string
		want      cliGateCall
	}{
		{
			name:    "Linux with a system user",
			sysuser: &sysuser,
			want: cliGateCall{
				Params: []string{
					"--restore",
					"--server=localhost",
					"--server-type=mysql",
					"--server-login=shop_user",
					"--server-port=3306",
					"--database=shop",
					"--backup-path=/tmp/dump.sql",
					"--sysuser=example_com",
				},
				Env: map[string]string{"PSA_PASSWORD": "the-db-password"},
			},
		},
		{
			name: "Linux with no system user",
			want: cliGateCall{
				Params: []string{
					"--restore",
					"--server=localhost",
					"--server-type=mysql",
					"--server-login=shop_user",
					"--server-port=3306",
					"--database=shop",
					"--backup-path=/tmp/dump.sql",
					"--sysuser=root",
				},
				Env: map[string]string{"PSA_PASSWORD": "the-db-password"},
			},
		},
		{
			name:      "Windows",
			isWindows: true,
			sysuser:   &sysuser,
			want: cliGateCall{
				Params: []string{
					"--restore",
					"-server=localhost",
					"-server-type=mysql",
					"-server-login=shop_user",
					"-server-pwd=the-db-password",
					"-port=3306",
					"-database=shop",
					"-backup-path=/tmp/dump.sql",
				},
				Env: map[string]string{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, c := startServer(t, map[string]reply{dbbackupPath: valueReply(t, http.StatusOK, cliGateFailure)})

			_ = NewDatabases(c).DeployDatabase(deployDB, deployDBUser, deployDBServer, "/tmp/dump.sql", tt.isWindows, tt.sysuser)

			checkCliGateRequest(t, s, dbbackupPath, tt.want)
		})
	}
}

func TestDeployDatabaseGivesNoErrorWhenTheCliGateSucceeds(t *testing.T) {
	_, c := startServer(t, map[string]reply{
		dbbackupPath: jsonReply(http.StatusOK, `{"code":0,"stdout":"Restore is done.\n","stderr":""}`),
	})

	err := NewDatabases(c).DeployDatabase(deployDB, deployDBUser, deployDBServer, "/tmp/dump.sql", false, nil)

	assert.NoError(t, err)
}

func TestDeployDatabaseGivesAnErrorWhenTheCliGateWritesToStderr(t *testing.T) {
	_, c := startServer(t, map[string]reply{
		dbbackupPath: jsonReply(http.StatusOK, `{"code":0,"stdout":"","stderr":"Warning: the dump has errors\n"}`),
	})

	err := NewDatabases(c).DeployDatabase(deployDB, deployDBUser, deployDBServer, "/tmp/dump.sql", false, nil)

	assert.EqualError(t, err, locales.L.Get("api.errors.cligate.error.responce", 0, "", "Warning: the dump has errors\n"))
}

func TestDeployDatabaseErrors(t *testing.T) {
	checkErrors(t, dbbackupPath, func(c *resty.Client) error {
		return NewDatabases(c).DeployDatabase(deployDB, deployDBUser, deployDBServer, "/tmp/dump.sql", false, nil)
	}, errorCases{
		needReauth:    true,
		statusError:   locales.L.Get("api.errors.failed.request", errorReplyError.Code, errorReplyError.Message, errorReplyError.Errors),
		cliGateFailed: locales.L.Get("api.errors.cligate.error.responce", cliGateFailure.Code, cliGateFailure.Stdout, cliGateFailure.Stderr),
	})
}
