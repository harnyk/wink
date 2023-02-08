package main

import "github.com/harnyk/wink/internal/app"

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
	a := app.NewApp()

	err := a.Run()
	if err != nil {
		panic(err)
	}
}

// func main() {
// 	usage := `Wink - command line time tracker.

// Usage:
//   wink ls
//   wink in [<time>]
//   wink out [<time>]
//   wink init
//   wink report [--start=<start>] [--end=<end>]
//   wink --version

// Commands:
//   ls   - list all my check-ins
//   in   - check in to work
//   out  - check out of work
//   init - setup the API key, and employee ID. Encrypt them using a password
//   report - generate a report for the current month
// `

// 	//seed a random number generator
// 	easteregg.Seed()

// 	arguments, _ := docopt.ParseArgs(usage, nil, version)

// 	command, err := getCommand(arguments)
// 	if err != nil {
// 		fmt.Println(usage)
// 		color.Red(err.Error())
// 		return
// 	}

// 	configFile, err := getConfigFileName()
// 	if err != nil {
// 		fmt.Println(usage)
// 		color.Red(err.Error())
// 		return
// 	}
// 	authPrompt := auth.NewAuthPrompt(configFile)

// 	switch Command(command) {
// 	case CmdInit:
// 		{
// 			if err := initStore(); err != nil {
// 				fmt.Println(usage)
// 				color.Red(err.Error())
// 				return
// 			}
// 		}
// 	case CmdLs:
// 		{
// 			if err := ls(authPrompt); err != nil {
// 				fmt.Println(usage)
// 				color.Red(err.Error())
// 				return
// 			}
// 		}
// 	case CmdIn:
// 		{
// 			time, err := getOptionalTime(arguments)
// 			if err != nil {
// 				fmt.Println(usage)
// 				color.Red(err.Error())
// 				return
// 			}

// 			if err := in(authPrompt, time); err != nil {
// 				fmt.Println(usage)
// 				color.Red(err.Error())
// 				return
// 			}

// 			color.Green("▓▓▓▓ Checked in ▓▓▓▓")

// 			eaphrase := easteregg.GetRandomCheckinPhrase(0.5)
// 			fmt.Println(eaphrase)
// 		}
// 	case CmdOut:
// 		{
// 			time, err := getOptionalTime(arguments)
// 			if err != nil {
// 				fmt.Println(usage)
// 				color.Red(err.Error())
// 				return
// 			}

// 			if err := out(authPrompt, time); err != nil {
// 				fmt.Println(usage)
// 				color.Red(err.Error())
// 				return
// 			}

// 			color.Green("▓▓▓▓ Checked out ▓▓▓▓")

// 			eaphrase := easteregg.GetRandomCheckoutPhrase(0.5)
// 			fmt.Println(eaphrase)
// 		}
// 	case CmdReport:
// 		{

// 			start, ok := arguments["--start"].(string)
// 			if !ok {
// 				start = now.BeginningOfMonth().Format("2006-01-02")
// 			}

// 			end, ok := arguments["--end"].(string)
// 			if !ok {
// 				end = now.EndOfMonth().Format("2006-01-02")
// 			}

// 			if err := doReport(authPrompt, start, end); err != nil {
// 				fmt.Println(usage)
// 				color.Red(err.Error())
// 				return
// 			}
// 		}
// 	default:
// 		{
// 			fmt.Println(usage)
// 			color.Red("Unknown command")
// 			return
// 		}
// 	}

// }

// func getCommand(arguments docopt.Opts) (Command, error) {
// 	commands := []Command{CmdLs, CmdIn, CmdOut, CmdInit, CmdReport}

// 	for _, command := range commands {
// 		commandSet, err := arguments.Bool(string(command))
// 		if err != nil {
// 			return "", err
// 		}
// 		if commandSet {
// 			return command, nil
// 		}
// 	}

// 	return "", fmt.Errorf("Unknown command")
// }

// func getOptionalTime(args docopt.Opts) (string, error) {
// 	time, err := args.String("<time>")
// 	if err != nil {
// 		return "", nil
// 	}

// 	if !api.IsValidTime(time) {
// 		return "", fmt.Errorf("Invalid time format. Please use 24h format, e.g. 12:00, 15:30")
// 	}

// 	return time, nil
// }

// func getConfigFileName() (string, error) {
// 	home, err := os.UserHomeDir()
// 	if err != nil {
// 		fmt.Println(err)
// 		return "", err
// 	}
// 	return filepath.Join(home, ".wink", "secrets"), nil
// }

// // init asks for the API key, employee ID and password. Save them in a file using the crypto store.
// func initStore() error {
// 	u := ui.NewUI()

// 	apiKey, err := u.AskString("Please enter your API key:")
// 	if err != nil {
// 		return err
// 	}

// 	employeeID, err := u.AskString("Please enter your employee ID:")
// 	if err != nil {
// 		return err
// 	}

// 	password, err := u.AskPassword("Please enter a password to encrypt your API key and employee ID:")
// 	if err != nil {
// 		return err
// 	}

// 	confPath, err := getConfigFileName()
// 	if err != nil {
// 		return err
// 	}

// 	store := cryptostore.NewCryptoStore[entities.Secrets](confPath)

// 	err = store.Store(entities.Secrets{
// 		APIKey:     apiKey,
// 		EmployeeID: employeeID,
// 	}, string(password))

// 	if err != nil {
// 		return err
// 	}

// 	fmt.Println("Your API key and employee ID have been saved")

// 	//lets try to load the record and display the API key (truncated) and employee ID
// 	loadedRecord, err := store.Load(string(password))
// 	if err != nil {
// 		return err
// 	}

// 	maxAPIKeyLength := 5
// 	if len(loadedRecord.APIKey) < maxAPIKeyLength {
// 		maxAPIKeyLength = len(loadedRecord.APIKey)
// 	}

// 	fmt.Printf("Your API key is: %s...\n", loadedRecord.APIKey[:maxAPIKeyLength])
// 	fmt.Printf("Your employee ID is: %s\n", loadedRecord.EmployeeID)

// 	return nil
// }

// // ls lists all my check-ins
// func ls(authPrompt auth.AuthPrompt) error {

// 	a, err := authPrompt.Get()
// 	if err != nil {
// 		return err
// 	}

// 	client := api.NewClient(a)

// 	// Get my check-ins
// 	checkInResult, err := client.GetTimesheet(time.Time{}, time.Time{})
// 	if err != nil {
// 		return err

// 	}

// 	fmt.Println()

// 	if len(checkInResult.Result) == 0 {
// 		return fmt.Errorf("no check-ins found")
// 	}

// 	// Print my check-ins
// 	for _, timeSheet := range checkInResult.Result {
// 		fmt.Println(timeSheet.TimesheetDate)
// 		actions := api.TimeSheetToActionsList(&timeSheet)
// 		for _, action := range actions {
// 			fmt.Printf(" - %s:\t%s\n", action.Type, action.Time)
// 		}
// 	}

// 	return nil
// }

// // time is optional. If not provided, it will use the current time
// func checkInOut(a api.Auth, action api.ActionType, timeStr string) error {
// 	client := api.NewClient(a)

// 	timeSheetResult, err := client.GetTimesheet(time.Time{}, time.Time{})
// 	if err != nil {
// 		return err
// 	}
// 	currentTimesheet := api.TimeSheet{}
// 	if len(timeSheetResult.Result) > 0 {
// 		currentTimesheet = timeSheetResult.Result[0]
// 	}

// 	actions := api.TimeSheetToActionsList(&currentTimesheet)

// 	switch action {
// 	case api.ActionTypeIn:
// 		{
// 			if !api.CanCheckIn(actions) {
// 				return fmt.Errorf("you can't check in")
// 			}
// 			fmt.Println("Checking in")
// 		}
// 	case api.ActionTypeOut:
// 		{
// 			if !api.CanCheckOut(actions) {
// 				return fmt.Errorf("you can't check out")
// 			}
// 			fmt.Println("Checking out")
// 		}
// 	}

// 	slot := api.GetNextSlotName(currentTimesheet)
// 	if slot == "" {
// 		return fmt.Errorf("timesheet is full")
// 	}

// 	if slot == "TimeIn1" {
// 		// create a new timesheet
// 		err := client.CreateNewTimesheet(timeStr)
// 		if err != nil {
// 			return err
// 		}
// 	} else {
// 		err = client.CheckInOut(slot, timeStr)
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }

// // in checks me in to work
// func in(authPrompt auth.AuthPrompt, time string) error {
// 	a, err := authPrompt.Get()
// 	if err != nil {
// 		return err
// 	}

// 	if err = checkInOut(a, api.ActionTypeIn, time); err != nil {
// 		return err
// 	}

// 	return nil
// }

// // out checks me out of work
// func out(authPrompt auth.AuthPrompt, time string) error {
// 	a, err := authPrompt.Get()
// 	if err != nil {
// 		return err
// 	}

// 	if err := checkInOut(a, api.ActionTypeOut, time); err != nil {
// 		return err
// 	}

// 	return nil
// }

// func doReport(authPrompt auth.AuthPrompt, start string, end string) error {
// 	a, err := authPrompt.Get()
// 	if err != nil {
// 		return err
// 	}

// 	client := api.NewClient(a)

// 	timeStart, err := time.Parse("2006-01-02", start)
// 	if err != nil {
// 		return err
// 	}

// 	timeEnd, err := time.Parse("2006-01-02", end)
// 	if err != nil {
// 		return err
// 	}

// 	reportData, err := client.GetTimesheet(timeStart, timeEnd)
// 	if err != nil {
// 		return err
// 	}

// 	fmt.Println()

// 	reportStr := report.RenderDailyReport(timeStart, timeEnd, reportData.Result)

// 	fmt.Println(reportStr)

// 	return nil
// }
