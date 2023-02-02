package report

import (
	"testing"
	"time"

	"github.com/harnyk/wink/internal/peopleapi"
)

func TestCalculateHours(t *testing.T) {
	type args struct {
		dayTimeSheet *peopleapi.TimeSheet
	}
	tests := []struct {
		name    string
		args    args
		want    *TimesheetDailyTotal
		wantErr bool
	}{
		{
			name: "8 hours with 1 break",
			args: args{
				dayTimeSheet: &peopleapi.TimeSheet{
					TimesheetDate: "2020-01-01",
					TimeIn1:       "08:00",
					TimeOut1:      "12:00",
					TimeIn2:       "13:00",
					TimeOut2:      "17:00",
				},
			},
			want: &TimesheetDailyTotal{
				Date:       time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				Hours:      8 * time.Hour,
				IsComplete: true,
			},
		},
		{
			name: "8 hours with 2 breaks",
			args: args{
				dayTimeSheet: &peopleapi.TimeSheet{
					TimesheetDate: "2020-01-01",
					TimeIn1:       "08:00",
					TimeOut1:      "12:00",
					TimeIn2:       "13:00",
					TimeOut2:      "14:00",
					TimeIn3:       "15:00",
					TimeOut3:      "18:00",
				},
			},
			want: &TimesheetDailyTotal{
				Date:       time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				Hours:      8 * time.Hour,
				IsComplete: true,
			},
		},
		{
			name: "incomplete timesheet",
			args: args{
				dayTimeSheet: &peopleapi.TimeSheet{
					TimesheetDate: "2020-01-01",
					TimeIn1:       "08:00",
					TimeOut1:      "12:00",
					TimeIn2:       "13:00",
				},
			},
			want: &TimesheetDailyTotal{
				Date:       time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				Hours:      4 * time.Hour,
				IsComplete: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CalculateHours(tt.args.dayTimeSheet)
			if (err != nil) != tt.wantErr {
				t.Errorf("CalculateHours() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CalculateHours() = %v, want %v", got, tt.want)
			}
		})
	}
}
