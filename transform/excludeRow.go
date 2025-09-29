package transform

import (
	"fmt"
	"strconv"
	"xlsx-processor/pkg/cell"
)

func (a *ActionExecutor) ExcludeRow() (err error) {
	file := a.File
	sheetName := a.SheetName
	row := a.Action.Value

	// Get the dimensions of the sheet
	_, startRowNum, _, _, endRowNum, _, err := cell.GetRange(file, sheetName)
	if err != nil {
		return err
	}
	rowNum, err := strconv.Atoi(row)
	if err != nil {
		return err
	}
	// Check to see if row is out of range
	if rowNum < startRowNum || rowNum > endRowNum {
		return fmt.Errorf("'%d' is out of range", rowNum)
	}
	// Remove the row
	err = file.RemoveRow(sheetName, rowNum)
	if err != nil {
		return err
	}

	return nil
}