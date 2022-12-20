package peopleapi

import (
	"time"

	"github.com/go-resty/resty/v2"
)

type Client interface {
	CreateNewTimesheet() error
	CheckInOut(slot string) error
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

func (c *client) CreateNewTimesheet() error {
	date := getTodayYYYYMMDD()
	now := getNowHHMM()

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

func (c *client) CheckInOut(slot string) error {
	date := getTodayYYYYMMDD()
	now := getNowHHMM()

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
