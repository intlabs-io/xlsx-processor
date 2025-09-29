package cell

import (
	"github.com/xuri/excelize/v2"
)

// Getting the style of a cell
func GetStyle(f *excelize.File, sheetName, cellReference string) (styling *excelize.Style, err error) {
	styleIndex, err := f.GetCellStyle(sheetName, cellReference)
	if err != nil {
		return nil, err
	}
	style, err := f.GetStyle(styleIndex)
	if err != nil {
		return nil, err
	}

	return style, nil
}
