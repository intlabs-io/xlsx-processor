package transform

import (
	"fmt"
	"xlsx-processor/pkg/cell"

	"github.com/xuri/excelize/v2"
)

func (a *ActionExecutor) RedactBgColor() (err error) {
	file := a.File
	sheetName := a.SheetName
	nonEmptyValueRedact := a.NonEmptyValueRedact
	colorHex := a.Action.Value
	foundBgColor := false

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
			// Get the background color of the cell
			bgColor, err := cell.GetBgColor(file, sheetName, cellName)
			if err != nil {
				return err
			}
			// Check if the background color is the same as the colorHex
			if bgColor == colorHex {
				foundBgColor = true
				// Redact the cell
				err = cell.SetValue(file, sheetName, cellName, nonEmptyValueRedact, "**redacted**")
				if err != nil {
					return err
				}
			}
		}
	}

	// If the background color was not found then return an error
	if !foundBgColor {
		return fmt.Errorf("'%s' was not found in the sheet", colorHex)
	}

	return nil
}
