package peopleapi

import (
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
)

type Client interface {
	CreateNewTimesheet(time string) error
	CheckInOut(slot string, time string) error
	GetTimesheet() (*GetTimesheetResponse, error)
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

func (c *client) GetTimesheet() (*GetTimesheetResponse, error) {
	checkInResponse := &GetTimesheetResponse{}

	date := getTodayYYYYMMDD()

	client := resty.New()
	_, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]string{
			"APIKey":     c.auth.APIKey,
			"EmployeeId": c.auth.EmployeeID,
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

func IsValidTime(t string) bool {
	_, err := time.Parse("15:04", t)
	return err == nil
}
