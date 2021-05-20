// Copyright 1999-2021. Plesk International GmbH.

package actions

import (
	"errors"
	"fmt"
	"github.com/pkg/browser"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/plesk/pleskapp/plesk/api"
	"github.com/plesk/pleskapp/plesk/api/factory"
	"github.com/plesk/pleskapp/plesk/config"
	"github.com/plesk/pleskapp/plesk/types"
	"github.com/plesk/pleskapp/plesk/utils"
)

const ADMIN_USER = "admin"

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
	login := ADMIN_USER
	pass, err := utils.RequestPassword("Enter \"admin\" user password for server " + host + ":")
	if err != nil {
		return err
	}

	auth := types.ServerAuth{
		Address:   host,
		Port:      "8443",
		IgnoreSsl: ignoreSsl,
		IsWindows: false,
		Login:     &login,
		Password:  &pass,
	}

	apiA := factory.GetAuthentication(auth)
	key, err := apiA.GetAPIKey(auth)
	if err != nil {
		return err
	}

	auth.APIKey = &key
	h := types.Server{
		Host:      auth.Address,
		IgnoreSsl: auth.IgnoreSsl,
		APIKey:    *auth.APIKey,
	}
	return ServerUpdate(h)
}

func ServerLogin(host types.Server, generateOnly bool) error {
	api := factory.GetAuthentication(host.GetServerAuth())
	url, err := api.GetLoginLink(host.GetServerAuth())
	if err != nil {
		return err
	}

	if generateOnly {
		fmt.Println("Generated one-time login link: " + url)
	} else {
		fmt.Println("Opening the browser with one-time login link...")
		err := browser.OpenURL(url)
		if err != nil {
			fmt.Println("Unable to open the browser, use one-time login link on your own: " + url)
		}
	}

	return nil
}

func ServerSsh(host types.Server) error {
	fmt.Printf("Login to %s using SSH...\n", host.Host)

	cmd := exec.Command("ssh", host.Host)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func ServerList() error {
	for _, i := range config.GetServers() {
		fmt.Printf(
			"Address: %s\nVersion: %s\nIPv4: %s\nIPv6: %s\n\n",
			i.Host,
			i.Info.Version,
			strings.Join(i.Info.IP.IPv4, ","),
			strings.Join(i.Info.IP.IPv6, ","),
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

	servers := config.GetServers()
	found := false
	for index, target := range servers {
		if target.Host == host.Host {
			servers[index] = host
			found = true
		}
	}

	if !found {
		servers = append([]types.Server{host}, servers...)
	}

	config.SetServers(servers)

	return nil
}

func ServerReauth(host types.Server) error {
	login := "admin"
	pass, err := utils.RequestPassword("Enter \"admin\" user password for server " + host.Host + ":")
	if err != nil {
		return err
	}

	auth := types.ServerAuth{
		Address:   host.Host,
		Port:      "8443",
		IgnoreSsl: host.IgnoreSsl,
		IsWindows: false,
		Login:     &login,
		Password:  &pass,
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

func DefaultServer() (string, error) {
	servers := config.GetServers()

	if len(servers) == 0 {
		return "", errors.New("context is not defined")
	}

	defaultServer := servers[0]

	return defaultServer.Host, nil
}
