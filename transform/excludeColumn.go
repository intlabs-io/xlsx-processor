package transform

import (
	"fmt"

	"xlsx-processor/pkg/cell"
)

func (a *ActionExecutor) ExcludeColumn() (err error) {
	file := a.File
	sheetName := a.SheetName
	col := a.Action.Value

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
	// Remove the column
	err = file.RemoveCol(sheetName, col)
	if err != nil {
		return err
	}

	return nil
}