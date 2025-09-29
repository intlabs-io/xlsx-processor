package cell

import (
	"fmt"
	"strings"

	"github.com/xuri/excelize/v2"
)

// Getting the start and end cells of a sheet
func GetRange(f *excelize.File, sheetName string) (startCell string, startRow int, startCol string, endCell string, endRow int, endCol string, err error) {
	sheetDimension, err := f.GetSheetDimension(sheetName)
	if err != nil {
		return "", 0, "", "", 0, "", err
	}

	// Check if GetSheetDimension returned only "A1" (indicating no proper range detected)
	if sheetDimension == "A1" {
		fmt.Println("GetSheetDimension returned A1, attempting to find actual data range...")
		actualStartCell, actualStartRow, actualStartCol, actualEndCell, actualEndRow, actualEndCol, actualErr := findActualDataRange(f, sheetName)
		if actualErr == nil && actualEndCell != "A1" {
			fmt.Printf("Found actual data range: %s:%s\n", actualStartCell, actualEndCell)
			return actualStartCell, actualStartRow, actualStartCol, actualEndCell, actualEndRow, actualEndCol, nil
		}
		fmt.Printf("Could not find larger data range, using original: %s\n", sheetDimension)
	}

	// Converting cells into numbers
	split := strings.Split(sheetDimension, ":")

	// Check if we have exactly 2 parts after splitting
	if len(split) != 2 {
		return "", 0, "", "", 0, "", fmt.Errorf("invalid sheet dimension format: %s", sheetDimension)
	}

	startCell, endCell = split[0], split[1]

	startCol, startRow, err = SplitReference(startCell)
	if err != nil {
		return "", 0, "", "", 0, "", err
	}

	endCol, endRow, err = SplitReference(endCell)
	if err != nil {
		return "", 0, "", "", 0, "", err
	}

	return startCell, startRow, startCol, endCell, endRow, endCol, nil
}

// findActualDataRange scans the sheet to find the actual range of data
// This is a fallback when GetSheetDimension returns only "A1"
func findActualDataRange(f *excelize.File, sheetName string) (startCell string, startRow int, startCol string, endCell string, endRow int, endCol string, err error) {
	// Get all rows from the sheet
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return "", 0, "", "", 0, "", err
	}

	if len(rows) == 0 {
		return "A1", 1, "A", "A1", 1, "A", nil
	}

	maxRow := 0
	maxCol := 0
	minRow := -1
	minCol := -1

	// Scan through all rows to find the actual data boundaries
	for rowIdx, row := range rows {
		rowNum := rowIdx + 1

		// Check if row has any non-empty cells
		hasData := false
		for colIdx, cellValue := range row {
			if strings.TrimSpace(cellValue) != "" {
				hasData = true
				colNum := colIdx + 1

				if minRow == -1 || rowNum < minRow {
					minRow = rowNum
				}
				if minCol == -1 || colNum < minCol {
					minCol = colNum
				}
				if rowNum > maxRow {
					maxRow = rowNum
				}
				if colNum > maxCol {
					maxCol = colNum
				}
			}
		}

		// If this row has data, update maxRow
		if hasData && rowNum > maxRow {
			maxRow = rowNum
		}
	}

	// If no data found, return A1
	if minRow == -1 {
		return "A1", 1, "A", "A1", 1, "A", nil
	}

	// Convert column numbers to letters
	startColLetter := numberToColumn(minCol)
	endColLetter := numberToColumn(maxCol)

	startCell = fmt.Sprintf("%s%d", startColLetter, minRow)
	endCell = fmt.Sprintf("%s%d", endColLetter, maxRow)

	return startCell, minRow, startColLetter, endCell, maxRow, endColLetter, nil
}

// numberToColumn converts a column number to Excel column letters (1=A, 26=Z, 27=AA, etc.)
func numberToColumn(colNum int) string {
	result := ""
	for colNum > 0 {
		colNum-- // Adjust for 0-based calculation
		result = string(rune('A'+colNum%26)) + result
		colNum /= 26
	}
	return result
}
