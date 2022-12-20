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
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/harnyk/wink/internal/cryptostore"
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

func getAuth() (Auth, error) {
	store := cryptostore.NewCryptoStore[Secrets](getConfigFileName())

	fmt.Println("Please enter your password:")

	password, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return Auth{}, err
	}

	record, err := store.Load(string(password))
	if err != nil {
		return Auth{}, err
	}

	return Auth{
		APIKey:     record.APIKey,
		EmployeeID: record.EmployeeID,
	}, nil
}

// ls lists all my check-ins
func ls(a Auth) {

	// Get my check-ins
	checkInResult, err := GetTimesheet(a)
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
		actions := TimeSheetToActionsList(timeSheet)
		for _, action := range actions {
			fmt.Printf(" - %s: %s\n", action.Type, action.Time)
		}
	}
}

func doAction(a Auth, action ActionType) error {
	timeSheetResult, err := GetTimesheet(a)
	if err != nil {
		return err
	}
	currentTimesheet := TimeSheet{}
	if len(timeSheetResult.Result) > 0 {
		currentTimesheet = timeSheetResult.Result[0]
	}

	actions := TimeSheetToActionsList(currentTimesheet)

	switch action {
	case ActionTypeIn:
		{
			if !CanCheckIn(actions) {
				return fmt.Errorf("you can't check in")
			}
			fmt.Println("Checking in")
		}
	case ActionTypeOut:
		{
			if !CanCheckOut(actions) {
				return fmt.Errorf("you can't check out")
			}
			fmt.Println("Checking out")
		}
	}

	slot := GetNextSlotName(currentTimesheet)
	if slot == "" {
		return fmt.Errorf("timesheet is full")
	}

	if slot == "TimeIn1" {
		// create a new timesheet
		err := CreateNewTimesheet(a)
		if err != nil {
			return err
		}
	} else {
		err = CheckInOut(a, slot)
		if err != nil {
			return err
		}
	}

	return nil
}

// in checks me in to work
func in(a Auth) {
	err := doAction(a, ActionTypeIn)
	if err != nil {
		fmt.Println(err)
		return
	}
}

// out checks me out of work
func out(a Auth) {
	err := doAction(a, ActionTypeOut)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func CreateNewTimesheet(auth Auth) error {
	date := getTodayYYYYMMDD()
	now := getNowHHMM()

	payload := map[string]string{
		"APIKey":        auth.APIKey,
		"EmployeeId":    auth.EmployeeID,
		"Action":        "CreateNewTimesheet",
		"TimesheetDate": date,
		"TimeIn1":       now,
	}

	client := resty.New()
	_, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(payload).
		SetResult(&EditResponse{}).
		Post("https://api.peoplehr.net/Timesheet")

	if err != nil {
		return err
	}

	return nil
}

func CheckInOut(auth Auth, slot string) error {
	date := getTodayYYYYMMDD()
	now := getNowHHMM()

	payload := map[string]string{
		"APIKey":        auth.APIKey,
		"EmployeeId":    auth.EmployeeID,
		"Action":        "UpdateTimesheet",
		"TimesheetDate": date,
	}

	payload[slot] = now

	client := resty.New()
	_, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(payload).
		SetResult(&EditResponse{}).
		Post("https://api.peoplehr.net/Timesheet")

	if err != nil {
		return err
	}

	return nil
}

func GetTimesheet(auth Auth) (*GetTimesheetResponse, error) {
	checkInResponse := &GetTimesheetResponse{}

	date := getTodayYYYYMMDD()

	client := resty.New()
	_, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]string{
			"APIKey":     auth.APIKey,
			"EmployeeId": auth.EmployeeID,
			"Action":     "GetTimesheetDetail",
			"EndDate":    date,
			"StartDate":  date,
		}).
		SetResult(checkInResponse).
		Post("https://api.peoplehr.net/Timesheet")

	if err != nil {
		return nil, err
	}

	return checkInResponse, nil
}

func getTodayYYYYMMDD() string {
	return time.Now().Format("2006-01-02")
}

func getNowHHMM() string {
	return time.Now().Format("15:04")
}
