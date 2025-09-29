package cell

import (
	"github.com/xuri/excelize/v2"
)

// Setting the value of a cell
func SetValue(f *excelize.File, sheetName string, cellReference string, nonEmptyValueRedact bool, setValue string) (err error) {
	cellValue, err := f.GetCellValue(sheetName, cellReference)
	if err != nil {
		return err
	}

	if nonEmptyValueRedact && cellValue != "" {
		err = f.SetCellValue(sheetName, cellReference, setValue)
		if err != nil {
			return err
		}
	} else if !nonEmptyValueRedact {
		err = f.SetCellValue(sheetName, cellReference, "")
		if err != nil {
			return err
		}
	}

	return nil
}
