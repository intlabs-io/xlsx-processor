package file

import (
	"fmt"
	"xlsx-processor/pkg/types"

	"github.com/xuri/excelize/v2"
)

// ConvertSheetToExcelizeFile converts a types.Sheet into an excelize.File
func ConvertSheetToExcelizeFile(sheet types.Sheet) (*excelize.File, error) {
	// Create a new excelize file
	f := excelize.NewFile()

	// Set the sheet name (rename the default sheet)
	defaultSheetName := f.GetSheetName(0)
	err := f.SetSheetName(defaultSheetName, sheet.SheetName)
	if err != nil {
		return nil, fmt.Errorf("failed to set sheet name: %w", err)
	}

	// If there are no cells, return the empty file
	if len(sheet.Cells) == 0 {
		return f, nil
	}

	// Iterate through all cells and populate the excelize file
	for rowIndex, row := range sheet.Cells {
		for colIndex, cell := range row {
			// Convert 0-based indices to 1-based for excelize
			excelRow := rowIndex + 1
			excelCol := colIndex + 1

			// Get the cell coordinate (e.g., "A1", "B2")
			cellCoord, err := excelize.CoordinatesToCellName(excelCol, excelRow)
			if err != nil {
				return nil, fmt.Errorf("failed to convert coordinates for cell [%d,%d]: %w", rowIndex, colIndex, err)
			}

			// Set the cell value
			err = f.SetCellValue(sheet.SheetName, cellCoord, cell.Value)
			if err != nil {
				return nil, fmt.Errorf("failed to set cell value at %s: %w", cellCoord, err)
			}

			// Set the cell style if it exists
			if cell.Style != nil {
				// Create a new style in the excelize file based on the existing style
				styleID, err := f.NewStyle(cell.Style)
				if err != nil {
					return nil, fmt.Errorf("failed to create style for cell %s: %w", cellCoord, err)
				}

				// Apply the style to the cell
				err = f.SetCellStyle(sheet.SheetName, cellCoord, cellCoord, styleID)
				if err != nil {
					return nil, fmt.Errorf("failed to set cell style at %s: %w", cellCoord, err)
				}
			}
		}
	}

	return f, nil
}
