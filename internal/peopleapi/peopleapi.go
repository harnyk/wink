package peopleapi

import (
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
)

type Client interface {
	CreateNewTimesheet(time string) error
	CheckInOut(slot string, time string) error
	GetTimesheet(startDate time.Time, endDate time.Time) (*GetTimesheetResponse, error)
}

func NewClient(auth Auth) Client {
	return &client{
		auth: auth,
	}
}

type client struct {
	auth Auth
}

func (c *client) CreateNewTimesheet(time string) error {
	date := getTodayYYYYMMDD()
	var now string

	if time != "" {
		if !IsValidTime(time) {
			return fmt.Errorf("invalid time format")
		}
		now = time
	} else {
		now = getNowHHMM()
	}

	payload := map[string]string{
		"APIKey":        c.auth.APIKey,
		"EmployeeId":    c.auth.EmployeeID,
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

func (c *client) CheckInOut(slot string, time string) error {
	date := getTodayYYYYMMDD()

	var now string

	if time != "" {
		if !IsValidTime(time) {
			return fmt.Errorf("invalid time format")
		}
		now = time
	} else {
		now = getNowHHMM()
	}

	payload := map[string]string{
		"APIKey":        c.auth.APIKey,
		"EmployeeId":    c.auth.EmployeeID,
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

func (c *client) GetTimesheet(
	startDate time.Time,
	endDate time.Time,
) (*GetTimesheetResponse, error) {
	timeSheetResponse := &GetTimesheetResponse{}

	// date := getTodayYYYYMMDD()
	var startDateS string
	var endDateS string

	if startDate.IsZero() {
		startDateS = getTodayYYYYMMDD()
	} else {
		startDateS = startDate.Format("2006-01-02")
	}

	if endDate.IsZero() {
		endDateS = getTodayYYYYMMDD()
	} else {
		endDateS = endDate.Format("2006-01-02")
	}

	client := resty.New()
	_, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]string{
			"APIKey":     c.auth.APIKey,
			"EmployeeId": c.auth.EmployeeID,
			"Action":     "GetTimesheetDetail",
			"EndDate":    endDateS,
			"StartDate":  startDateS,
		}).
		SetResult(timeSheetResponse).
		Post("https://api.peoplehr.net/Timesheet")

	if err != nil {
		return nil, err
	}

	if timeSheetResponse.IsError {
		return nil, fmt.Errorf("server response: %s", timeSheetResponse.Message)
	}

	return timeSheetResponse, nil
}

func getTodayYYYYMMDD() string {
	return time.Now().Format("2006-01-02")
}

func getNowHHMM() string {
	return time.Now().Format("15:04")
}

func IsValidTime(t string) bool {
	_, err := time.Parse("15:04", t)
	return err == nil
}
