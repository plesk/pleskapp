// Copyright 1999-2021. Plesk International GmbH.

package config

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/plesk/pleskapp/plesk/internal/locales"
	"os"
	"strings"
	"sync"

	"github.com/plesk/pleskapp/plesk/internal/types"
	"github.com/plesk/pleskapp/plesk/internal/utils"
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
		return fmt.Errorf("cannot save configuration")
	}

	str, err := json.MarshalIndent(globalConfig.config, "", "  ")
	if err != nil {
		return err
	}

	f.Truncate(0)
	f.Seek(0, 0)

	_, err = f.Write(str)
	return err
}

func GetServer(host string) (*types.Server, error) {
	servers := GetServers()

	var foundServer *types.Server
	found := 0
	for index, server := range servers {
		if strings.HasPrefix(server.Host, host) {
			found++
			foundServer = &servers[index]
		}
	}

	if found == 1 {
		return foundServer, nil
	}

	if found > 1 {
		return nil, fmt.Errorf(locales.L.Get("errors.multiple.servers", host))
	}

	return nil, types.ServerNotFound{Server: host}
}

func GetServerByArgs(args []string) (*types.Server, error) {
	var serverName string
	if len(args) == 0 {
		serverName, _ = DefaultServer()
	} else {
		serverName = args[0]
	}

	server, err := GetServer(serverName)
	if err != nil {
		return nil, err
	}
	return server, nil
}

func DefaultServer() (string, error) {
	servers := GetServers()

	if len(servers) == 0 {
		return "", errors.New("context is not defined")
	}

	defaultServer := servers[0]

	return defaultServer.Host, nil
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

func SetDomains(host *types.Server, newData []types.Domain) {
	servers, _ := utils.FilterServers(GetServers(), host.Host)
	host.Domains = newData
	SetServers(append([]types.Server{*host}, servers...))
}

func SetDomain(host types.Server, domain types.Domain) {
	domains, _ := utils.FilterDomains(host.Domains, domain.Name)
	SetDomains(&host, append(domains, domain))
}
