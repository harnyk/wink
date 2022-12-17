// This is a simple command-line tool that interacts with the PeopleHR API.
//It has three commands:
// ls - list all my check-ins
// in - check in to work
// out - check out of work

package main

import (
	"fmt"
	"os"
	"time"

	"github.com/go-resty/resty/v2"
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
	default:
		fmt.Println("Unknown command")
	}
}

// ls lists all my check-ins
func ls() {
	// Get my check-ins
	checkInResult, err := GetCheckIns()
	if err != nil {
		fmt.Println(err)
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

type TimeSheetResult struct {
	TimeIn1  string
	TimeOut1 string

	TimeIn2  string
	TimeOut2 string

	TimeIn3  string
	TimeOut3 string

	TimeIn4  string
	TimeOut4 string

	TimeIn5  string
	TimeOut5 string

	TimeIn6  string
	TimeOut6 string

	TimeIn7  string
	TimeOut7 string

	TimeIn8  string
	TimeOut8 string

	TimeIn9  string
	TimeOut9 string

	TimeIn10  string
	TimeOut10 string

	TimeIn11  string
	TimeOut11 string

	TimeIn12  string
	TimeOut12 string

	TimeIn13  string
	TimeOut13 string

	TimeIn14  string
	TimeOut14 string

	TimeIn15  string
	TimeOut15 string

	TimesheetDate string
}

type EditResponse struct {
	Message string `json:"Message"`
	Status  uint32 `json:"Status"`
	IsError bool   `json:"isError"`
}

type CheckInResponse struct {
	Message string            `json:"Message"`
	Result  []TimeSheetResult `json:"Result"`
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
func GetCheckIns() (*CheckInResponse, error) {
	checkInResponse := &CheckInResponse{}

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
