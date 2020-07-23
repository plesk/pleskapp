// Copyright 1999-2020. Plesk International GmbH. All rights reserved.

package cmd

import (
	"errors"
	"os"

	"git.plesk.ru/projects/SBX/repos/pleskapp/actions"
	"git.plesk.ru/projects/SBX/repos/pleskapp/config"
	"git.plesk.ru/projects/SBX/repos/pleskapp/features"
	"git.plesk.ru/projects/SBX/repos/pleskapp/locales"
	"git.plesk.ru/projects/SBX/repos/pleskapp/utils"
	"github.com/spf13/cobra"
)

var registerCmd = &cobra.Command{
	Use:   locales.L.Get("app.register.cmd"),
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
				utils.Log.Error(locales.L.Get("app.register.flag.feature.unknown", f))
			}
		}

		cmd.SilenceUsage = true
		return utils.Log.PrintSuccessOrError("app.register.success", nil, actions.AppAdd(*server, *domain, kFeat, sPath, path, overwrite))
	},
	Args: cobra.ExactArgs(3),
}

func init() {
	registerCmd.Flags().StringSliceP("features", "f", []string{"php74", "nginx"}, locales.L.Get("app.register.features.flag"))
	registerCmd.Flags().StringP("target-path", "d", "/", locales.L.Get("app.register.target.path.flag"))
	registerCmd.Flags().BoolP("overwrite", "o", false, locales.L.Get("app.register.overwrite.flag"))

	AppsCmd.AddCommand(registerCmd)
}
