package transform

import (
	"fmt"
	"strconv"

	"github.com/xuri/excelize/v2"

	"xlsx-processor/pkg/cell"
)

func (a *ActionExecutor) RedactRow() (err error) {
	file := a.File
	sheetName := a.SheetName
	nonEmptyValueRedact := a.NonEmptyValueRedact
	row := a.Action.Value
	/*
		User input validation
	*/
	if len(row) == 0 {
		return fmt.Errorf("'%s' is invalid", row)
	}

	// Get the dimensions of the sheet
	_, startRowNum, _, _, endRowNum, _, err := cell.GetRange(file, sheetName)
	if err != nil {
		return err
	}
	// Convert the row string to an integer
	rowNum, err := strconv.Atoi(row)
	if err != nil {
		return fmt.Errorf("'%s' is not a valid row number", row)
	}
	// Check to see if row is out of range
	if rowNum < startRowNum || rowNum > endRowNum {
		return fmt.Errorf("'%s' is out of range", row)
	}

	/*
		Redact the row
	*/

	// Find the row index that corresponds to the row number
	targetRowIndex := rowNum - 1 // Convert to 0-based index

	// Get all columns to iterate through them
	cols, err := file.GetCols(sheetName)
	if err != nil {
		return err
	}

	// Iterate through each column and redact the cell in the target row
	for colIndex, col := range cols {
		// Skip if the column doesn't have enough cells
		if len(col) <= targetRowIndex {
			continue
		}

		// Get the cell coordinates
		cellName, err := excelize.CoordinatesToCellName(colIndex+1, rowNum)
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
