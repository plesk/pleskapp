// Copyright 1999-2022. Plesk International GmbH.

package actions

import (
	"github.com/plesk/pleskapp/plesk/internal/api/factory"
	"github.com/plesk/pleskapp/plesk/internal/config"
	"github.com/plesk/pleskapp/plesk/internal/types"
	"github.com/plesk/pleskapp/plesk/internal/utils"
)

func FindCachedFtpUser(domain types.Domain) *types.FtpUser {
	if len(domain.FTPUsers) != 0 {
		return &domain.FTPUsers[0]
	}
	return nil
}

func FtpUserCreate(host types.Server, domain types.Domain, user *types.FtpUser) (*types.FtpUser, error) {
	if user == nil {
		user = &types.FtpUser{
			Login:    utils.GenerateUsername(16),
			Password: utils.GeneratePassword(32),
		}
	}

	api := factory.GetFTPUserManagement(host.GetServerAuth())
	_, err := api.CreateFtpUser(domain.Name, *user)
	if err != nil {
		return nil, err
	}

	domain.FTPUsers = append(domain.FTPUsers, *user)
	config.SetDomain(&host, domain)

	return user, nil
}
