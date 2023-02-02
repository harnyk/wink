package report

import (
	"fmt"
	"time"

	"github.com/harnyk/wink/internal/peopleapi"
)

type TimesheetDailyTotal struct {
	Date       time.Time
	Hours      time.Duration
	IsComplete bool
}

func CalculateHours(dayTimeSheet *peopleapi.TimeSheet) (*TimesheetDailyTotal, error) {
	actionsList := peopleapi.TimeSheetToActionsList(dayTimeSheet)
	var totalHours time.Duration

	currentExpectedAction := peopleapi.ActionTypeIn

	currentTimeIn := time.Time{}
	err := error(nil)
	isComplete := false

	for _, action := range actionsList {
		if action.Type != currentExpectedAction {
			return nil, fmt.Errorf("invalid action sequence")
		}

		if currentExpectedAction == peopleapi.ActionTypeIn {
			currentTimeIn, err = time.Parse("15:04", action.Time)
			if err != nil {
				return nil, err
			}
		} else {
			timeOut, err := time.Parse("15:04", action.Time)
			if err != nil {
				return nil, err
			}

			totalHours += timeOut.Sub(currentTimeIn)
		}

		if currentExpectedAction == peopleapi.ActionTypeIn {
			currentExpectedAction = peopleapi.ActionTypeOut
			isComplete = false
		} else {
			currentExpectedAction = peopleapi.ActionTypeIn
			isComplete = true
		}
	}

	date, err := time.Parse("2006-01-02", dayTimeSheet.TimesheetDate)
	if err != nil {
		return nil, err
	}

	return &TimesheetDailyTotal{
		Date:       date,
		Hours:      totalHours,
		IsComplete: isComplete,
	}, nil
}
