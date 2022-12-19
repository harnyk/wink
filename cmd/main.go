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
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please specify a command")
		return
	}

	switch os.Args[1] {
	case "ls":
		{
			ls()
		}
	case "in":
		{
			in()
			ls()
		}
	case "out":
		{
			out()
			ls()
		}
	case "init":
		{
			initStore()
		}
	default:
		fmt.Println("Unknown command")
	}
}

func getConfigFileName() string {
	home := os.Getenv("HOME")
	return filepath.Join(home, ".wink", "config.json")
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
	var password string
	fmt.Scanln(&password)

	// password will be used as the key to encrypt the API key and employee ID

	store := cryptostore.NewCryptoStore(getConfigFileName())

	err := store.Store(cryptostore.CryptoStoreRecord{
		APIKey:     apiKey,
		EmployeeID: employeeID,
	}, password)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Your API key and employee ID have been saved")

	//lets try to load the record and display the API key (truncated) and employee ID
	loadedRecord, err := store.Load(password)
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

// ls lists all my check-ins
func ls() {
	// Get my check-ins
	checkInResult, err := GetCheckIns()
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

// in checks me in to work
func in() {
	checkInResult, err := GetCheckIns()
	if err != nil {
		fmt.Println(err)
		return
	}

	actions := TimeSheetToActionsList(checkInResult.Result[0])

	if !CanCheckIn(actions) {
		fmt.Println("You can't check in")
		return
	}

	fmt.Println("Checking in")

	slot := GetNextSlotName(checkInResult.Result[0])
	if slot == "" {
		fmt.Println("You can't check in")
		return
	}

	err = CheckInOut(slot)
	if err != nil {
		fmt.Println(err)
		return
	}
}

// out checks me out of work
func out() {
	checkInResult, err := GetCheckIns()
	if err != nil {
		fmt.Println(err)
		return
	}

	actions := TimeSheetToActionsList(checkInResult.Result[0])

	if !CanCheckOut(actions) {
		fmt.Println("You can't check out")
		return
	}

	fmt.Println("Checking out")

	slot := GetNextSlotName(checkInResult.Result[0])
	if slot == "" {
		fmt.Println("You can't check out")
		return
	}

	err = CheckInOut(slot)
	if err != nil {
		fmt.Println(err)
		return
	}
}

// POST https://api.peoplehr.net/Timesheet
/*

APIKey: "****"
Action: "UpdateTimesheet"
EmployeeId: "*****"
TimeOut1: "17:22"
TimesheetDate: "2022-12-16"
*/

func CheckInOut(slot string) error {
	date := getTodayYYYYMMDD()
	now := getNowHHMM()

	payload := map[string]string{
		"APIKey":        apiKey,
		"EmployeeId":    employeeId,
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

// The request body for the GetCheckIns endpoint:
// POST https://api.peoplehr.net/Timesheet
// {
// APIKey: ""****"
// Action: "GetTimesheetDetail"
// EmployeeId: "****"
// EndDate: "YYYY-MM-DD"
// StartDate: "YYYY-MM-DD"
// }
func GetCheckIns() (*GetTimesheetResponse, error) {
	checkInResponse := &GetTimesheetResponse{}

	date := getTodayYYYYMMDD()

	client := resty.New()
	_, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]string{
			"APIKey":     apiKey,
			"EmployeeId": employeeId,
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
