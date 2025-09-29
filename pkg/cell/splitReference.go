package cell

import (
	"fmt"
	"regexp"
	"strconv"
)

// Splitting a cell reference into column and row
func SplitReference(cellReference string) (column string, row int, err error) {
	// Use regular expression to extract letters (column) and numbers (row)
	r := regexp.MustCompile("([A-Za-z]+)([0-9]+)")
	matches := r.FindStringSubmatch(cellReference)

	if len(matches) != 3 {
		return "", 0, fmt.Errorf("invalid cell reference: %s", cellReference)
	}

	column = matches[1]

	row, err = strconv.Atoi(matches[2])
	if err != nil {
		return "", 0, fmt.Errorf("failed to convert row to integer: %v", err)
	}

	return column, row, nil
}
