package app

import (
	"fmt"
	"time"

	"github.com/fatih/color"
	"github.com/harnyk/wink/internal/auth"
	"github.com/harnyk/wink/internal/cryptostore"
	"github.com/harnyk/wink/internal/easteregg"
	"github.com/harnyk/wink/internal/entities"
	"github.com/harnyk/wink/internal/peopleapi"
	"github.com/harnyk/wink/internal/report"
	"github.com/harnyk/wink/internal/ui"
	"github.com/jinzhu/now"

	"github.com/spf13/cobra"
)

type App interface {
	Run() error
}

type app struct {
	authPrompt     auth.AuthPrompt
	version        Version
	configFileName ConfigFileName
}

func NewApp(
	authPrompt auth.AuthPrompt,
	appVersion Version,
	configFileName ConfigFileName,
) App {
	return &app{
		authPrompt:     authPrompt,
		version:        appVersion,
		configFileName: configFileName,
	}
}

func (a *app) Run() error {
	// 	//seed a random number generator
	easteregg.Seed()

	rootCmd := &cobra.Command{
		Use:   "wink",
		Short: "Wink is a command line tool to check in and out of work",
		Long:  "Wink is a command line tool to check in and out of work",
		RunE: func(cmd *cobra.Command, args []string) error {
			if cmd.Flag("version").Value.String() == "true" {
				return a.doVersion()
			}

			return cmd.Help()
		},
	}
	rootCmd.Flags().BoolP("version", "v", false, "Print the version number of wink")

	lsCmd := &cobra.Command{
		Use:   "ls",
		Short: "List all my check-ins",
		Long:  "List all my check-ins",
		RunE: func(cmd *cobra.Command, args []string) error {
			return a.doList()
		},
	}

	inCmd := &cobra.Command{
		Use:   "in",
		Short: "Check in to work",
		Long:  "Check in to work",
		RunE: func(cmd *cobra.Command, args []string) error {
			var timeArg string
			if len(args) > 0 {
				timeArg = args[0]
			}

			return a.doCheckInOut(timeArg, peopleapi.ActionTypeIn)
		},
	}

	outCmd := &cobra.Command{
		Use:   "out",
		Short: "Check out of work",
		Long:  "Check out of work",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var timeArg string
			if len(args) > 0 {
				timeArg = args[0]
			}

			return a.doCheckInOut(timeArg, peopleapi.ActionTypeOut)
		},
	}

	initCmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize wink",
		Long:  "Initialize wink",
		RunE: func(cmd *cobra.Command, args []string) error {
			return a.doInit()
		},
	}

	reportCmd := &cobra.Command{
		Use:   "report",
		Short: "Generate a report",
		Long:  "Generate a report",
		RunE: func(cmd *cobra.Command, args []string) error {
			var start time.Time
			var end time.Time
			var err error

			if cmd.Flag("start").Value.String() == "" {
				start = now.BeginningOfMonth()
			} else {
				start, err = time.Parse("2006-01-02", cmd.Flag("start").Value.String())
				if err != nil {
					return err
				}
			}

			if cmd.Flag("end").Value.String() == "" {
				end = time.Now()
			} else {
				end, err = time.Parse("2006-01-02", cmd.Flag("end").Value.String())
				if err != nil {
					return err
				}
			}

			return a.doReport(start, end)

		},
	}
	reportCmd.Flags().StringP("start", "s", "", "Start date")
	reportCmd.Flags().StringP("end", "e", "", "End date")

	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version number of wink",
		Long:  "Print the version number of wink",
		RunE: func(cmd *cobra.Command, args []string) error {
			return a.doVersion()
		},
	}

	rootCmd.AddCommand(inCmd)
	rootCmd.AddCommand(outCmd)
	rootCmd.AddCommand(lsCmd)
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(reportCmd)
	rootCmd.AddCommand(versionCmd)

	return rootCmd.Execute()
}

func (a *app) doCheckInOut(timeFlag string, action peopleapi.ActionType) error {
	var checkInTime time.Time
	var err error

	if timeFlag == "" {
		checkInTime = time.Now()
	} else {
		checkInTime, err = time.Parse("15:04", timeFlag)
		if err != nil {
			return err
		}
	}

	au, err := a.authPrompt.Get()
	if err != nil {
		return err
	}

	if err = checkInOut(au, action, checkInTime); err != nil {
		return err
	}

	switch action {
	case peopleapi.ActionTypeIn:
		{
			printSuccess(fmt.Sprintf("Checked in at %s", checkInTime.Format("15:04")))
			fmt.Println(easteregg.GetRandomCheckinPhrase(0.5))
		}
	case peopleapi.ActionTypeOut:
		{
			printSuccess(fmt.Sprintf("Checked out at %s", checkInTime.Format("15:04")))
			fmt.Println(easteregg.GetRandomCheckoutPhrase(0.5))
		}
	}

	return nil

}

func (a *app) doList() error {
	authData, err := a.authPrompt.Get()
	if err != nil {
		return err
	}

	client := peopleapi.NewClient(authData)

	// Get my check-ins
	checkInResult, err := client.GetTimesheet(time.Time{}, time.Time{})
	if err != nil {
		return err

	}

	fmt.Println()

	if len(checkInResult.Result) == 0 {
		return fmt.Errorf("no check-ins found")
	}

	// Print my check-ins
	for _, timeSheet := range checkInResult.Result {
		fmt.Println(timeSheet.TimesheetDate)
		actions := peopleapi.TimeSheetToActionsList(&timeSheet)
		for _, action := range actions {
			fmt.Printf(" - %s:\t%s\n", action.Type, action.Time)
		}
	}

	return nil
}

func (a *app) doInit() error {

	u := ui.NewUI()

	apiKey, err := u.AskString("Please enter your API key:")
	if err != nil {
		return err
	}

	employeeID, err := u.AskString("Please enter your employee ID:")
	if err != nil {
		return err
	}

	password, err := u.AskPassword("Please enter a password to encrypt your API key and employee ID:")
	if err != nil {
		return err
	}

	store := cryptostore.NewCryptoStore[entities.Secrets](string(a.configFileName))

	err = store.Store(entities.Secrets{
		APIKey:     apiKey,
		EmployeeID: employeeID,
	}, string(password))

	if err != nil {
		return err
	}

	//lets try to load the record and display the API key (truncated) and employee ID
	loadedRecord, err := store.Load(string(password))
	if err != nil {
		return err
	}

	maxAPIKeyLength := 5
	if len(loadedRecord.APIKey) < maxAPIKeyLength {
		maxAPIKeyLength = len(loadedRecord.APIKey)
	}

	fmt.Printf("Your API key is: %s...\n", loadedRecord.APIKey[:maxAPIKeyLength])
	fmt.Printf("Your employee ID is: %s\n", loadedRecord.EmployeeID)

	printSuccess("Successfully initialized wink")

	return nil

}

func (a *app) doReport(timeStart, timeEnd time.Time) error {
	authData, err := a.authPrompt.Get()
	if err != nil {
		return err
	}

	client := peopleapi.NewClient(authData)

	reportData, err := client.GetTimesheet(timeStart, timeEnd)
	if err != nil {
		return err
	}

	fmt.Println()

	reportStr := report.RenderDailyReport(timeStart, timeEnd, reportData.Result)

	fmt.Println(reportStr)

	return nil
}

func (a *app) doVersion() error {
	fmt.Println(a.version)
	return nil
}

// ------------------------

func printSuccess(message string) {
	color.Green("▓▓▓▓ " + message + " ▓▓▓▓")
}

func checkInOut(authData peopleapi.Auth, action peopleapi.ActionType, checkInTime time.Time) error {

	timeStr := checkInTime.Format("15:04")

	client := peopleapi.NewClient(authData)

	timeSheetResult, err := client.GetTimesheet(time.Time{}, time.Time{})
	if err != nil {
		return err
	}
	currentTimesheet := peopleapi.TimeSheet{}
	if len(timeSheetResult.Result) > 0 {
		currentTimesheet = timeSheetResult.Result[0]
	}

	actions := peopleapi.TimeSheetToActionsList(&currentTimesheet)

	switch action {
	case peopleapi.ActionTypeIn:
		{
			if !peopleapi.CanCheckIn(actions) {
				return fmt.Errorf("you can't check in")
			}
			fmt.Println("Checking in")
		}
	case peopleapi.ActionTypeOut:
		{
			if !peopleapi.CanCheckOut(actions) {
				return fmt.Errorf("you can't check out")
			}
			fmt.Println("Checking out")
		}
	}

	slot := peopleapi.GetNextSlotName(currentTimesheet)
	if slot == "" {
		return fmt.Errorf("timesheet is full")
	}

	if slot == "TimeIn1" {
		// create a new timesheet
		err := client.CreateNewTimesheet(timeStr)
		if err != nil {
			return err
		}
	} else {
		err = client.CheckInOut(slot, timeStr)
		if err != nil {
			return err
		}
	}

	return nil
}
