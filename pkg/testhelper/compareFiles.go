package testhelper

import (
	"testing"
	"xlsx-processor/pkg/sheet"

	"github.com/go-playground/assert/v2"
	"github.com/xuri/excelize/v2"
)

func CompareSheet(t *testing.T, expected, actual *excelize.File, sheetName *string) {
	expectedSheet, err := sheet.ParseSheetToCsv(expected, sheetName)
	if err != nil {
		t.Fatalf("failed to convert expected sheet to csv: %v", err)
	}

	actualSheet, err := sheet.ParseSheetToCsv(actual, sheetName)
	if err != nil {
		t.Fatalf("failed to convert actual sheet to csv: %v", err)
	}

	assert.Equal(t, expectedSheet, actualSheet)
}
