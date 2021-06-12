// Copyright 1999-2021. Plesk International GmbH.

package actions

import (
	"fmt"
	"github.com/pkg/browser"
	"os"
	"os/exec"
	"sort"
	"text/tabwriter"
	"time"

	"github.com/plesk/pleskapp/plesk/internal/api"
	"github.com/plesk/pleskapp/plesk/internal/api/factory"
	"github.com/plesk/pleskapp/plesk/internal/config"
	"github.com/plesk/pleskapp/plesk/internal/types"
	"github.com/plesk/pleskapp/plesk/internal/utils"
)

const AdminUser = "admin"

func getServerInfo(a types.ServerAuth) (*api.ServerInfo, *types.ServerIPAddresses, *[]types.DatabaseServer, error) {
	apiI := factory.GetServerInfo(a)
	info, err := apiI.GetInfo()
	if err != nil {
		return nil, nil, nil, err
	}

	ipAddr, err := apiI.GetIPAddresses()
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
	login := AdminUser
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

func ServerSSH(host types.Server) error {
	fmt.Printf("Login to %s using SSH...\n", host.Host)

	cmd := exec.Command("ssh", host.Host)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func ServerList() error {
	servers := config.GetServers()
	sort.Slice(servers, func(i, j int) bool { return servers[i].Host < servers[j].Host })

	w := tabwriter.NewWriter(os.Stdout, 1, 1, 2, ' ', 0)
	_, _ = fmt.Fprintln(w, "HOST\tPLATFORM\tVERSION\tIPV4\tIPV6")

	for _, server := range servers {
		ipv4 := "-"
		if len(server.Info.IP.IPv4) > 0 {
			ipv4 = server.Info.IP.IPv4[0]
		}

		ipv6 := "-"
		if len(server.Info.IP.IPv6) > 0 {
			ipv6 = server.Info.IP.IPv6[0]
		}

		platform := "Linux"
		if server.Info.IsWindows {
			platform = "Windows"
		}

		_, _ = fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
			server.Host,
			platform,
			server.Info.Version,
			ipv4,
			ipv6,
		)
	}

	return w.Flush()
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
