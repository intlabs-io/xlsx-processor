package sheet

import (
	"encoding/hex"
	"fmt"

	"xlsx-processor/pkg/cell"
	"xlsx-processor/pkg/types"

	"github.com/xuri/excelize/v2"
)

// getSheetTabColor extracts the hex color code from sheet properties
func getSheetTabColor(sheetProps excelize.SheetPropsOptions) (string, error) {
	if sheetProps.TabColorRGB == nil {
		return "", nil
	}

	b, err := hex.DecodeString(*sheetProps.TabColorRGB)
	if err != nil {
		return "", err
	}

	if len(b) < 3 {
		return "", nil
	}

	// Skip the optional alpha byte if len==4
	start := 0
	if len(b) == 4 {
		start = 1
	}
	return fmt.Sprintf("%02X%02X%02X", b[start], b[start+1], b[start+2]), nil
}

// collectUniqueColor adds a color to the slice if it's not already present
func collectUniqueColor(colors []string, newColor string) []string {
	for _, color := range colors {
		if color == newColor {
			return colors
		}
	}
	return append(colors, newColor)
}

// processCellStyle extracts style information from a cell and updates color collections
func processCellStyle(f *excelize.File, sheetName, cellName string, style *excelize.Style, textColors, bgColors []string) ([]string, []string) {
	if style == nil {
		return textColors, bgColors
	}

	// Process text color using the improved GetTextColor function
	textColor, err := cell.GetTextColor(f, sheetName, cellName)
	if err == nil && len(textColor) > 0 && textColor != "000000" && textColor != "FFFFFF" {
		textColors = collectUniqueColor(textColors, textColor)
	}

	// Process background color using the improved GetBgColor function
	bgColor, err := cell.GetBgColor(f, sheetName, cellName)
	if err == nil && len(bgColor) > 0 && bgColor != "FFFFFF" {
		bgColors = collectUniqueColor(bgColors, bgColor)
	}

	return textColors, bgColors
}

// collectSheetMinimals collects basic information about all sheets
func collectSheetMinimals(f *excelize.File, sheetNames []string) ([]types.SheetMinimal, error) {
	sheetMinimals := make([]types.SheetMinimal, len(sheetNames))
	for i, sheet := range sheetNames {
		// Get sheet tab color
		sheetProps, err := f.GetSheetProps(sheet)
		if err != nil {
			return nil, fmt.Errorf("failed to get sheet properties: %w", err)
		}

		hexCode, err := getSheetTabColor(sheetProps)
		if err != nil {
			return nil, fmt.Errorf("failed to process sheet tab color: %w", err)
		}

		sheetMinimals[i] = types.SheetMinimal{
			SheetName:     sheet,
			SheetTabColor: hexCode,
		}
	}

	return sheetMinimals, nil
}

// Converting an xlsx sheet to a csv object
func ParseSheetToCsv(f *excelize.File, sheetName *string) (*types.Sheet, error) {
	sheetNames := f.GetSheetList()

	// Find the index of the requested sheet
	sheetIndex := 0
	if sheetName != nil {
		for i, s := range sheetNames {
			if s == *sheetName {
				sheetIndex = i
				break
			}
		}
	}

	currentSheet := sheetNames[sheetIndex]

	// Get rows and process styles
	rows, err := f.GetRows(currentSheet)
	if err != nil {
		return nil, fmt.Errorf("failed to get rows: %w", err)
	}

	styledRows := make([][]types.StyledCell, len(rows))
	for rowIndex, row := range rows {
		styledRow := make([]types.StyledCell, len(row))

		for colIndex, cellValue := range row {
			cellName, err := excelize.CoordinatesToCellName(colIndex+1, rowIndex+1)
			if err != nil {
				return nil, fmt.Errorf("failed to convert coordinates: %w", err)
			}

			styleIndex, err := f.GetCellStyle(currentSheet, cellName)
			if err != nil {
				return nil, fmt.Errorf("failed to get cell style: %w", err)
			}

			style, err := f.GetStyle(styleIndex)
			if err != nil {
				style = nil
			}

			styledCell := types.StyledCell{
				Value: cellValue,
				Style: style,
			}

			styledRow[colIndex] = styledCell
		}

		styledRows[rowIndex] = styledRow
	}

	sheet := &types.Sheet{
		SheetName: currentSheet,
		Cells:     styledRows,
	}

	return sheet, nil
}

func ParseSheetsToCsvAndAttributes(f *excelize.File) ([]*types.Sheet, *types.Attributes, error) {
	sheetNames := f.GetSheetList()

	sheets := make([]*types.Sheet, len(sheetNames))

	// Unique textColors and bgColors for all the sheets
	var textColors, bgColors []string
	// Parse all the sheets and collect the attributes for all the sheets
	for i, sheetName := range sheetNames {
		// Get rows and process styles
		rows, err := f.GetRows(sheetName)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to get rows: %w", err)
		}

		styledRows := make([][]types.StyledCell, len(rows))

		for rowIndex, row := range rows {
			styledRow := make([]types.StyledCell, len(row))

			for colIndex, cellValue := range row {
				cellName, err := excelize.CoordinatesToCellName(colIndex+1, rowIndex+1)
				if err != nil {
					return nil, nil, fmt.Errorf("failed to convert coordinates: %w", err)
				}

				styleIndex, err := f.GetCellStyle(sheetName, cellName)
				if err != nil {
					return nil, nil, fmt.Errorf("failed to get cell style: %w", err)
				}

				style, err := f.GetStyle(styleIndex)
				if err != nil {
					style = nil
				}

				styledCell := types.StyledCell{
					Value: cellValue,
					Style: style,
				}

				textColors, bgColors = processCellStyle(f, sheetName, cellName, style, textColors, bgColors)
				styledRow[colIndex] = styledCell
			}

			styledRows[rowIndex] = styledRow
		}

		sheet := &types.Sheet{
			SheetName: sheetName,
			Cells:     styledRows,
		}

		sheets[i] = sheet
	}

	sheetMinimals, err := collectSheetMinimals(f, sheetNames)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to collect sheet minimals: %w", err)
	}

	attributes := &types.Attributes{
		SheetMinimals: sheetMinimals,
		TextColors:    textColors,
		BgColors:      bgColors,
	}

	return sheets, attributes, nil
}
