package sheet

import (
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/xuri/excelize/v2"
)

func TestParseSheetToCsv(t *testing.T) {
	t.Run("parse specific sheet", func(t *testing.T) {
		file, err := excelize.OpenFile("../../assets/goldenFiles/test.xlsx")
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}
		defer file.Close()

		sheetName := "Forecasting"
		sheet, err := ParseSheetToCsv(file, &sheetName)

		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		// Verify we got the correct sheet
		if sheet.SheetName != sheetName {
			t.Fatalf("Expected sheet name %s, got: %s", sheetName, sheet.SheetName)
		}

		// Verify we have some content
		if len(sheet.Cells) == 0 {
			t.Fatal("Expected sheet to have some cells")
		}
	})

	t.Run("parse first sheet when sheetName is nil", func(t *testing.T) {
		file, err := excelize.OpenFile("../../assets/goldenFiles/test.xlsx")
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}
		defer file.Close()

		sheet, err := ParseSheetToCsv(file, nil)

		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		// Verify we got a sheet (should be the first one)
		if sheet.SheetName == "" {
			t.Fatal("Expected sheet to have a name")
		}

		// Verify we have some content
		if len(sheet.Cells) == 0 {
			t.Fatal("Expected sheet to have some cells")
		}
	})
}

func TestParseSheetsToCsvAndAttributes(t *testing.T) {
	t.Run("parse all sheets and collect attributes", func(t *testing.T) {
		file, err := excelize.OpenFile("../../assets/goldenFiles/test.xlsx")
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}
		defer file.Close()

		sheets, attributes, err := ParseSheetsToCsvAndAttributes(file)

		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		// Verify we got multiple sheets
		if len(sheets) == 0 {
			t.Fatal("Expected at least one sheet")
		}

		// Verify sheet names
		sheetMinimals := attributes.SheetMinimals
		expectedSheetNames := []string{"Forecasting", "CONTROL"}
		actualSheetNames := make([]string, len(sheetMinimals))
		for i, minimal := range sheetMinimals {
			actualSheetNames[i] = minimal.SheetName
		}
		EqualNotSorted(t, expectedSheetNames, actualSheetNames)

		// Verify we got the same number of sheets and minimals
		if len(sheets) != len(sheetMinimals) {
			t.Fatalf("Expected %d sheets to match %d sheet minimals", len(sheets), len(sheetMinimals))
		}

		// Verify colors are collected
		textColors := attributes.TextColors
		EqualNotSorted(t, []string{"FF0000", "0070C0"}, textColors)

		bgColors := attributes.BgColors
		EqualNotSorted(t, []string{"000000", "808080", "4472C4", "FFC000", "FF0000", "70AD47", "7030A0", "F2F2F2", "DAE3F3", "FFF2CC", "FBD0CF", "E2F0D9", "DDCBFF", "E7E6E6"}, bgColors)
	})
}

func EqualNotSorted(t *testing.T, expected, actual []string) {
	assert.Equal(t, len(expected), len(actual))
	for _, exp := range expected {
		assert.Equal(t, true, contains(actual, exp))
	}
}

// contains checks if a string slice contains a specific string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
