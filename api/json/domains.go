// Copyright 1999-2020. Plesk International GmbH. All rights reserved.

package json

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"

	"git.plesk.ru/~abashurov/pleskapp/api"
	"git.plesk.ru/~abashurov/pleskapp/features"
	"git.plesk.ru/~abashurov/pleskapp/locales"
	"git.plesk.ru/~abashurov/pleskapp/types"
	"git.plesk.ru/~abashurov/pleskapp/utils"
)

type jsonDomains struct {
	auth   api.Auth
	client *http.Client
}

type createDomainRequest struct {
	Name            string          `json:"name"`
	Description     string          `json:"description"`
	HostingType     string          `json:"hosting_type"`
	HostingSettings hostingSettings `json:"hosting_settings"`
	IPAddresses     []string        `json:"ip_addresses"`
}

type hostingSettings struct {
	FtpLogin    string `json:"ftp_login"`
	FtpPassword string `json:"ftp_password"`
}

type actionDomainResponse struct {
	ID   int    `json:"id"`
	GUID string `json:"guid"`
}

func NewDomains(a api.Auth) jsonDomains {
	return jsonDomains{
		auth:   a,
		client: getClient(a.GetIgnoreSsl()),
	}
}

func (j jsonDomains) CreateDomain(d string, ipa types.ServerIPAddresses) (*api.DomainInfo, error) {
	if len(ipa.IPv4) > 1 || len(ipa.IPv6) > 1 {
		return nil, errors.New(locales.L.Get("errors.ip.address.class.limit"))
	}

	var ip = []string{}

	for _, a := range ipa.IPv4 {
		ip = append(ip, a)
	}
	for _, a := range ipa.IPv6 {
		ip = append(ip, a)
	}

	p := createDomainRequest{
		Name:        d,
		Description: "",
		HostingType: "virtual",
		HostingSettings: hostingSettings{
			FtpLogin:    utils.GenUsername(16),
			FtpPassword: utils.GenPassword(32),
		},
		IPAddresses: ip,
	}

	jd, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", api.GetApiUrl(j.auth, "/api/v2/domains"), bytes.NewBuffer(jd))
	if err != nil {
		return nil, err
	}

	addBasicHeaders(req, j.auth.GetApiKey())

	res, err := doAndThenCheckAuthFailure(j.client, req, j.auth.GetAddress())
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var status actionDomainResponse
	e, err := tryParseResponceOrParseError(data, &status)
	if err != nil {
		return nil, err
	}

	if e != nil || status.GUID == "" {
		if e.Code != 0 || len(e.Errors) != 0 {
			_ = j.RemoveDomain(d)

			return nil, jsonErrorToError(*e)
		}
	}

	info, err := j.GetDomain(d)
	return &info, err
}

func (j jsonDomains) AddDomainFeatures(domain string, featureList []string) error {
	var featureProvider = features.FeatureProvider{
		IsWindows: j.auth.GetIsWindows(),
	}

	for _, featureStr := range featureList {
		feature := features.GetFeatureByString(featureStr)

		if feature != nil {
			packet, err := featureProvider.GetFeaturePackage(domain, *feature)
			if err != nil {
				return err
			}

			var req, _ = http.NewRequest("POST", api.GetApiUrl(j.auth, "/api/v2/cli/domain/call"), bytes.NewBuffer(packet))
			addBasicHeaders(req, j.auth.GetApiKey())

			res, err := doAndThenCheckAuthFailure(j.client, req, j.auth.GetAddress())
			if err != nil {
				return err
			}

			data, err := ioutil.ReadAll(res.Body)
			if err != nil {
				return err
			}

			var r cliGateResponce
			e, err := tryParseResponceOrParseError(data, &r)
			if err != nil {
				return err
			}

			if e != nil {
				if e.Code != 0 || len(e.Errors) != 0 {
					return jsonErrorToError(*e)
				}
			}

			if r.Code != 0 {
				return jsonCliGateResponceToError(r)
			}
		}
	}

	return nil
}

func (j jsonDomains) GetDomain(d string) (api.DomainInfo, error) {
	req, err := http.NewRequest("GET", api.GetApiUrl(j.auth, "/api/v2/domains/?name="+d), bytes.NewBuffer([]byte{}))
	if err != nil {
		return api.DomainInfo{}, err
	}

	addBasicHeaders(req, j.auth.GetApiKey())

	var ds []domainInfo
	res, err := doAndThenCheckAuthFailure(j.client, req, j.auth.GetAddress())
	if err != nil {
		return api.DomainInfo{}, err
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return api.DomainInfo{}, err
	}

	err = json.Unmarshal(data, &ds)
	if err != nil {
		return api.DomainInfo{}, err
	}

	if len(ds) == 0 {
		return api.DomainInfo{}, errors.New(locales.L.Get("errors.domain.unknown", d))
	}

	return api.DomainInfo{
		ID:             ds[0].ID,
		Name:           ds[0].Name,
		HostingType:    ds[0].HostingType,
		ParentDomainID: ds[0].BaseDomainID,
		GUID:           ds[0].GUID,
		WWWRoot:        ds[0].WWWRoot,
	}, nil
}

func (j jsonDomains) RemoveDomain(d string) error {
	info, err := j.GetDomain(d)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("DELETE", api.GetApiUrl(j.auth, "/api/v2/domains/"+strconv.Itoa(info.ID)), bytes.NewBuffer([]byte{}))
	if err != nil {
		return err
	}

	addBasicHeaders(req, j.auth.GetApiKey())

	res, err := doAndThenCheckAuthFailure(j.client, req, j.auth.GetAddress())
	if err != nil {
		return err
	}

	var data, _ = ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	var status actionDomainResponse
	e, err := tryParseResponceOrParseError(data, &status)
	if err != nil {
		return err
	}

	if e != nil || status.GUID == "" {
		if e.Code != 0 || len(e.Errors) != 0 {
			return jsonErrorToError(*e)
		}
	}

	return nil
}

func (j jsonDomains) ListDomains() ([]api.DomainInfo, error) {
	var req, _ = http.NewRequest("GET", api.GetApiUrl(j.auth, "/api/v2/domains"), bytes.NewBuffer([]byte{}))
	addBasicHeaders(req, j.auth.GetApiKey())

	var d []domainInfo
	res, err := doAndThenCheckAuthFailure(j.client, req, j.auth.GetAddress())
	if err != nil {
		return jsonDomainInfoToInfo(d), err
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return jsonDomainInfoToInfo(d), err
	}

	err = json.Unmarshal(data, &d)
	if err != nil {
		return jsonDomainInfoToInfo(d), err
	}

	return jsonDomainInfoToInfo(d), nil
}
