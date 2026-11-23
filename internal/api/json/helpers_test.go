// Copyright 1999-2026. WebPros International GmbH.

package json

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type reply struct {
	Status      int             `json:"status"`
	ContentType string          `json:"content_type"`
	Body        json.RawMessage `json:"body"`
}

func jsonReply(status int, body string) reply {
	return reply{Status: status, Body: json.RawMessage(body)}
}

func valueReply(t *testing.T, status int, v any) reply {
	t.Helper()
	body, err := json.Marshal(v)
	require.NoError(t, err)
	return reply{Status: status, Body: body}
}

func readReply(t *testing.T, name string) reply {
	t.Helper()
	path := filepath.Join("testdata", name)
	text, err := os.ReadFile(path)
	require.NoError(t, err, "read %s", path)

	var r reply
	require.NoError(t, json.Unmarshal(text, &r), "decode %s", path)
	return r
}

type recordedRequest struct {
	Method string
	Path   string
	Query  url.Values
	Header http.Header
	Body   []byte
}

type testServer struct {
	*httptest.Server

	mu       sync.Mutex
	requests []recordedRequest
}

func startServer(t *testing.T, replies map[string]reply) (*testServer, *resty.Client) {
	t.Helper()
	s := &testServer{}
	s.Server = httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		s.mu.Lock()
		s.requests = append(s.requests, recordedRequest{
			Method: r.Method,
			Path:   r.URL.Path,
			Query:  r.URL.Query(),
			Header: r.Header.Clone(),
			Body:   body,
		})
		s.mu.Unlock()

		answer, ok := replies[r.URL.Path]
		if !ok {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusNotFound)
			_, _ = fmt.Fprintf(w, "the test server has no reply for %q\n", r.URL.Path)
			return
		}

		status := answer.Status
		if status == 0 {
			status = http.StatusOK
		}
		contentType := answer.ContentType
		if contentType == "" {
			contentType = "application/json"
		}
		w.Header().Set("Content-Type", contentType)
		w.WriteHeader(status)
		_, _ = w.Write(answer.Body)
	}))
	t.Cleanup(s.Close)

	c := resty.NewWithClient(s.Client()).
		SetBaseURL(s.URL).
		SetHeader("Content-Type", "application/json")

	return s, c
}

func (s *testServer) Requests() []recordedRequest {
	s.mu.Lock()
	defer s.mu.Unlock()
	return append([]recordedRequest(nil), s.requests...)
}

func (s *testServer) onlyRequest(t *testing.T) recordedRequest {
	t.Helper()
	requests := s.Requests()
	require.Len(t, requests, 1)
	return requests[0]
}

type cliGateCall struct {
	Params []string
	Env    map[string]string
}

func decodeCliGateCall(t *testing.T, body []byte) cliGateCall {
	t.Helper()
	var fields map[string]json.RawMessage
	require.NoError(t, json.Unmarshal(body, &fields), "decode the body %s", body)

	names := make([]string, 0, len(fields))
	for name := range fields {
		names = append(names, name)
	}
	require.ElementsMatch(t, []string{"params", "env"}, names, "the fields of the body %s", body)

	call := cliGateCall{Env: map[string]string{}}
	require.NoError(t, json.Unmarshal(fields["params"], &call.Params), "decode the params of the body %s", body)
	require.NoError(t, json.Unmarshal(fields["env"], &call.Env), "decode the env of the body %s", body)
	return call
}

var (
	errorReplyError = jsonError{
		Code:    1013,
		Message: "Bad request.",
		Errors:  []jsonErrorInner{{Property: "params", Message: "The value is not valid."}},
	}
	cliGateFailure = cliGateResponce{Code: 1, Stdout: "partial output", Stderr: "Unknown option\n"}
)

type errorCases struct {
	needReauth    bool
	statusError   string
	cliGateFailed string
}

func checkErrors(t *testing.T, path string, call func(c *resty.Client) error, want errorCases) {
	t.Helper()
	t.Run("HTTP 403", func(t *testing.T) {
		_, c := startServer(t, map[string]reply{path: readReply(t, "forbidden.json")})
		assert.Equal(t, authError{server: c.BaseURL, needReauth: want.needReauth}, call(c))
	})
	t.Run("another error status", func(t *testing.T) {
		_, c := startServer(t, map[string]reply{path: valueReply(t, http.StatusBadRequest, errorReplyError)})
		assert.EqualError(t, call(c), want.statusError)
	})
	t.Run("CLI gate code is not 0", func(t *testing.T) {
		_, c := startServer(t, map[string]reply{path: valueReply(t, http.StatusOK, cliGateFailure)})
		assert.EqualError(t, call(c), want.cliGateFailed)
	})
}

func checkCliGateRequest(t *testing.T, s *testServer, path string, want cliGateCall) {
	t.Helper()
	r := s.onlyRequest(t)
	assert.Equal(t, http.MethodPost, r.Method)
	assert.Equal(t, path, r.Path)
	assert.Equal(t, want, decodeCliGateCall(t, r.Body))
}
