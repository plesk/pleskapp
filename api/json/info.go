// Copyright 1999-2020. Plesk International GmbH. All rights reserved.

package json

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"git.plesk.ru/projects/SBX/repos/pleskapp/api"
	"git.plesk.ru/projects/SBX/repos/pleskapp/types"
)

type jsonInfo struct {
	auth   api.Auth
	client *http.Client
}

type serverInfo struct {
	Platform           string `json:"platform"`
	Hostname           string `json:"hostname"`
	Guid               string `json:"guid"`
	PanelVersion       string `json:"panel_version"`
	PanelRevision      string `json:"panel_revision"`
	PanelBuildDate     string `json:"panel_build_date"`
	PanelUpdateVersion string `json:"panel_update_version"`
	ExtensionVersion   string `json:"extension_version"`
	ExtensionRelease   string `json:"extension_release"`
}

type serverIPAddresses struct {
	IPv4      string `json:"ipv4"`
	IPv6      string `json:"ipv6"`
	Netmask   string `json:"netmask"`
	Interface string `json:"interface"`
	Type      string `json:"type"`
}

func NewInfo(a api.Auth) jsonInfo {
	return jsonInfo{
		auth:   a,
		client: getClient(a.GetIgnoreSsl()),
	}
}

func (j jsonInfo) GetInfo() (api.ServerInfo, error) {
	req, err := http.NewRequest("GET", api.GetApiUrl(j.auth, "/api/v2/server"), bytes.NewBuffer([]byte{}))
	if err != nil {
		return api.ServerInfo{}, err
	}

	addBasicHeaders(req, j.auth.GetApiKey())

	res, err := doAndThenCheckAuthFailure(j.client, req, j.auth.GetAddress())
	if err != nil {
		return api.ServerInfo{}, err
	}

	d, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return api.ServerInfo{}, err
	}

	var info serverInfo
	e, err := tryParseResponceOrParseError(d, &info)
	if err != nil {
		return api.ServerInfo{}, err
	}

	if e != nil {
		return api.ServerInfo{}, fmt.Errorf("Failed to check server version using provided API key: [%d: %s]", e.Code, e.Message)
	}

	return api.ServerInfo{
		IsWindows: info.Platform == "Windows",
		Version:   info.PanelVersion + "." + info.PanelUpdateVersion,
	}, nil
}

func (j jsonInfo) GetIpAddresses() (types.ServerIPAddresses, error) {
	ipAddresses := types.ServerIPAddresses{
		IPv4: []string{},
		IPv6: []string{},
	}

	var req, err = http.NewRequest("GET", api.GetApiUrl(j.auth, "/api/v2/server/ips"), bytes.NewBuffer([]byte{}))
	if err != nil {
		return ipAddresses, err
	}

	addBasicHeaders(req, j.auth.GetApiKey())

	res, err := doAndThenCheckAuthFailure(j.client, req, j.auth.GetAddress())
	if err != nil {
		return ipAddresses, err
	}

	d, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return ipAddresses, err
	}

	var addr []serverIPAddresses

	e, err := tryParseResponceOrParseError(d, &addr)
	if err != nil {
		return ipAddresses, err
	}

	if e != nil {
		return ipAddresses, fmt.Errorf("Failed to check server version using provided API key: [%d: %s]", e.Code, e.Message)
	}

	for _, a := range addr {
		if a.IPv4 != "" {
			ipAddresses.IPv4 = append(ipAddresses.IPv4, a.IPv4)
		}
		if a.IPv6 != "" {
			ipAddresses.IPv6 = append(ipAddresses.IPv6, a.IPv6)
		}
	}

	return ipAddresses, nil
}
