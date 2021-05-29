// Copyright 1999-2021. Plesk International GmbH.

package actions

import (
	"errors"
	"fmt"
	"os"
	"sort"
	"text/tabwriter"

	"github.com/plesk/pleskapp/plesk/api/factory"
	"github.com/plesk/pleskapp/plesk/config"
	"github.com/plesk/pleskapp/plesk/locales"
	"github.com/plesk/pleskapp/plesk/types"
	"github.com/plesk/pleskapp/plesk/utils"
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
	if !v4Valid || !v6Valid {
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
	domains := host.Domains
	sort.Slice(domains, func(i, j int) bool { return domains[i].Name < domains[j].Name })

	w := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
	fmt.Fprintln(w, "DOMAIN\tGUID")
	for _, domain := range domains {
		fmt.Fprintf(w, "%s\t%s\n", domain.Name, domain.GUID)
	}
	w.Flush()

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
