package types

type TransformError struct {
	Message         string `json:"message"`
	RuleIndex       *int   `json:"ruleIndex,omitempty"`
	ActionIndex     *int   `json:"actionIndex,omitempty"`
	Key             string `json:"key,omitempty"`
}

type Action struct {
	Operation  string `json:"operation"`
	Value      string `json:"value"`
	ActionType string `json:"actionType"`
}

type PageCondition struct {
	SheetName           string `json:"sheetName"`
	IncludeFormulas     bool   `json:"includeFormulas"`
	NonEmptyValueRedact bool   `json:"nonEmptyValueRedact"`
}

/*
Examples:

redact by range -
operation: "range"
value: "C4:D9"

redact by value -
operation: "value"
value: "1.00%"

redact by text color -
operation: "textColor"
value: "0070C0"

redact by bg color -
operation: "bgColor"
value: "0070C0"

exclude column -
operation: "column"
value: "C"

exclude row -
operation: "row"
value: "4"
*/

/*
Request:

source: ...
configs: [
	{
		sheetName: '',
		actions: []
		includeFormulas: false
	},
	{
		sheetName: '',
		actions: []
		includeFormulas: false
	}
]
*/

type Rule struct {
	PageCondition PageCondition `json:"pageCondition"`
	Actions       []Action      `json:"actions"`
}