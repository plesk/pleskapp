// Copyright 1999-2021. Plesk International GmbH.

package cmd

import (
	"github.com/plesk/pleskapp/plesk/locales"
	"github.com/spf13/cobra"
)

var AppsCmd = &cobra.Command{Use: "apps", Short: locales.L.Get("app.description")}
