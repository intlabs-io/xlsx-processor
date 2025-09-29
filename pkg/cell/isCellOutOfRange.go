package cell

import (
	"github.com/xuri/excelize/v2"
)

// Checking if a cell is out of range
func IsOutOfRange(f *excelize.File, sheetName string, cellReference string) (bool, error) {
	_, startRow, startCol, _, endRow, endCol, err := GetRange(f, sheetName)
	if err != nil {
		return true, err
	}
	startColAsNum := ColumnToNumber(startCol)
	endColAsNum := ColumnToNumber(endCol)

	cellCol, cellRow, err := SplitReference(cellReference)
	if err != nil {
		return true, err
	}
	cellColAsNum := ColumnToNumber(cellCol)

	// Checking if row is out of range
	if cellRow < startRow || cellRow > endRow {
		return true, nil
	}

	// Checking if col is out of range
	if cellColAsNum < startColAsNum || cellColAsNum > endColAsNum {
		return true, nil
	}

	return false, nil
}
