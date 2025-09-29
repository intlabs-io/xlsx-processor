package sheet

import (
	"github.com/xuri/excelize/v2"
)

// Clearing the formulas of all the cells in the given sheet
func ClearFormulas(f *excelize.File, sheetName string) (err error) {
	cols, err := f.GetCols(sheetName)
	if err != nil {
		return err
	}

	for colIndex, col := range cols {
		// Iterate over each cell in the row
		for rowIndex, _ := range col {
			// Getting the cell column and row pair, eg: A1
			cellName, err := excelize.CoordinatesToCellName(colIndex+1, rowIndex+1)
			if err != nil {
				return err
			}
			// Clearing the formula
			err = f.SetCellFormula(sheetName, cellName, "")
			if err != nil {
				return err
			}
		}
	}

	return nil
}