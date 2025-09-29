package file

import (
	"testing"
	"xlsx-processor/pkg/types"

	"github.com/go-playground/assert/v2"
	"github.com/xuri/excelize/v2"
)

func TestConvertSheetToExcelizeFile(t *testing.T) {
	t.Run("convert basic sheet with values and styles", func(t *testing.T) {
		// Create a sample types.Sheet
		sheet := types.Sheet{
			SheetName: "TestSheet",
			Cells: [][]types.StyledCell{
				{
					{Value: "Header1", Style: &excelize.Style{Font: &excelize.Font{Bold: true}}},
					{Value: "Header2", Style: &excelize.Style{Font: &excelize.Font{Bold: true}}},
				},
				{
					{Value: "Row1Col1", Style: nil},
					{Value: "Row1Col2", Style: nil},
				},
				{
					{Value: "Row2Col1", Style: nil},
					{Value: "Row2Col2", Style: &excelize.Style{Font: &excelize.Font{Color: "FF0000"}}},
				},
			},
		}

		// Convert to excelize file
		f, err := ConvertSheetToExcelizeFile(sheet)
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		// Verify sheet name
		sheetNames := f.GetSheetList()
		assert.Equal(t, 1, len(sheetNames))
		assert.Equal(t, "TestSheet", sheetNames[0])

		// Verify cell values
		val, err := f.GetCellValue("TestSheet", "A1")
		if err != nil {
			t.Fatalf("Failed to get cell value: %v", err)
		}
		assert.Equal(t, "Header1", val)

		val, err = f.GetCellValue("TestSheet", "B1")
		if err != nil {
			t.Fatalf("Failed to get cell value: %v", err)
		}
		assert.Equal(t, "Header2", val)

		val, err = f.GetCellValue("TestSheet", "A2")
		if err != nil {
			t.Fatalf("Failed to get cell value: %v", err)
		}
		assert.Equal(t, "Row1Col1", val)

		val, err = f.GetCellValue("TestSheet", "B3")
		if err != nil {
			t.Fatalf("Failed to get cell value: %v", err)
		}
		assert.Equal(t, "Row2Col2", val)
	})

	t.Run("convert empty sheet", func(t *testing.T) {
		// Create an empty sheet
		sheet := types.Sheet{
			SheetName: "EmptySheet",
			Cells:     [][]types.StyledCell{},
		}

		// Convert to excelize file
		f, err := ConvertSheetToExcelizeFile(sheet)
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		// Verify sheet name
		sheetNames := f.GetSheetList()
		assert.Equal(t, 1, len(sheetNames))
		assert.Equal(t, "EmptySheet", sheetNames[0])
	})

	t.Run("convert sheet with jagged rows", func(t *testing.T) {
		// Create a sheet with rows of different lengths
		sheet := types.Sheet{
			SheetName: "JaggedSheet",
			Cells: [][]types.StyledCell{
				{
					{Value: "A1", Style: nil},
					{Value: "B1", Style: nil},
					{Value: "C1", Style: nil},
				},
				{
					{Value: "A2", Style: nil},
				},
				{
					{Value: "A3", Style: nil},
					{Value: "B3", Style: nil},
				},
			},
		}

		// Convert to excelize file
		f, err := ConvertSheetToExcelizeFile(sheet)
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		// Verify values
		val, err := f.GetCellValue("JaggedSheet", "C1")
		if err != nil {
			t.Fatalf("Failed to get cell value: %v", err)
		}
		assert.Equal(t, "C1", val)

		val, err = f.GetCellValue("JaggedSheet", "A2")
		if err != nil {
			t.Fatalf("Failed to get cell value: %v", err)
		}
		assert.Equal(t, "A2", val)

		// C2 should be empty
		val, err = f.GetCellValue("JaggedSheet", "C2")
		if err != nil {
			t.Fatalf("Failed to get cell value: %v", err)
		}
		assert.Equal(t, "", val)
	})
}
