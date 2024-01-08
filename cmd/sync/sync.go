// Copyright 1999-2024. Plesk International GmbH.

package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/plesk/pleskapp/plesk/internal/actions"
	"github.com/plesk/pleskapp/plesk/internal/config"
	"github.com/plesk/pleskapp/plesk/internal/locales"
	"github.com/spf13/cobra"
)

var SyncCmd = &cobra.Command{
	Use:   "sync [SERVER] [DOMAIN] [FILE ...]",
	Short: locales.L.Get("files.upload.description"),
	RunE: func(cmd *cobra.Command, args []string) error {
		var lastErr error = nil
		overwrite, _ := cmd.Flags().GetBool("overwrite")
		dry, _ := cmd.Flags().GetBool("dryRun")

		server, err := config.GetServer(args[0])
		if err != nil {
			return err
		}

		domain, err := config.GetDomain(*server, args[1])
		if err != nil {
			return err
		}

		// TODO: Optimize batch upload
		for i, p := range args {
			if i < 2 {
				continue
			}

			path, err := filepath.Abs(p)
			if err != nil {
				lastErr = err
				fmt.Println(locales.L.Get("errors.abspath.failed", p, err.Error()))
				continue
			}

			s, err := os.Stat(path)
			if err != nil {
				lastErr = err
				fmt.Println(locales.L.Get("errors.stat.failed", path, err.Error()))
				continue
			}

			if s.IsDir() {
				err = actions.UploadDirectory(*server, *domain, overwrite, dry, path, nil)
				if err != nil {
					lastErr = err
					fmt.Println(locales.L.Get("errors.upload.failed", path, err.Error()))
				}
			} else {
				if !dry {
					err = actions.UploadFile(*server, *domain, overwrite, path)
					if err != nil {
						lastErr = err
						fmt.Println(locales.L.Get("errors.upload.failed", path, err.Error()))
					}
				}
			}
		}

		cmd.SilenceUsage = true
		if lastErr == nil {
			fmt.Println(locales.L.Get("files.upload.success"))
		}

		return err
	},
	Args: cobra.MinimumNArgs(3),
}

func init() {
	SyncCmd.Flags().BoolP("overwrite", "f", false, locales.L.Get("files.upload.flag.overwrite"))
	SyncCmd.Flags().BoolP("dryRun", "n", false, locales.L.Get("files.upload.flag.dry-run"))
}
