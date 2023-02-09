package report

import (
	"fmt"
	"strings"
	"time"

	"github.com/fatih/color"

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
			currentTimeIn, err = time.Parse("15:04:05", action.Time)
			if err != nil {
				return nil, err
			}
		} else {
			timeOut, err := time.Parse("15:04:05", action.Time)
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

func RenderDailyReport(dateStart time.Time, dateEnd time.Time, timeSheets []peopleapi.TimeSheet) string {
	dimmed := color.New(color.Faint).SprintFunc()
	bold := color.New(color.Bold).SprintFunc()

	perDateTotals := make(map[string]TimesheetDailyTotal)

	for _, timeSheet := range timeSheets {
		timesheetDailyTotal, err := CalculateHours(&timeSheet)
		if err != nil {
			continue
		}

		perDateTotals[timesheetDailyTotal.Date.Format("2006-01-02")] = *timesheetDailyTotal
	}

	var report strings.Builder

	report.WriteString(dimmed("-----------------------------------------------\n"))

	report.WriteString(color.CyanString("# Daily report"))
	report.WriteString("\n")

	report.WriteString(dimmed("From : "))
	report.WriteString(dateStart.Format("02-Jan-2006"))
	report.WriteString("\n")
	report.WriteString(dimmed("To   : "))
	report.WriteString(dateEnd.Format("02-Jan-2006"))

	report.WriteString("\n")
	report.WriteString("\n")

	for date := dateStart; date.Before(dateEnd); date = date.AddDate(0, 0, 1) {
		report.WriteString(date.Format("02-Jan"))
		report.WriteString(" ")
		report.WriteString(renderWeekDay(date))
		report.WriteString(": ")

		timesheetDailyTotal, ok := perDateTotals[date.Format("2006-01-02")]
		if !ok {
			report.WriteString(dimmed("-"))
			report.WriteString("\n")
			continue
		}

		if timesheetDailyTotal.IsInvalidSequence {
			report.WriteString(color.RedString("Invalid sequence"))
			report.WriteString("\n")
			continue
		}

		report.WriteString(bold(fmt.Sprintf("%.1fh", timesheetDailyTotal.Hours.Hours())))
		if !timesheetDailyTotal.IsComplete {
			report.WriteString(" ")
			report.WriteString(color.YellowString("(incomplete)"))
		}
		report.WriteString("\n")
	}

	report.WriteString(dimmed("\n-----------------------------------------------\n"))

	return report.String()
}

func renderWeekDay(date time.Time) string {

	str := date.Format("Mon")

	if date.Weekday() == time.Sunday || date.Weekday() == time.Saturday {
		return color.RedString(str)
	}

	return date.Format("Mon")
}
