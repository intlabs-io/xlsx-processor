package cell

import (
	"testing"

	"github.com/xuri/excelize/v2"
)

func TestGetTextColor(t *testing.T) {
	t.Run("get text colors from test file", func(t *testing.T) {
		file, err := excelize.OpenFile("../../assets/goldenFiles/test.xlsx")
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}
		defer file.Close()

		sheetName := "Forecasting"

		// Test cases with expected results based on the test file
		testCases := []struct {
			cellRef     string
			description string
		}{
			{"A1", "First cell"},
			{"B1", "Second cell"},
			{"C1", "Third cell"},
			{"A2", "Second row first cell"},
			{"B2", "Second row second cell"},
			{"A3", "Third row first cell"},
			{"B3", "Third row second cell"},
			{"C3", "Third row third cell"},
			{"D3", "Third row fourth cell"},
		}

		t.Logf("Testing text colors in sheet: %s", sheetName)
		t.Logf("%-10s %-25s %s", "Cell", "Description", "Text Color")
		t.Logf("%-10s %-25s %s", "----", "-----------", "----------")

		for _, tc := range testCases {
			textColor, err := GetTextColor(file, sheetName, tc.cellRef)
			if err != nil {
				t.Errorf("Expected no error for cell %s, got: %v", tc.cellRef, err)
				continue
			}

			t.Logf("%-10s %-25s %s", tc.cellRef, tc.description, textColor)

			// Verify we get a valid hex color (6 characters) or default
			if textColor != "" && len(textColor) != 6 {
				t.Errorf("Expected valid 6-character hex color for cell %s, got: %s (length: %d)", tc.cellRef, textColor, len(textColor))
			}
		}
	})

	t.Run("test different sheets", func(t *testing.T) {
		file, err := excelize.OpenFile("../../assets/goldenFiles/test.xlsx")
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}
		defer file.Close()

		sheets := file.GetSheetList()
		t.Logf("Available sheets: %v", sheets)

		for _, sheetName := range sheets {
			t.Logf("\nTesting sheet: %s", sheetName)

			// Test a few cells in each sheet
			testCells := []string{"A1", "B1", "A2", "B2"}
			for _, cellRef := range testCells {
				textColor, err := GetTextColor(file, sheetName, cellRef)
				if err != nil {
					t.Logf("  %s: error - %v", cellRef, err)
					continue
				}
				t.Logf("  %s: %s", cellRef, textColor)
			}
		}
	})

	t.Run("compare with known colors from attributes test", func(t *testing.T) {
		file, err := excelize.OpenFile("../../assets/goldenFiles/test.xlsx")
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}
		defer file.Close()

		// These are the text colors we expect to find based on the attributes test
		expectedTextColors := []string{"FF0000", "0070C0"}

		t.Logf("Expected text colors from test: %v", expectedTextColors)
		t.Logf("Scanning all sheets to find these colors...")

		foundColors := make(map[string][]string) // color -> list of cells where found

		sheets := file.GetSheetList()
		for _, sheetName := range sheets {
			// Get sheet dimensions to scan all cells
			rows, err := file.GetRows(sheetName)
			if err != nil {
				continue
			}

			for rowIndex, row := range rows {
				for colIndex := range row {
					cellRef, err := excelize.CoordinatesToCellName(colIndex+1, rowIndex+1)
					if err != nil {
						continue
					}

					textColor, err := GetTextColor(file, sheetName, cellRef)
					if err != nil || textColor == "000000" || textColor == "FFFFFF" {
						continue // Skip errors and default colors
					}

					// Check if this is one of our expected colors
					for _, expectedColor := range expectedTextColors {
						if textColor == expectedColor {
							foundColors[textColor] = append(foundColors[textColor], sheetName+":"+cellRef)
							break
						}
					}
				}
			}
		}

		t.Logf("\nFound expected text colors:")
		for color, locations := range foundColors {
			t.Logf("  %s: found in %d locations: %v", color, len(locations), locations[:min(5, len(locations))])
		}

		t.Logf("\nSummary: Found %d out of %d expected text colors", len(foundColors), len(expectedTextColors))
	})

	t.Run("test invalid inputs", func(t *testing.T) {
		file, err := excelize.OpenFile("../../assets/goldenFiles/test.xlsx")
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}
		defer file.Close()

		// Test invalid sheet name
		textColor, err := GetTextColor(file, "NonExistentSheet", "A1")
		if err == nil {
			t.Logf("Invalid sheet returned color: %s (this might be expected behavior)", textColor)
		} else {
			t.Logf("Invalid sheet returned error: %v", err)
		}

		// Test invalid cell reference
		textColor, err = GetTextColor(file, "Forecasting", "InvalidCell")
		if err == nil {
			t.Logf("Invalid cell returned color: %s (this might be expected behavior)", textColor)
		} else {
			t.Logf("Invalid cell returned error: %v", err)
		}
	})
}
