package main

import (
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/harnyk/wink/internal/app"
	"github.com/harnyk/wink/internal/auth"
)

// this will be replaced in the goreleaser build
var version = "development"

type Command string

const (
	CmdLs     Command = "ls"
	CmdIn     Command = "in"
	CmdOut    Command = "out"
	CmdInit   Command = "init"
	CmdReport Command = "report"
)

func main() {
	fname, err := getConfigFileName()
	if err != nil {
		exitWithError(err)
	}
	authPrompt := auth.NewAuthPrompt(fname)

	a := app.NewApp(authPrompt, app.Version(version), app.ConfigFileName(fname))

	err = a.Run()
	if err != nil {
		exitWithError(err)
	}
}

func exitWithError(err error) {
	color.Red("▓▓▓▓ " + err.Error() + " ▓▓▓▓")
	os.Exit(1)
}

func getConfigFileName() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".wink", "secrets"), nil
}
