// Copyright 1999-2024. Plesk International GmbH.

package json

import (
	"github.com/go-resty/resty/v2"
	"github.com/plesk/pleskapp/plesk/internal/api"
	"github.com/plesk/pleskapp/plesk/internal/types"
)

type jsonInfo struct {
	client *resty.Client
}

type serverInfo struct {
	Platform           string `json:"platform"`
	Hostname           string `json:"hostname"`
	GUID               string `json:"guid"`
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

func NewInfo(c *resty.Client) jsonInfo {
	return jsonInfo{
		client: c,
	}
}

func (j jsonInfo) GetInfo() (api.ServerInfo, error) {
	res, err := j.client.R().
		SetResult(&serverInfo{}).
		SetError(&jsonError{}).
		Get("/api/v2/server")

	if err != nil {
		return api.ServerInfo{}, err
	}

	if res.IsSuccess() {
		r := res.Result().(*serverInfo)
		return api.ServerInfo{
			IsWindows: r.Platform == "Windows",
			Version:   r.PanelVersion + "." + r.PanelUpdateVersion,
			Revision:  r.PanelRevision,
			BuildDate: r.PanelBuildDate,
		}, nil
	}

	if res.StatusCode() == 403 {
		return api.ServerInfo{}, authError{server: j.client.BaseURL, needReauth: true}
	}

	var r = res.Error().(*jsonError)
	return api.ServerInfo{}, jsonErrorToError(*r)
}

func (j jsonInfo) GetIPAddresses() (types.ServerIPAddresses, error) {
	res, err := j.client.R().
		SetResult([]serverIPAddresses{}).
		SetError(&jsonError{}).
		Get("/api/v2/server/ips")

	if err != nil {
		return types.ServerIPAddresses{}, err
	}

	if res.IsSuccess() {
		var r = res.Result().(*[]serverIPAddresses)
		ip := types.ServerIPAddresses{
			IPv4: []string{},
			IPv6: []string{},
		}
		for _, a := range *r {
			if a.IPv4 != "" {
				ip.IPv4 = append(ip.IPv4, a.IPv4)
			}
			if a.IPv6 != "" {
				ip.IPv6 = append(ip.IPv6, a.IPv6)
			}
		}

		return ip, nil
	}

	if res.StatusCode() == 403 {
		return types.ServerIPAddresses{}, authError{server: j.client.BaseURL, needReauth: true}
	}

	var r = res.Error().(*jsonError)
	return types.ServerIPAddresses{}, jsonErrorToError(*r)
}
