package transform

import (
	"fmt"
	"strings"
	"xlsx-processor/pkg/cell"

	"github.com/xuri/excelize/v2"
)

func (a *ActionExecutor) RedactRange() (err error) {
	file := a.File
	sheetName := a.SheetName
	nonEmptyValueRedact := a.NonEmptyValueRedact
	rangeString := a.Action.Value
	split := strings.Split(rangeString, ":")
	// Checking if value is formatted properly
	if len(split) <= 1 || len(split[0]) <= 1 || len(split[1]) <= 1 {
		return fmt.Errorf("'%s' is invalid", rangeString)
	}

	startCell, endCell := split[0], split[1]

	// Checking if start and end cells are out of range
	isStartCellOutOfRange, err := cell.IsOutOfRange(file, sheetName, startCell)
	if err != nil {
		return err
	}
	isEndCellOutOfRange, err := cell.IsOutOfRange(file, sheetName, endCell)
	if err != nil {
		return err
	}
	if isStartCellOutOfRange {
		return fmt.Errorf("'%s' starting cell is out of range", rangeString)
	}

	if isEndCellOutOfRange {
		return fmt.Errorf("'%s' ending cell is out of range", rangeString)
	}

	// Splitting cell reference into column and row
	startCol, startRowNum, err := cell.SplitReference(startCell)
	if err != nil {
		return err
	}
	startColNum := cell.ColumnToNumber(startCol)

	resetStartRowNum := startRowNum

	// Splitting cell reference into column and row
	endCol, endRowNum, err := cell.SplitReference(endCell)
	if err != nil {
		return err
	}
	endColNum := cell.ColumnToNumber(endCol)

	// Iterating over the range and redacting the values
	for startColNum <= endColNum {
		startRowNum = resetStartRowNum
		for startRowNum <= endRowNum {
			// Getting the cell column and row pair, eg: A1
			cellName, _ := excelize.CoordinatesToCellName(startColNum, startRowNum)
			// Redacting the cell
			err := cell.SetValue(file, sheetName, cellName, nonEmptyValueRedact, "**redacted**")
			if err != nil {
				return err
			}
			startRowNum++
		}
		startColNum++
	}

	return nil
}