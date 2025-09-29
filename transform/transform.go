package transform

import (
	"fmt"
	"slices"
	"xlsx-processor/pkg/sheet"
	"xlsx-processor/pkg/types"

	"github.com/xuri/excelize/v2"
)

// Constants for the operations and action types
const (
	REDACT    string = "REDACT"
	EXCLUDE   string = "EXCLUDE"
	VALUE     string = "VALUE"
	RANGE     string = "RANGE"
	TEXT_COLOR string = "TEXT_COLOR"
	BG_COLOR   string = "BG_COLOR"
	COLUMN    string = "COLUMN"
	ROW       string = "ROW"
)

type RulesExecutor struct {
	File *excelize.File
	rules *[]types.Rule
}

func MakeRulesExecutor(file *excelize.File, rules []types.Rule) *RulesExecutor {
	return &RulesExecutor{
		File: file,
		rules: &rules,
	}
}

func (r *RulesExecutor) Execute() *types.TransformError {
	file := r.File
	rules := r.rules

	for ruleIndex, rule := range *rules {
		var sheetName string = rule.PageCondition.SheetName
		var nonEmptyValueRedact bool = rule.PageCondition.NonEmptyValueRedact

		doesSheetNameExist := slices.Contains(file.GetSheetList(), sheetName)
		if !doesSheetNameExist {
			fmt.Println("Sheet name does not exist, skipping rule", sheetName)
			return nil
		}

		// Removing formulas if the rule is set to not include them
		if !rule.PageCondition.IncludeFormulas {
			err := sheet.ClearFormulas(file, sheetName)
			if err != nil {
				return &types.TransformError{
					Message: err.Error(),
					RuleIndex: &ruleIndex,
					Key: "includeFormulas",
				}
			}
		}

		for actionIndex, action := range rule.Actions {
			// Initialize the operations
			actionExecutor := MakeActionExecutor(file, sheetName, nonEmptyValueRedact, &action, actionIndex, ruleIndex)
			// Execute the action
			transformErr := actionExecutor.Execute()
			if transformErr != nil {
				return transformErr
			}
		}
	}
	return nil
}
