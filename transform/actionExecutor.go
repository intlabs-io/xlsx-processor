package transform

import (
	"xlsx-processor/pkg/types"

	"github.com/xuri/excelize/v2"
)

// Actions represents a collection of methods that operate on an excelize file
// with different action types: redact and exclude
type ActionExecutor struct {
	File                *excelize.File
	SheetName           string
	NonEmptyValueRedact bool
	Action              *types.Action
	ActionIndex         int
	RuleIndex           int
}

// MakeActionExecutor creates a new Actions instance
func MakeActionExecutor(f *excelize.File, sheetName string, nonEmptyValueRedact bool, action *types.Action, actionIndex int, ruleIndex int) *ActionExecutor {
	return &ActionExecutor{
		File:                f,
		SheetName:           sheetName,
		NonEmptyValueRedact: nonEmptyValueRedact,
		Action:              action,
		ActionIndex:         actionIndex,
		RuleIndex:           ruleIndex,
	}
}

// newTransformError creates a new TransformError with the given message and key
func (a *ActionExecutor) newTransformError(message, key string) *types.TransformError {
	return &types.TransformError{
		Message:     message,
		ActionIndex: &a.ActionIndex,
		RuleIndex:   &a.RuleIndex,
		Key:         key,
	}
}

// Execute executes the appropriate operation based on the action type
func (a *ActionExecutor) Execute() *types.TransformError {
	switch a.Action.ActionType {
	case REDACT:
		transformErr := a.ExecuteRedact()
		if transformErr != nil {
			return transformErr
		}
	case EXCLUDE:
		transformErr := a.ExecuteExclude()
		if transformErr != nil {
			return transformErr
		}
	default:
		return a.newTransformError("Invalid action type", "actionType")
	}

	return nil
}

// ExecuteRedact handles all redact operations
func (a *ActionExecutor) ExecuteRedact() *types.TransformError {
	// Skip empty values
	if a.Action.Value == "" {
		return nil
	}

	switch a.Action.Operation {
	case RANGE:
		if err := a.RedactRange(); err != nil {
			return a.newTransformError(err.Error(), "value")
		}
	case VALUE:
		if err := a.RedactValue(); err != nil {
			return a.newTransformError(err.Error(), "value")
		}
	case TEXT_COLOR:
		if err := a.RedactTextColor(); err != nil {
			return a.newTransformError(err.Error(), "value")
		}
	case BG_COLOR:
		if err := a.RedactBgColor(); err != nil {
			return a.newTransformError(err.Error(), "value")
		}
	case COLUMN:
		if err := a.RedactColumn(); err != nil {
			return a.newTransformError(err.Error(), "value")
		}
	case ROW:
		if err := a.RedactRow(); err != nil {
			return a.newTransformError(err.Error(), "value")
		}
	default:
		return a.newTransformError("Invalid operation", "operation")
	}

	return nil
}

// ExecuteExclude handles all exclude operations
func (a *ActionExecutor) ExecuteExclude() *types.TransformError {
	// Skip empty values
	if a.Action.Value == "" {
		return nil
	}

	switch a.Action.Operation {
	case ROW:
		if err := a.ExcludeRow(); err != nil {
			return a.newTransformError(err.Error(), "value")
		}
	case COLUMN:
		if err := a.ExcludeColumn(); err != nil {
			return a.newTransformError(err.Error(), "value")
		}
	default:
		return a.newTransformError("Invalid operation", "operation")
	}

	return nil
}
