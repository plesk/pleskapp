// Copyright 1999-2020. Plesk International GmbH. All rights reserved.

package actions

import (
	"errors"
	"fmt"

	"git.plesk.ru/~abashurov/pleskapp/api/factory"
	"git.plesk.ru/~abashurov/pleskapp/config"
	"git.plesk.ru/~abashurov/pleskapp/locales"
	"git.plesk.ru/~abashurov/pleskapp/types"
	"git.plesk.ru/~abashurov/pleskapp/utils"
)

func validateIps(val types.ServerIPAddresses, comp types.ServerIPAddresses) (bool, bool) {
	var v4Valid bool = false
	var v6Valid bool = false
	if len(val.IPv4) == 0 {
		v4Valid = true
	} else {
		for _, i := range comp.IPv4 {
			if i == val.IPv4[0] {
				v4Valid = true
			}
		}
	}

	if len(val.IPv6) == 0 {
		v6Valid = true
	} else {
		for _, i := range comp.IPv6 {
			if i == val.IPv6[0] {
				v6Valid = true
			}
		}
	}

	return v4Valid, v6Valid
}

func DomainAdd(host types.Server, domain string, ipa types.ServerIPAddresses) error {
	_, err := config.GetDomain(host, domain)
	if err == nil {
		return errors.New(locales.L.Get("errors.domain.already.exists", domain))
	}

	if len(ipa.IPv4) != 1 && len(ipa.IPv6) != 1 {
		return errors.New(locales.L.Get("errors.ip.address.required"))
	}

	v4Valid, v6Valid := validateIps(ipa, host.Info.IP)
	if v4Valid != true || v6Valid != true {
		return errors.New(locales.L.Get("errors.ip.address.not.cached", host.Host))
	}

	api := factory.GetDomainManagement(host.GetServerAuth())
	d, err := api.CreateDomain(domain, ipa)
	if err == nil {
		config.SetDomain(host, types.Domain{
			Name: d.Name,
			GUID: d.GUID,
		})
	}

	return err
}

func DomainList(host types.Server) error {
	for _, domain := range host.Domains {
		utils.Log.Print(fmt.Sprintf("%s\t%s\n", domain.Name, domain.GUID))
	}

	return nil
}

func DomainReload(host types.Server) error {
	api := factory.GetDomainManagement(host.GetServerAuth())
	domains, err := api.ListDomains()
	if err != nil {
		return err
	}

	var newdomains []types.Domain
	for _, d := range domains {
		newdomains = append(newdomains, types.Domain{
			Name: d.Name,
			GUID: d.GUID,
		})
	}

	for _, d := range host.Domains {
		others, this := utils.FilterDomains(newdomains, d.Name)
		if len(this) == 0 {
			// This domain no longer exists on server
			newdomains = others
			continue
		}

		// This domain still exists, and may have more info in config
		newdomains = append(others, d)
	}

	config.SetDomains(host, newdomains)
	return nil
}

func DomainDelete(host types.Server, domain types.Domain) error {
	api := factory.GetDomainManagement(host.GetServerAuth())
	return api.RemoveDomain(domain.Name)
}
