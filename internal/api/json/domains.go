// Copyright 1999-2022. Plesk International GmbH.

package json

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/plesk/pleskapp/plesk/internal/api"
	"github.com/plesk/pleskapp/plesk/internal/features"
	"github.com/plesk/pleskapp/plesk/internal/locales"
	"github.com/plesk/pleskapp/plesk/internal/types"
	"github.com/plesk/pleskapp/plesk/internal/utils"
)

type jsonDomains struct {
	client *resty.Client
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

func NewDomains(c *resty.Client) jsonDomains {
	return jsonDomains{
		client: c,
	}
}

// TODO: Make sure REST API returns sysuser on GET /api/v2/domains
func (j jsonDomains) getDomainSysUser(d string) (string, error) {
	p := cliGateRequest{
		Params: []string{
			"--info",
			d,
		},
		Env: map[string]string{},
	}
	req, err := json.Marshal(p)
	if err != nil {
		return "", err
	}

	res, err := j.client.R().
		SetBody(req).
		SetResult(&cliGateResponce{}).
		SetError(&jsonError{}).
		Post("/api/v2/cli/domain/call")

	if err != nil {
		return "", err
	}

	if res.IsSuccess() {
		var r *cliGateResponce = res.Result().(*cliGateResponce)
		if r.Code != 0 {
			return "", jsonCliGateResponceToError(*r)
		}

		for _, l := range strings.Split(r.Stdout, "\n") {
			if strings.HasPrefix(l, "FTP Login") {
				p := strings.Split(l, " ")
				return p[len(p)-1], nil
			}
		}

		return "", errors.New(locales.L.Get("api.errors.domain.info.not.found"))
	}

	if res.StatusCode() == 403 {
		return "", authError{server: j.client.HostURL, needReauth: true}
	}

	var r *jsonError = res.Error().(*jsonError)
	return "", errors.New(locales.L.Get("api.errors.domain.info.failed", r.Code, r.Message))
}

func (j jsonDomains) CreateDomain(d string, ipa types.ServerIPAddresses) (*api.DomainInfo, error) {
	if len(ipa.IPv4) > 1 || len(ipa.IPv6) > 1 {
		return nil, errors.New(locales.L.Get("errors.ip.address.class.limit"))
	}

	var ip []string
	ip = append(ip, ipa.IPv4...)
	ip = append(ip, ipa.IPv6...)

	p := createDomainRequest{
		Name:        d,
		Description: "",
		HostingType: "virtual",
		HostingSettings: hostingSettings{
			FtpLogin:    utils.GenerateUsername(16),
			FtpPassword: utils.GeneratePassword(32),
		},
		IPAddresses: ip,
	}

	req, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}

	res, err := j.client.R().
		SetBody(req).
		SetResult(&actionDomainResponse{}).
		SetError(&jsonError{}).
		Post("/api/v2/domains")

	if err != nil {
		return nil, err
	}

	if res.IsSuccess() {
		var _ *actionDomainResponse = res.Result().(*actionDomainResponse)
		info, err := j.GetDomain(d)
		return &info, err
	}

	if res.StatusCode() == 403 {
		return nil, authError{server: j.client.HostURL, needReauth: true}
	}

	var r *jsonError = res.Error().(*jsonError)
	return nil, jsonErrorToError(*r)
}

func (j jsonDomains) AddDomainFeatures(domain string, featureList []string, isWin bool) error {
	var featureProvider = features.FeatureProvider{
		IsWindows: isWin,
	}

	for _, featureStr := range featureList {
		feature := features.GetFeatureByString(featureStr)

		if feature != nil {
			p, err := featureProvider.GetFeaturePackage(domain, *feature)
			if err != nil {
				return err
			}

			res, err := j.client.R().
				SetBody(p).
				SetResult(&cliGateResponce{}).
				SetError(&jsonError{}).
				Post("/api/v2/cli/domain/call")

			if err != nil {
				return err
			}

			if res.IsSuccess() {
				var r *cliGateResponce = res.Result().(*cliGateResponce)
				if r.Code == 0 {
					return nil
				}

				return jsonCliGateResponceToError(*r)
			}

			if res.StatusCode() == 403 {
				return authError{server: j.client.HostURL, needReauth: true}
			}

			var r *jsonError = res.Error().(*jsonError)
			return jsonErrorToError(*r)
		}
	}

	return nil
}

func (j jsonDomains) GetDomain(d string) (api.DomainInfo, error) {
	res, err := j.client.R().
		SetResult([]domainInfo{}).
		SetError(&jsonError{}).
		SetQueryParam("name", d).
		Get("/api/v2/domains")

	if err != nil {
		return api.DomainInfo{}, err
	}

	if res.IsSuccess() {
		var r *[]domainInfo = res.Result().(*[]domainInfo)
		if len(*r) == 0 {
			return api.DomainInfo{}, errors.New(locales.L.Get("errors.domain.unknown", d))
		}

		s, err := j.getDomainSysUser(d)
		if err != nil {
			return api.DomainInfo{}, err
		}

		return api.DomainInfo{
			ID:             (*r)[0].ID,
			Name:           (*r)[0].Name,
			HostingType:    (*r)[0].HostingType,
			ParentDomainID: (*r)[0].BaseDomainID,
			GUID:           (*r)[0].GUID,
			WWWRoot:        (*r)[0].WWWRoot,
			Sysuser:        s,
		}, nil
	}

	if res.StatusCode() == 403 {
		return api.DomainInfo{}, authError{server: j.client.HostURL, needReauth: true}
	}

	var r *jsonError = res.Error().(*jsonError)
	return api.DomainInfo{}, jsonErrorToError(*r)
}

func (j jsonDomains) RemoveDomain(d string) error {
	info, err := j.GetDomain(d)
	if err != nil {
		return err
	}

	res, err := j.client.R().
		SetError(&jsonError{}).
		Delete("/api/v2/domains/" + strconv.Itoa(info.ID))

	if err != nil {
		return err
	}

	if res.StatusCode() == 403 {
		return authError{server: j.client.HostURL, needReauth: true}
	}

	if res.IsError() {
		var r *jsonError = res.Error().(*jsonError)
		return jsonErrorToError(*r)
	}

	return nil
}

func (j jsonDomains) ListDomains() ([]api.DomainInfo, error) {
	res, err := j.client.R().
		SetResult([]domainInfo{}).
		SetError(&jsonError{}).
		Get("/api/v2/domains")

	if err != nil {
		return []api.DomainInfo{}, err
	}

	if res.IsSuccess() {
		var r *[]domainInfo = res.Result().(*[]domainInfo)
		return jsonDomainInfoToInfo(*r), nil
	}

	if res.StatusCode() == 403 {
		return []api.DomainInfo{}, authError{server: j.client.HostURL, needReauth: true}
	}

	var r *jsonError = res.Error().(*jsonError)
	return []api.DomainInfo{}, jsonErrorToError(*r)
}
