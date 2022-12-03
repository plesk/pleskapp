// Copyright 1999-2022. Plesk International GmbH.

package cmd

import (
	"github.com/plesk/pleskapp/plesk/internal/locales"
	"github.com/spf13/cobra"
)

var AppsCmd = &cobra.Command{Use: "apps", Short: locales.L.Get("app.description")}
