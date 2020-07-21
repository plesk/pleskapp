// Copyright 1999-2020. Plesk International GmbH. All rights reserved.

package cmd

import (
	"git.plesk.ru/~abashurov/pleskapp/locales"
	"github.com/spf13/cobra"
)

var DomainsCmd = &cobra.Command{Use: "domains", Short: locales.L.Get("domain.description")}
