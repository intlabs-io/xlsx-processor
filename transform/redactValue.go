package transform

import (
	"xlsx-processor/pkg/cell"

	"github.com/xuri/excelize/v2"
)

func (a *ActionExecutor) RedactValue() (err error) {
	file := a.File
	sheetName := a.SheetName
	nonEmptyValueRedact := a.NonEmptyValueRedact
	valueToRedact := a.Action.Value

	cols, err := file.GetCols(sheetName)
	if err != nil {
		return err
	}

	for colIndex, col := range cols {
		// Iterate over each cell in the row
		for rowIndex, cellValue := range col {
			// Access individual cell
			cellName, err := excelize.CoordinatesToCellName(colIndex+1, rowIndex+1)
			if err != nil {
				return err
			}
			// Check if the cell value is the same as the valueToRedact
			if cellValue == valueToRedact {
				// Redact the cell
				err = cell.SetValue(file, sheetName, cellName, nonEmptyValueRedact, "**redacted**")
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}