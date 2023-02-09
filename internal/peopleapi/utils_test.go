package peopleapi

import (
	"reflect"
	"testing"
)

func TestTimeSheetToActionsList(t *testing.T) {
	type args struct {
		timeSheet *TimeSheet
	}
	tests := []struct {
		name string
		args args
		want []Action
	}{
		{
			name: "TestTimeSheetToActionsList",
			args: args{
				timeSheet: &TimeSheet{
					TimeIn1:  "09:00",
					TimeOut1: "10:00",
					TimeIn2:  "11:00",
					TimeOut2: "12:00",
				},
			},
			want: []Action{
				{
					Type: ActionTypeIn,
					Time: "09:00",
				},
				{
					Type: ActionTypeOut,
					Time: "10:00",
				},
				{
					Type: ActionTypeIn,
					Time: "11:00",
				},
				{
					Type: ActionTypeOut,
					Time: "12:00",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TimeSheetToActionsList(tt.args.timeSheet); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TimeSheetToActionsList() = %v, want %v", got, tt.want)
			}
		})
	}
}
