// Copyright 1999-2020. Plesk International GmbH. All rights reserved.

package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"git.plesk.ru/projects/SBX/repos/pleskapp/types"
	"git.plesk.ru/projects/SBX/repos/pleskapp/utils"
)

var globalConfig config = config{}

type config struct {
	mutex  sync.Mutex
	config types.Config
}

func New(f *os.File) error {
	var c types.Config
	if f != nil {
		buf := bytes.NewBuffer([]byte{})
		_, err := buf.ReadFrom(f)
		if err != nil {
			return err
		}

		err = json.Unmarshal(buf.Bytes(), &c)
		if err != nil {
			return err
		}
	}

	globalConfig.config = c
	return nil
}

func Save(f *os.File) error {
	if f == nil {
		return fmt.Errorf("Cannot save configuration")
	}

	str, err := json.Marshal(globalConfig.config)
	if err != nil {
		return err
	}

	f.Truncate(0)
	f.Seek(0, 0)

	_, err = f.Write(str)
	return err
}

func GetServer(host string) (*types.Server, error) {
	_, server := utils.FilterServers(GetServers(), host)
	if len(server) == 0 {
		return nil, types.ServerNotFound{Server: host}
	}

	return &server[0], nil
}

func GetServers() []types.Server {
	globalConfig.mutex.Lock()
	defer globalConfig.mutex.Unlock()
	return globalConfig.config.Server
}

func GetDomain(host types.Server, domain string) (*types.Domain, error) {
	_, d := utils.FilterDomains(host.Domains, domain)
	if len(d) == 0 {
		return nil, types.DomainNotFound{Domain: domain, Server: host.Host}
	}

	return &d[0], nil
}

func GetDatabase(host types.Server, dbn string) (*types.Database, error) {
	for _, d := range host.Domains {
		_, db := utils.FilterDatabases(d.Databases, dbn)
		if len(db) != 0 {
			return &db[0], nil
		}
	}
	return nil, types.DatabaseNotFound{DbName: dbn, Server: host.Host}
}

func DeleteDatabase(host types.Server, dbn string) {
	var domain *types.Domain

	for _, d := range host.Domains {
		keep, remove := utils.FilterDatabases(d.Databases, dbn)
		if len(remove) != 0 {
			domain = &d
			domain.Databases = keep
		}
	}

	if domain != nil {
		SetDomain(host, *domain)
	}
}

func SetServers(newData []types.Server) {
	globalConfig.mutex.Lock()
	defer globalConfig.mutex.Unlock()
	globalConfig.config.Server = newData
}

func SetDomains(host types.Server, newData []types.Domain) {
	servers, _ := utils.FilterServers(GetServers(), host.Host)
	host.Domains = newData
	SetServers(append(servers, host))
}

func SetDomain(host types.Server, domain types.Domain) {
	domains, _ := utils.FilterDomains(host.Domains, domain.Name)
	SetDomains(host, append(domains, domain))
}
