// Copyright 1999-2020. Plesk International GmbH. All rights reserved.

package actions

import (
	"fmt"
	"time"

	"git.plesk.ru/projects/SBX/repos/pleskapp/api"
	"git.plesk.ru/projects/SBX/repos/pleskapp/api/factory"
	"git.plesk.ru/projects/SBX/repos/pleskapp/config"
	"git.plesk.ru/projects/SBX/repos/pleskapp/types"
	"git.plesk.ru/projects/SBX/repos/pleskapp/utils"
)

func getServerInfo(a types.ServerAuth) (*api.ServerInfo, *types.ServerIPAddresses, *[]types.DatabaseServer, error) {
	apiI := factory.GetServerInfo(a)
	info, err := apiI.GetInfo()
	if err != nil {
		return nil, nil, nil, err
	}

	ipAddr, err := apiI.GetIpAddresses()
	if err != nil {
		return nil, nil, nil, err
	}

	apiD := factory.GetDatabaseManagement(a)
	dbs, err := apiD.ListDatabaseServers()
	if err != nil {
		return nil, nil, nil, err
	}

	dbst := []types.DatabaseServer{}
	for _, i := range dbs {
		dbst = append(dbst, types.DatabaseServer{
			ID:        i.ID,
			Host:      i.Host,
			Port:      i.Port,
			Type:      i.Type,
			IsDefault: i.IsDefault,
			IsLocal:   i.IsLocal,
		})
	}

	return &info, &ipAddr, &dbst, err
}

func ServerAdd(host string, ignoreSsl bool) error {
	_, err := config.GetServer(host)
	if err == nil {
		return fmt.Errorf("Server with address " + host + " is already registered")
	}
	pass, err := utils.RequestPassword("Enter \"admin\" user password for server " + host)
	if err != nil {
		return err
	}

	auth := types.ServerAuth{
		Address:   host,
		Port:      "8443",
		IgnoreSsl: ignoreSsl,
		IsWindows: false,
		Login:     "admin",
		Password:  pass,
	}

	apiA := factory.GetAuthentication(auth)
	key, err := apiA.GetAPIKey(auth)
	if err != nil {
		return err
	}

	auth.APIKey = key
	h := types.Server{
		Host:      auth.Address,
		IgnoreSsl: auth.IgnoreSsl,
		APIKey:    auth.APIKey,
	}
	return ServerUpdate(h)
}

func ServerLogin(host types.Server) error {
	api := factory.GetAuthentication(host.GetServerAuth())
	str, err := api.GetLoginLink(host.GetServerAuth())
	if err != nil {
		return err
	}

	fmt.Println("Generated one-time login link: " + str)

	return nil
}

func ServerList() error {
	servers := config.GetServers()

	fmt.Printf("Address\tVersion\tAPI key\tIP v4\tIP v6\n")
	for i := 0; i < len(servers); i++ {
		server := servers[i]
		fmt.Printf(
			"%s\t%s\t%s\t%s\t%s\n",
			server.Host,
			server.Info.Version,
			server.APIKey,
			server.Info.IP.IPv4,
			server.Info.IP.IPv6,
		)
	}

	return nil
}

// ServerUpdate reloads and recaches server info
func ServerUpdate(host types.Server) error {
	info, ipAddr, dbs, err := getServerInfo(host.GetServerAuth())
	if err != nil {
		return err
	}

	host.Info.IP = *ipAddr
	host.Info.IsWindows = info.IsWindows
	host.Info.Version = info.Version
	host.DatabaseServers = *dbs
	host.UpdatedAt = time.Now()

	s, _ := utils.FilterServers(config.GetServers(), host.Host)
	config.SetServers(append(s, host))

	return nil
}

func ServerReauth(host types.Server) error {
	pass, err := utils.RequestPassword("Enter \"admin\" user password for server " + host.Host)
	if err != nil {
		return err
	}

	auth := types.ServerAuth{
		Address:   host.Host,
		Port:      "8443",
		IgnoreSsl: host.IgnoreSsl,
		IsWindows: false,
		Login:     "admin",
		Password:  pass,
	}

	apiA := factory.GetAuthentication(auth)
	key, err := apiA.GetAPIKey(auth)
	if err != nil {
		return err
	}

	host.APIKey = key
	return ServerUpdate(host)
}

func ServerRemove(host types.Server) error {
	var keepServers, removeServers = utils.FilterServers(config.GetServers(), host.Host)

	if len(removeServers) == 0 {
		return types.ServerNotFound{Server: host.Host}
	}

	for _, server := range removeServers {
		fmt.Println("Removing server and its API key")

		var auth types.ServerAuth = server.GetServerAuth()
		var api api.Authentication = factory.GetAuthentication(auth)

		api.RemoveAPIKey(auth)
	}

	config.SetServers(keepServers)

	return nil
}
