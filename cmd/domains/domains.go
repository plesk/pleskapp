// Copyright 1999-2020. Plesk International GmbH.

package cmd

import (
	"git.plesk.ru/projects/SBX/repos/pleskapp/locales"
	"github.com/spf13/cobra"
)

var DomainsCmd = &cobra.Command{Use: "domains", Short: locales.L.Get("domain.description")}
