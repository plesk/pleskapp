// Copyright 1999-2021. Plesk International GmbH.

package types

import "time"

type Server struct {
	Host            string           `json:"host"`
	IgnoreSsl       bool             `json:"ignore_ssl"`
	APIKey          string           `json:"api_key"`
	Info            ServerInfo       `json:"info"`
	UpdatedAt       time.Time        `json:"updated_at"`
	Domains         []Domain         `json:"domain"`
	DatabaseServers []DatabaseServer `json:"database_servers"`
}

func (s Server) GetServerAuth() ServerAuth {
	return ServerAuth{
		Address:   s.Host,
		Port:      "8443",
		IgnoreSsl: s.IgnoreSsl,
		IsWindows: s.Info.IsWindows,
		APIKey:    &s.APIKey,
	}
}

func (s Server) GetDatabaseServer(id int) *DatabaseServer {
	for _, i := range s.DatabaseServers {
		if i.ID == id {
			return &i
		}
	}
	return nil
}

func (s Server) GetDatabaseServerByType(dbt string) *DatabaseServer {
	var server *DatabaseServer

	for _, i := range s.DatabaseServers {
		if i.Type == dbt && i.IsDefault {
			return &i
		} else if i.Type == dbt {
			server = &i
		}
	}
	return server
}

type ServerInfo struct {
	IsWindows bool              `json:"is_windows"`
	Version   string            `json:"version"`
	IP        ServerIPAddresses `json:"addresses"`
}

type ServerIPAddresses struct {
	IPv4 []string `json:"ip_v4"`
	IPv6 []string `json:"ip_v6"`
}

type Domain struct {
	Name          string         `json:"name"`
	GUID          string         `json:"guid"`
	FTPUsers      []FtpUser      `json:"ftp_users"`
	Databases     []Database     `json:"databases"`
	DatabaseUsers []DatabaseUser `json:"database_users"`
}

type App struct {
	TargetPath string   `json:"target_path"`
	Features   []string `json:"features"`
	Server     string   `json:"server"`
	Domain     string   `json:"domain"`
}

type Config struct {
	Server []Server `json:"servers"`
}

type FtpUser struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Database struct {
	ID               int    `json:"id"`
	Name             string `json:"name"`
	Type             string `json:"type"`
	DatabaseServerID int    `json:"server_id"`
}

type NewDatabase struct {
	Name             string `json:"name"`
	Type             string `json:"type"`
	ParentDomain     string `json:"parent_domain"`
	DatabaseServerID int    `json:"server_id"`
}

type DatabaseUser struct {
	ID         int    `json:"id"`
	Login      string `json:"login"`
	Password   string `json:"password"`
	DatabaseID int    `json:"database_id"`
}

type NewDatabaseUser struct {
	Login            string  `json:"login"`
	Password         string  `json:"password"`
	ParentDomain     *string `json:"parent_domain"`
	DatabaseServerID *int    `json:"database_server_id"`
	DatabaseID       *int    `json:"database_id"`
}

type DatabaseServer struct {
	ID        int    `json:"id"`
	Host      string `json:"host"`
	Port      int    `json:"port"`
	Type      string `json:"type"`
	IsDefault bool   `json:"is_default"`
	IsLocal   bool   `json:"is_local"`
}

type ServerAuth struct {
	Address   string
	Port      string
	IgnoreSsl bool
	IsWindows bool
	Login     *string
	Password  *string
	APIKey    *string
}

// GetAddress impl Auth
func (a ServerAuth) GetAddress() string { return a.Address }

// GetPort impl Auth
func (a ServerAuth) GetPort() string { return a.Port }

// GetIgnoreSsl impl Auth
func (a ServerAuth) GetIgnoreSsl() bool { return a.IgnoreSsl }

// GetIsWindows impl Auth
func (a ServerAuth) GetIsWindows() bool { return a.IsWindows }

// GetLogin impl Auth
func (a ServerAuth) GetLogin() *string { return a.Login }

// GetPassword impl Auth
func (a ServerAuth) GetPassword() *string { return a.Password }

// GetApiKey impl Auth
func (a ServerAuth) GetApiKey() *string { return a.APIKey }
