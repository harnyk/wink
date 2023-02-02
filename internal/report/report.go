package report

import (
	"time"

	"github.com/harnyk/wink/internal/peopleapi"
)

type TimesheetDailyTotal struct {
	Date              time.Time
	Hours             time.Duration
	IsComplete        bool
	IsInvalidSequence bool
}

func CalculateHours(dayTimeSheet *peopleapi.TimeSheet) (*TimesheetDailyTotal, error) {
	date, err := time.Parse("2006-01-02", dayTimeSheet.TimesheetDate)
	if err != nil {
		return nil, err
	}

	actionsList := peopleapi.TimeSheetToActionsList(dayTimeSheet)
	var totalHours time.Duration

	currentExpectedAction := peopleapi.ActionTypeIn

	currentTimeIn := time.Time{}
	isComplete := false

	for _, action := range actionsList {
		if action.Type != currentExpectedAction {
			return &TimesheetDailyTotal{
				Date:              date,
				Hours:             totalHours,
				IsComplete:        false,
				IsInvalidSequence: true,
			}, nil
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

	return &TimesheetDailyTotal{
		Date:              date,
		Hours:             totalHours,
		IsComplete:        isComplete,
		IsInvalidSequence: false,
	}, nil
}

func RenderDailyReport(timeSheetsList []peopleapi.TimeSheet) {

}
