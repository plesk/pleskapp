// Copyright 1999-2020. Plesk International GmbH. All rights reserved.

package cmd

import (
	"git.plesk.ru/~abashurov/pleskapp/locales"
	"github.com/spf13/cobra"
)

var AppsCmd = &cobra.Command{Use: "apps", Short: locales.L.Get("app.description")}
