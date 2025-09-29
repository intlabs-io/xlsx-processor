package transform

import (
	"fmt"
	"testing"
	"xlsx-processor/pkg/types"

	"github.com/go-playground/assert/v2"
	"github.com/xuri/excelize/v2"
)

type formulaTestCase struct {
	name            string
	inputFile       string
	includeFormulas bool
	expectedFormula string
	sheetName       string
	cellReference   string
}

func executeFormulaRule(filePath string, includeFormulas bool, sheetName string) (file *excelize.File, err error) {
	file, err = excelize.OpenFile(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	rules := []types.Rule{
		{
			PageCondition: types.PageCondition{
				SheetName:       sheetName,
				IncludeFormulas: includeFormulas,
			},
			Actions: []types.Action{},
		},
	}

	rulesExecutor := MakeRulesExecutor(file, rules)
	transformErr := rulesExecutor.Execute()
	if transformErr != nil {
		return nil, fmt.Errorf("%s", transformErr.Message)
	}

	return file, nil
}

func runFormulaTestCase(t *testing.T, tc formulaTestCase) {
	t.Helper()

	file, err := executeFormulaRule(tc.inputFile, tc.includeFormulas, tc.sheetName)
	if err != nil {
		t.Fatalf("failed to execute formula rule: %v", err)
	}

	formulas, err := file.GetCellFormula(tc.sheetName, tc.cellReference)
	if err != nil {
		t.Fatalf("failed to get cell formula: %v", err)
	}
	assert.Equal(t, formulas, tc.expectedFormula)
}

func TestCheckFormulas(t *testing.T) {
	testCases := []formulaTestCase{
		{
			name:            "check formulas are cleared",
			inputFile:       "../assets/goldenFiles/testCheckFormulas.xlsx",
			includeFormulas: false,
			expectedFormula: "",
			sheetName:       "Sheet1",
			cellReference:   "E2",
		},
		{
			name:            "check formulas exist",
			inputFile:       "../assets/goldenFiles/testCheckFormulas.xlsx",
			includeFormulas: true,
			expectedFormula: "A2+B2",
			sheetName:       "Sheet1",
			cellReference:   "E2",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			runFormulaTestCase(t, tc)
		})
	}
}
