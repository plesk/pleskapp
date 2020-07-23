// Copyright 1999-2020. Plesk International GmbH. All rights reserved.

package features

import (
	"encoding/json"
	"errors"
	"strings"

	"git.plesk.ru/projects/SBX/repos/pleskapp/locales"
)

type Feature string

const (
	Php70 Feature = "php70"
	Php71 Feature = "php71"
	Php72 Feature = "php72"
	Php73 Feature = "php73"
	Php74 Feature = "php74"
	Php   Feature = "php74"
	Nginx Feature = "nginx"
)

type FeatureProvider struct {
	IsWindows bool
}

func GetFeatureByString(specifier string) *Feature {
	var t Feature
	switch specifier {
	case "php":
		fallthrough
	case "php74":
		t = Php74
		return &t
	case "php73":
		t = Php73
		return &t
	case "php72":
		t = Php72
		return &t
	case "php71":
		t = Php71
		return &t
	case "php70":
		t = Php70
		return &t
	case "nginx":
		t = Nginx
		return &t
	}

	return nil
}

func (f FeatureProvider) GetFeaturePackage(domain string, feature Feature) ([]byte, error) {
	var packet map[string][]string = map[string][]string{"params": []string{}, "env": []string{}}

	switch feature {
	case Php74:
		fallthrough
	case Php73:
		fallthrough
	case Php72:
		fallthrough
	case Php71:
		fallthrough
	case Php70:
		if !f.IsWindows {
			packet["params"] = []string{"--update", domain, "-php", "true", "-php_handler_id", "plesk-" + string(feature) + "-fpm"}
		} else {
			ver := strings.Split(strings.TrimPrefix(string(feature), "php"), "")
			packet["params"] = []string{"--update", domain, "-php", "true", "-php_handler_id", "fastcgi-" + ver[0] + "." + ver[1]}
		}
	case Nginx:
		if f.IsWindows {
			return []byte{}, errors.New(locales.L.Get("errors.feature.not.supported", string(feature)))
		}
		packet["params"] = []string{"--update", domain, "-nginx-serve-php", "true"}
	default:
		return []byte{}, errors.New(locales.L.Get("errors.feature.unknown", string(feature)))
	}

	jsonPacket, err := json.Marshal(packet)
	return jsonPacket, err
}
