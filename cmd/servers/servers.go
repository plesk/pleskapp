// Copyright 1999-2020. Plesk International GmbH. All rights reserved.

package cmd

import (
	"git.plesk.ru/projects/SBX/repos/pleskapp/locales"
	"github.com/spf13/cobra"
)

var ServersCmd = &cobra.Command{Use: "servers", Short: locales.L.Get("server.description")}
