package main

import "encoding/json"

type TimeSheet struct {
	TimeIn1  string
	TimeOut1 string

	TimeIn2  string
	TimeOut2 string

	TimeIn3  string
	TimeOut3 string

	TimeIn4  string
	TimeOut4 string

	TimeIn5  string
	TimeOut5 string

	TimeIn6  string
	TimeOut6 string

	TimeIn7  string
	TimeOut7 string

	TimeIn8  string
	TimeOut8 string

	TimeIn9  string
	TimeOut9 string

	TimeIn10  string
	TimeOut10 string

	TimeIn11  string
	TimeOut11 string

	TimeIn12  string
	TimeOut12 string

	TimeIn13  string
	TimeOut13 string

	TimeIn14  string
	TimeOut14 string

	TimeIn15  string
	TimeOut15 string

	TimesheetDate string
}

type EditResponse struct {
	Message string `json:"Message"`
	Status  uint32 `json:"Status"`
	IsError bool   `json:"isError"`
}

type GetTimesheetResponse struct {
	Message string      `json:"Message"`
	Result  []TimeSheet `json:"Result"`
}

//Unmarshal JSON into GetTimesheetResponse allowing for empty strings at Result field

func (gtsr *GetTimesheetResponse) UnmarshalJSON(data []byte) error {
	type Alias GetTimesheetResponse
	aux := &struct {
		Result json.RawMessage `json:"Result"`
		*Alias
	}{
		Alias: (*Alias)(gtsr),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	// if aux.Result is an empty string, we need to set it to an empty array
	if string(aux.Result) == `""` {
		aux.Result = []byte("[]")
	}
	return json.Unmarshal(aux.Result, &gtsr.Result)
}
