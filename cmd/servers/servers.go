// Copyright 1999-2020. Plesk International GmbH.

package cmd

import (
	"github.com/plesk/pleskapp/plesk/locales"
	"github.com/spf13/cobra"
)

var ServersCmd = &cobra.Command{Use: "servers", Short: locales.L.Get("server.description")}
