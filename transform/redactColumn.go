package transform

import (
	"fmt"

	"github.com/xuri/excelize/v2"

	"xlsx-processor/pkg/cell"
)

func (a *ActionExecutor) RedactColumn() (err error) {
	file := a.File
	sheetName := a.SheetName
	nonEmptyValueRedact := a.NonEmptyValueRedact
	col := a.Action.Value
	/*
		User input validation
	*/
	if len(col) == 0 {
		return fmt.Errorf("'%s' is invalid", col)
	}

	// Get the dimensions of the sheet
	_, _, startCol, _, _, endCol, err := cell.GetRange(file, sheetName)
	if err != nil {
		return err
	}
	// Convert the column letters to numbers
	startColAsNum := cell.ColumnToNumber(startCol)
	endColAsNum := cell.ColumnToNumber(endCol)
	colAsNum := cell.ColumnToNumber(col)
	// Check to see if column is out of range
	if colAsNum < startColAsNum || colAsNum > endColAsNum {
		return fmt.Errorf("'%s' is out of range", col)
	}

	/*
		Redact the column
	*/

	// Find the column index that corresponds to the column letter
	targetColIndex := colAsNum - 1 // Convert to 0-based index
	
	// Get all rows to iterate through them
	rows, err := file.GetRows(sheetName)
	if err != nil {
		return err
	}
	
	// Iterate through each row and redact the cell in the target column
	for rowIndex, row := range rows {
		// Skip if the row doesn't have enough cells
		if len(row) <= targetColIndex {
			continue
		}
		
		// Get the cell coordinates
		cellName, err := excelize.CoordinatesToCellName(colAsNum, rowIndex+1)
		if err != nil {
			return err
		}
		// Redact the cell
		err = cell.SetValue(file, sheetName, cellName, nonEmptyValueRedact, "**redacted**")
		if err != nil {
			return err
		}
	}

	return nil
}