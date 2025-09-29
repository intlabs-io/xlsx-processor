package transform

import (
	"testing"
	"fmt"

	"xlsx-processor/pkg/testhelper"
	"xlsx-processor/pkg/types"

	"github.com/xuri/excelize/v2"
)

type testCase struct {
	name       string
	inputFile  string
	outputFile string
	action     *types.Action
	sheetName  string
}

func executeAction(filePath string, action *types.Action, sheetName string) (mutatedFile *excelize.File, err error) {
	file, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	actionExecutor := MakeActionExecutor(file, sheetName, true, action, 0, 0)
	transformErr := actionExecutor.Execute()
	if transformErr != nil {
		return nil, fmt.Errorf("%s", transformErr.Message)
	}

	return file, nil
}

func runTestCase(t *testing.T, tc testCase) {
	t.Helper()
	
	mutatedFile, err := executeAction(tc.inputFile, tc.action, tc.sheetName)
	if err != nil {
		t.Fatalf("failed to execute action: %v", err)
	}

	expectedFile, err := excelize.OpenFile(tc.outputFile)
	if err != nil {
		t.Fatalf("failed to load expected file: %v", err)
	}
	defer expectedFile.Close()

	testhelper.CompareSheet(t, expectedFile, mutatedFile, &tc.sheetName)
}

func TestActionExecutor(t *testing.T) {
	testCases := []testCase{
		{
			name:       "1 exclude column",
			inputFile:  "../assets/goldenFiles/testActionExecutor.xlsx",
			outputFile: "../assets/goldenFiles/testActionExecutor1Exclude.xlsx",
			action: &types.Action{
				ActionType: EXCLUDE,
				Operation:  COLUMN,
				Value:      "B",
			},
			sheetName: "Forecasting",
		},
		{
			name:       "2 exclude row",
			inputFile:  "../assets/goldenFiles/testActionExecutor.xlsx",
			outputFile: "../assets/goldenFiles/testActionExecutor2Exclude.xlsx",
			action: &types.Action{
				ActionType: EXCLUDE,
				Operation:  ROW,
				Value:      "2",
			},
			sheetName: "Forecasting",
		},
		{
			name:       "3 redact bg color",
			inputFile:  "../assets/goldenFiles/testActionExecutor.xlsx",
			outputFile: "../assets/goldenFiles/testActionExecutor3Redact.xlsx",
			action: &types.Action{
				ActionType: REDACT,
				Operation:  BG_COLOR,
				Value:      "4472C4",
			},
			sheetName: "Forecasting",
		},
		{
			name:       "4 redact text color",
			inputFile:  "../assets/goldenFiles/testActionExecutor.xlsx",
			outputFile: "../assets/goldenFiles/testActionExecutor4Redact.xlsx",
			action: &types.Action{
				ActionType: REDACT,
				Operation:  TEXT_COLOR,
				Value:      "0070C0",
			},
			sheetName: "Forecasting",
		},
		{
			name:       "5 redact column",
			inputFile:  "../assets/goldenFiles/testActionExecutor.xlsx",
			outputFile: "../assets/goldenFiles/testActionExecutor5Redact.xlsx",
			action: &types.Action{
				ActionType: REDACT,
				Operation:  COLUMN,
				Value:      "C",
			},
			sheetName: "Forecasting",
		},
		{
			name:       "6 redact row",
			inputFile:  "../assets/goldenFiles/testActionExecutor.xlsx",
			outputFile: "../assets/goldenFiles/testActionExecutor6Redact.xlsx",
			action: &types.Action{
				ActionType: REDACT,
				Operation:  ROW,
				Value:      "2",
			},
			sheetName: "Forecasting",
		},
		{
			name:       "7 redact range",
			inputFile:  "../assets/goldenFiles/testActionExecutor.xlsx",
			outputFile: "../assets/goldenFiles/testActionExecutor7Redact.xlsx",
			action: &types.Action{
				ActionType: REDACT,
				Operation:  RANGE,
				Value:      "B3:C5",
			},
			sheetName: "Forecasting",
		},
		{
			name:       "8 redact value",
			inputFile:  "../assets/goldenFiles/testActionExecutor.xlsx",
			outputFile: "../assets/goldenFiles/testActionExecutor8Redact.xlsx",
			action: &types.Action{
				ActionType: REDACT,
				Operation:  VALUE,
				Value:      "OFF",
			},
			sheetName: "Forecasting",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			runTestCase(t, tc)
		})
	}
}
