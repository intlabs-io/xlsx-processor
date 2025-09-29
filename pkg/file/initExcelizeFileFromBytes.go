package file

import (
	"bytes"

	"github.com/xuri/excelize/v2"
)

// Initializing an excelize file from bytes
func InitFileFromBytes(fileContents []byte) (*excelize.File, error) {
	reader := bytes.NewReader(fileContents)
	f, err := excelize.OpenReader(reader)
	if err != nil {
		return nil, err
	}

	return f, nil
}
