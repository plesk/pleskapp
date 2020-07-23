// Copyright 1999-2020. Plesk International GmbH. All rights reserved.

package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	"git.plesk.ru/projects/SBX/repos/pleskapp/types"
)

var localConfig config = config{}

type lconfig struct {
	config types.App
}

func LoadLocal(f *os.File) (*lconfig, error) {
	var c types.App
	if f != nil {
		buf := bytes.NewBuffer([]byte{})
		_, err := buf.ReadFrom(f)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(buf.Bytes(), &c)
		if err != nil {
			return nil, err
		}
	}

	return &lconfig{config: c}, nil
}

func (c *lconfig) Save(f *os.File) error {
	if f == nil {
		return fmt.Errorf("Cannot save configuration")
	}

	str, err := json.Marshal(c.config)
	if err != nil {
		return err
	}

	err = f.Truncate(0)
	if err != nil {
		return err
	}

	_, err = f.Write(str)
	return err
}

func (c *lconfig) SetTargetPath(p string) {
	c.config.TargetPath = p
}

func (c *lconfig) SetFeatures(f []string) {
	c.config.Features = f
}

func (c *lconfig) GetTargetPath() string {
	return c.config.TargetPath
}

func (c *lconfig) GetFeatures() []string {
	return c.config.Features
}
