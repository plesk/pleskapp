// Copyright 1999-2021. Plesk International GmbH.

package main

import (
	"fmt"
	"os"

	"github.com/plesk/pleskapp/plesk/cmd"
	"github.com/plesk/pleskapp/plesk/config"
	"github.com/plesk/pleskapp/plesk/locales"
	"github.com/plesk/pleskapp/plesk/utils"
)

var (
	revision  string
	buildTime string
	version   string
)

func init() {
	cmd.Revision = revision
	cmd.BuildTime = buildTime
	cmd.Version = version
}

func main() {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(fmt.Errorf("failed obtaining current user home directory: %s", err))
	}

	path := home + "/.pleskrc"

	f, err := os.Open(path)
	if err != nil {
		config.New(nil)
	} else {
		defer f.Close()
		err = config.New(f)
		if err != nil {
			panic(err)
		}
	}

	exitCode := 0
	err = cmd.Execute()
	if err != nil {
		utils.Log.Error(locales.L.Get("errors.execution.failed.generic", err.Error()))
		exitCode = 1
	}

	f, err = os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600)
	if err != nil {
		fmt.Println("Could not open the configuration for saving")
		panic(err)
	}

	err = config.Save(f)
	if err != nil {
		fmt.Println("Could not save the configuration")
		panic(err)
	}

	os.Exit(exitCode)
}
