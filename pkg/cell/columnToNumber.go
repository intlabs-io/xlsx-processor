package cell

import (
	"strings"
)

// Converting a column letter to the number of the column
func ColumnToNumber(column string) int {
	column = strings.ToUpper(column)
	// Validate input - ensure it contains only A-Z characters
	for _, char := range column {
		if char < 'A' || char > 'Z' {
			return 0
		}
	}
	result := 0
	for _, char := range column {
		result = result*26 + int(char-'A') + 1
	}
	return result
}
