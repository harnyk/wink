package peopleapi

import "reflect"

// ActionType is the type of action: In or Out
type ActionType string

const (
	ActionTypeIn  ActionType = "In"
	ActionTypeOut ActionType = "Out"
)

type Action struct {
	Type ActionType
	Time string
}

func TimeSheetToActionsList(timeSheet *TimeSheet) []Action {
	var actions []Action

	fields := reflect.ValueOf(*timeSheet)

	for i := 0; i < fields.NumField(); i++ {
		field := fields.Field(i)

		if field.String() == "" {
			continue
		}

		if field.String() != "" {
			fieldName := fields.Type().Field(i).Name

			actionType := ActionTypeIn

			if fieldName[0:6] == "TimeIn" {
				actionType = ActionTypeIn
			} else if fieldName[0:7] == "TimeOut" {
				actionType = ActionTypeOut
			} else {
				continue
			}

			actions = append(actions, Action{Type: actionType, Time: field.String()})
		}
	}

	return actions
}

func CanCheckIn(actions []Action) bool {
	if len(actions) == 0 {
		return true
	}
	lastAction := actions[len(actions)-1]
	return lastAction.Type == ActionTypeOut
}

func CanCheckOut(actions []Action) bool {
	if len(actions) == 0 {
		return false
	}
	lastAction := actions[len(actions)-1]
	return lastAction.Type == ActionTypeIn
}

func GetNextSlotName(timeSheet TimeSheet) string {
	//find the first field that is empty and start with TimeIn or TimeOut
	fields := reflect.ValueOf(timeSheet)

	for i := 0; i < fields.NumField(); i++ {
		field := fields.Field(i)

		if field.String() == "" {
			fieldName := fields.Type().Field(i).Name

			if fieldName[0:6] == "TimeIn" || fieldName[0:7] == "TimeOut" {
				return fieldName
			}
		}
	}

	return ""
}
