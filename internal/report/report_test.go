package report_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/harnyk/wink/internal/peopleapi"
	"github.com/harnyk/wink/internal/report"
)

func TestCalculateHours(t *testing.T) {
	type args struct {
		dayTimeSheet *peopleapi.TimeSheet
	}
	tests := []struct {
		name    string
		args    args
		want    *report.TimesheetDailyTotal
		wantErr bool
	}{
		{
			name: "8 hours with 1 break",
			args: args{
				dayTimeSheet: &peopleapi.TimeSheet{
					TimesheetDate: "2020-01-01",
					TimeIn1:       "08:00:00",
					TimeOut1:      "12:00:00",
					TimeIn2:       "13:00:00",
					TimeOut2:      "17:00:00",
				},
			},
			want: &report.TimesheetDailyTotal{
				Date:              time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				Hours:             8 * time.Hour,
				IsComplete:        true,
				IsInvalidSequence: false,
			},
		},
		{
			name: "8 hours with 2 breaks",
			args: args{
				dayTimeSheet: &peopleapi.TimeSheet{
					TimesheetDate: "2020-01-01",
					TimeIn1:       "08:00:00",
					TimeOut1:      "12:00:00",
					TimeIn2:       "13:00:00",
					TimeOut2:      "14:00:00",
					TimeIn3:       "15:00:00",
					TimeOut3:      "18:00:00",
				},
			},
			want: &report.TimesheetDailyTotal{
				Date:              time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				Hours:             8 * time.Hour,
				IsComplete:        true,
				IsInvalidSequence: false,
			},
		},
		{
			name: "incomplete timesheet",
			args: args{
				dayTimeSheet: &peopleapi.TimeSheet{
					TimesheetDate: "2020-01-01",
					TimeIn1:       "08:00:00",
					TimeOut1:      "12:00:00",
					TimeIn2:       "13:00:00",
				},
			},
			want: &report.TimesheetDailyTotal{
				Date:              time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				Hours:             4 * time.Hour,
				IsComplete:        false,
				IsInvalidSequence: false,
			},
		},
		{
			name: "broken actions sequence",
			args: args{
				dayTimeSheet: &peopleapi.TimeSheet{
					TimesheetDate: "2020-01-01",
					TimeIn1:       "08:00:00",
					TimeIn2:       "09:00:00",
				},
			},
			want: &report.TimesheetDailyTotal{
				Date:              time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				Hours:             0 * time.Hour,
				IsComplete:        false,
				IsInvalidSequence: true,
			},
		},
		{
			name: "date parsing error",
			args: args{
				dayTimeSheet: &peopleapi.TimeSheet{
					TimesheetDate: "2020-01-01-01",
					TimeIn1:       "08:00:00",
					TimeOut1:      "12:00:00",
				},
			},
			wantErr: true,
		},
		{
			name: "time parsing error",
			args: args{
				dayTimeSheet: &peopleapi.TimeSheet{
					TimesheetDate: "2020-01-01",
					TimeIn1:       "08:XX:YY",
					TimeOut1:      "12:00:00",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := report.CalculateHours(tt.args.dayTimeSheet)
			if (err != nil) != tt.wantErr {
				t.Errorf("CalculateHours() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CalculateHours() got = %v, want %v", got, tt.want)
			}
		})
	}
}
