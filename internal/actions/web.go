// Copyright 1999-2022. Plesk International GmbH.

package actions

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
)

const defaultPort = "8080"

type Stack string

const (
	StackPHP     Stack = "PHP"
	StackJS      Stack = "JavaScript"
	StackRuby    Stack = "Ruby"
	StackStatic  Stack = "static"
	StackUnknown Stack = "unknown"
)

func DetectStack() Stack {
	if fileExists("Gemfile") && fileExists("config.ru") {
		return StackRuby
	}

	if fileExists("composer.json") {
		return StackPHP
	}

	if fileExists("package.json") {
		return StackJS
	}

	if fileExists("index.htm") || fileExists("index.html") {
		return StackStatic
	}

	return StackUnknown
}

func RunServer(stack Stack) error {
	port := defaultPort

	switch stack {
	case StackRuby:
		return runRuby(port)
	case StackPHP:
		return runPHP(port)
	case StackJS:
		return runJS()
	default:
		return runStatic(port)
	}
}

func fileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func runCommand(name string, arg ...string) error {
	command := exec.Command(name, arg...)
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	return command.Run()
}

func runPHP(port string) error {
	args := []string{"-S", "127.0.0.1:" + port}

	if fileExists("public") {
		args = append(args, "-t", "public")
	}

	return runCommand("php", args...)
}

func runJS() error {
	if fileExists("yarn.lock") {
		return runCommand("yarn", "start")
	}

	return runCommand("npm", "start")
}

func runRuby(port string) error {
	return runCommand("bundle", "exec", "rackup", "--port", port)
}

func runStatic(port string) error {
	fmt.Println("Static server started:", "http://127.0.0.1:"+port)
	return http.ListenAndServe(":"+port, http.FileServer(http.Dir(".")))
}
