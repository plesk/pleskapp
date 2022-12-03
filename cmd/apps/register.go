// Copyright 1999-2022. Plesk International GmbH.

package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/plesk/pleskapp/plesk/internal/actions"
	"github.com/plesk/pleskapp/plesk/internal/config"
	"github.com/plesk/pleskapp/plesk/internal/features"
	"github.com/plesk/pleskapp/plesk/internal/locales"
	"github.com/spf13/cobra"
)

var registerCmd = &cobra.Command{
	Use:   "register [SERVER] [DOMAIN] [PATH]",
	Short: locales.L.Get("app.register.description"),
	RunE: func(cmd *cobra.Command, args []string) error {
		sPath, _ := cmd.Flags().GetString("target-path")
		feat, _ := cmd.Flags().GetStringSlice("features")
		overwrite, _ := cmd.Flags().GetBool("overwrite")
		path := args[2]

		server, err := config.GetServer(args[0])
		if err != nil {
			return err
		}

		domain, err := config.GetDomain(*server, args[1])
		if err != nil {
			return err
		}

		stat, err := os.Stat(path)
		if err != nil || !stat.IsDir() {
			return errors.New(locales.L.Get("app.register.flag.invalid"))
		}

		var kFeat []string
		for _, f := range feat {
			var k *features.Feature = features.GetFeatureByString(f)

			if k != nil {
				kFeat = append(kFeat, f)
			} else {
				fmt.Println(locales.L.Get("app.register.flag.feature.unknown", f))
			}
		}

		cmd.SilenceUsage = true
		err = actions.AppAdd(*server, *domain, kFeat, sPath, path, overwrite)

		if err == nil {
			fmt.Println(locales.L.Get("app.register.success", path))
		}

		return err
	},
	Args: cobra.ExactArgs(3),
}

func init() {
	registerCmd.Flags().StringSliceP("features", "f", []string{"php74", "nginx"}, locales.L.Get("app.register.features.flag"))
	registerCmd.Flags().StringP("target-path", "d", "/", locales.L.Get("app.register.target.path.flag"))
	registerCmd.Flags().BoolP("overwrite", "o", false, locales.L.Get("app.register.overwrite.flag"))

	AppsCmd.AddCommand(registerCmd)
}
