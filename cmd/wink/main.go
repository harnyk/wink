// This is a simple command-line tool that interacts with the PeopleHR API.
//It has three commands:
// ls - list all my check-ins
// in - check in to work
// out - check out of work
// init - ask for the API key, employee ID and password. Save them in a file using the crypto store.

package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/harnyk/wink/internal/cryptostore"
	api "github.com/harnyk/wink/internal/peopleapi"
	"golang.org/x/term"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please specify a command")
		return
	}

	if os.Args[1] == "init" {
		initStore()
		return
	}

	a, err := getAuth()
	if err != nil {
		fmt.Println(err)
		return
	}

	switch os.Args[1] {
	case "ls":
		{
			ls(a)
		}
	case "in":
		{
			in(a)
			ls(a)
		}
	case "out":
		{
			out(a)
			ls(a)
		}
	default:
		fmt.Println("Unknown command")
	}
}

func getConfigFileName() string {
	home := os.Getenv("HOME")
	return filepath.Join(home, ".wink", "secrets")
}

// init asks for the API key, employee ID and password. Save them in a file using the crypto store.
func initStore() {
	fmt.Println("Please enter your API key:")
	var apiKey string
	fmt.Scanln(&apiKey)

	fmt.Println("Please enter your employee ID:")
	var employeeID string
	fmt.Scanln(&employeeID)

	fmt.Println("Please enter your password:")
	password, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Println(err)
		return
	}

	// password will be used as the key to encrypt the API key and employee ID

	store := cryptostore.NewCryptoStore[Secrets](getConfigFileName())

	err = store.Store(Secrets{
		APIKey:     apiKey,
		EmployeeID: employeeID,
	}, string(password))

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Your API key and employee ID have been saved")

	//lets try to load the record and display the API key (truncated) and employee ID
	loadedRecord, err := store.Load(string(password))
	if err != nil {
		fmt.Println(err)
		return
	}

	maxAPIKeyLength := 5
	if len(loadedRecord.APIKey) < maxAPIKeyLength {
		maxAPIKeyLength = len(loadedRecord.APIKey)
	}

	fmt.Printf("Your API key is: %s...\n", loadedRecord.APIKey[:maxAPIKeyLength])
	fmt.Printf("Your employee ID is: %s\n", loadedRecord.EmployeeID)
}

func getAuth() (api.Auth, error) {
	store := cryptostore.NewCryptoStore[Secrets](getConfigFileName())

	fmt.Println("Please enter your password:")

	password, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return api.Auth{}, err
	}

	record, err := store.Load(string(password))
	if err != nil {
		return api.Auth{}, err
	}

	return api.Auth{
		APIKey:     record.APIKey,
		EmployeeID: record.EmployeeID,
	}, nil
}

// ls lists all my check-ins
func ls(a api.Auth) {

	client := api.NewClient(a)

	// Get my check-ins
	checkInResult, err := client.GetTimesheet()
	if err != nil {
		fmt.Println(err)
		return
	}

	if len(checkInResult.Result) == 0 {
		fmt.Println("No check-ins found")
		return
	}

	// Print my check-ins
	for _, timeSheet := range checkInResult.Result {
		fmt.Println(timeSheet.TimesheetDate)
		actions := api.TimeSheetToActionsList(timeSheet)
		for _, action := range actions {
			fmt.Printf(" - %s: %s\n", action.Type, action.Time)
		}
	}
}

func doAction(a api.Auth, action api.ActionType) error {
	client := api.NewClient(a)

	timeSheetResult, err := client.GetTimesheet()
	if err != nil {
		return err
	}
	currentTimesheet := api.TimeSheet{}
	if len(timeSheetResult.Result) > 0 {
		currentTimesheet = timeSheetResult.Result[0]
	}

	actions := api.TimeSheetToActionsList(currentTimesheet)

	switch action {
	case api.ActionTypeIn:
		{
			if !api.CanCheckIn(actions) {
				return fmt.Errorf("you can't check in")
			}
			fmt.Println("Checking in")
		}
	case api.ActionTypeOut:
		{
			if !api.CanCheckOut(actions) {
				return fmt.Errorf("you can't check out")
			}
			fmt.Println("Checking out")
		}
	}

	slot := api.GetNextSlotName(currentTimesheet)
	if slot == "" {
		return fmt.Errorf("timesheet is full")
	}

	if slot == "TimeIn1" {
		// create a new timesheet
		err := client.CreateNewTimesheet()
		if err != nil {
			return err
		}
	} else {
		err = client.CheckInOut(slot)
		if err != nil {
			return err
		}
	}

	return nil
}

// in checks me in to work
func in(a api.Auth) {
	err := doAction(a, api.ActionTypeIn)
	if err != nil {
		fmt.Println(err)
		return
	}
}

// out checks me out of work
func out(a api.Auth) {
	err := doAction(a, api.ActionTypeOut)
	if err != nil {
		fmt.Println(err)
		return
	}
}
