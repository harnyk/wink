package app

import (
	"log"

	"github.com/spf13/cobra"
)

type App interface {
	Run() error
}

type app struct {
}

func NewApp() App {
	return &app{}
}

func (a *app) Run() error {
	//parse the args using cobra

	//commands:
	// - ls
	// - in
	// - out
	// - init
	// - report
	// - version
	// - help

	rootCmd := &cobra.Command{
		Use:   "wink",
		Short: "Wink is a command line tool to check in and out of work",
		Long:  "Wink is a command line tool to check in and out of work",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	lsCmd := &cobra.Command{
		Use:   "ls",
		Short: "List all check-ins",
		Long:  "List all check-ins",
		RunE: func(cmd *cobra.Command, args []string) error {
			log.Println("ls")
			return nil
		},
	}

	rootCmd.AddCommand(lsCmd)

	return rootCmd.Execute()
}
