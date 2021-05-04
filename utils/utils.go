// Copyright 1999-2020. Plesk International GmbH.

package utils

import (
	"fmt"
	randn "math/rand"
	"strings"
	"syscall"
	"time"

	"github.com/plesk/pleskapp/plesk/types"
	"golang.org/x/crypto/ssh/terminal"
)

func FilterDomains(elements []types.Domain, filterOut string) ([]types.Domain, []types.Domain) {
	var keep []types.Domain
	var remove []types.Domain

	for _, element := range elements {
		if element.Name == filterOut {
			remove = append(remove, element)
		} else {
			keep = append(keep, element)
		}
	}

	return keep, remove
}

func FilterDatabases(elements []types.Database, filterOut string) ([]types.Database, []types.Database) {
	var keep []types.Database
	var remove []types.Database

	for _, element := range elements {
		if element.Name == filterOut {
			remove = append(remove, element)
		} else {
			keep = append(keep, element)
		}
	}

	return keep, remove
}

func FilterServers(elements []types.Server, filterOut string) ([]types.Server, []types.Server) {
	var keep []types.Server
	var remove []types.Server

	for _, element := range elements {
		if element.Host == filterOut {
			remove = append(remove, element)
		} else {
			keep = append(keep, element)
		}
	}

	return keep, remove
}

func GenPassword(length int) string {
	var charsets = [][]string{
		strings.Split("abcdefghijklmnopqrstuvwxyz", ""),
		strings.Split("0123456789", ""),
		// Symbol "&" may cause REST API to fail
		strings.Split("!@#$%^*()-=+_", ""),
	}

	randn.Seed(time.Now().Unix())
	var pw string = ""
	for i := 0; i < length; i++ {
		pw += charsets[i%3][randn.Intn(len(charsets[i%3]))]
	}

	pwb := []rune(pw)
	randn.Shuffle(len(pwb), func(i, j int) { pwb[i], pwb[j] = pwb[j], pwb[i] })

	return string(pwb)
}

var allowedChars = "abcdefghijklmnopqrstuvwxyz"

func GenUsername(length int) string {
	var charset = strings.Split("abcdefghijklmnopqrstuvwxyz", "")

	randn.Seed(time.Now().Unix())
	var pw string = ""
	for i := 0; i < length; i++ {
		pw += charset[randn.Intn(len(charset))]
	}

	pwb := []rune(pw)
	randn.Shuffle(len(pwb), func(i, j int) { pwb[i], pwb[j] = pwb[j], pwb[i] })

	return string(pwb)
}

func RequestPassword(reason string) (string, error) {
	fmt.Println(reason)

	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(bytePassword)), err
}

func strRev(s []rune) []rune {
	var rs []rune

	for i := len(s) - 1; i >= 0; i-- {
		rs = append(rs, s[i])
	}

	return rs
}

func StrSplitRN(s string, sep string, n int) []string {
	rev := strRev([]rune(s))
	spl := strings.SplitN(string(rev), sep, n)

	var sr []string
	for i := len(spl) - 1; i >= 0; i-- {
		sr = append(sr, string(strRev([]rune(spl[i]))))
	}

	return sr
}

func GetClientRootName(s string) (string, string) {
	spl := StrSplitRN(s, "/", 2)
	if len(spl) == 1 {
		return "/", spl[0]
	}

	return spl[0] + "/", spl[1]
}
