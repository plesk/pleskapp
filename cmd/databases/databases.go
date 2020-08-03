// Copyright 1999-2020. Plesk International GmbH.

package cmd

import (
	"github.com/plesk/pleskapp/locales"
	"github.com/spf13/cobra"
)

var DatabasesCmd = &cobra.Command{Use: "databases", Short: locales.L.Get("database.description")}
