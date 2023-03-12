package report

import "math"

type TimesheetDailyTotalJSON struct {
	Date              string  `json:"date"`
	Hours             float64 `json:"hours"`
	IsComplete        bool    `json:"is_complete"`
	IsInvalidSequence bool    `json:"is_invalid_sequence"`
}

func NewTimesheetDailyTotalJSON(t *TimesheetDailyTotal) TimesheetDailyTotalJSON {
	return TimesheetDailyTotalJSON{
		Date:              t.Date.Format("2006-01-02"),
		Hours:             math.Round(t.Duration.Hours()*10) / 10,
		IsComplete:        t.IsComplete,
		IsInvalidSequence: t.IsInvalidSequence,
	}
}
