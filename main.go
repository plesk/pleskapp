// Copyright 1999-2020. Plesk International GmbH. All rights reserved.

package main

import (
	"fmt"
	"os"

	"git.plesk.ru/projects/SBX/repos/pleskapp/cmd"
	"git.plesk.ru/projects/SBX/repos/pleskapp/config"
)

func main() {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(fmt.Errorf("Failed obtaining current user home directory: %s", err))
	}

	path := home + "/.pleskapprc"

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
		fmt.Println(err.Error())
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
