package transform

import (
	"fmt"

	"xlsx-processor/pkg/cell"

	"github.com/xuri/excelize/v2"
)

func (a *ActionExecutor) RedactTextColor() (err error) {
	file := a.File
	sheetName := a.SheetName
	nonEmptyValueRedact := a.NonEmptyValueRedact
	colorHex := a.Action.Value
	foundTextColor := false

	cols, err := file.GetCols(sheetName)
	if err != nil {
		return err
	}
	
	for colIndex, col := range cols {
		// Iterate over each cell in the row
		for rowIndex, _ := range col {
			// Access individual cell
			cellName, err := excelize.CoordinatesToCellName(colIndex+1, rowIndex+1)
			if err != nil {
				return err
			}
			// Get the text color of the cell
			textColor, err := cell.GetTextColor(file, sheetName, cellName)
			if err != nil {
				return err
			}
			// Check if the text color is the same as the colorHex
			if textColor == colorHex {
				foundTextColor = true
				// Redact the cell
				err = cell.SetValue(file, sheetName, cellName, nonEmptyValueRedact, "**redacted**")
				if err != nil {
					return err
				}
			}
		}
	}

	// If the text color was not found then return an error
	if !foundTextColor {
		return fmt.Errorf("'%s' was not found in sheet '%s'", colorHex, sheetName)
	}

	return nil
}