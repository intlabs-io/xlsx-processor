package storage

import (
	"encoding/json"
	"fmt"
	"xlsx-processor/pkg/file"
	"xlsx-processor/pkg/types"

	"github.com/xuri/excelize/v2"
)

func downloadProxy(storageType string, input types.Input) ([]byte, error) {
	switch storageType {
	case "S3", "FILE":
		return DownloadFromS3Input(input)
	case "ONEDRIVE":
		return DownloadFromOneDrive(input)
	case "GOOGLEDRIVE":
		return DownloadFromGoogleDrive(input)
	default:
		return nil, fmt.Errorf("storage type not found")
	}
}

func GetFileFromJson(storageType string, input types.Input) (*excelize.File, error) {
	fileBytes, err := downloadProxy(storageType, input)
	if err != nil {
		return nil, fmt.Errorf("failed to download file: %w", err)
	}

	var sheetContent types.Sheet
	err = json.Unmarshal(fileBytes, &sheetContent)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal file: %w", err)
	}

	// Convert the complete JSON content to an excelize file
	f, err := file.ConvertSheetToExcelizeFile(sheetContent)
	if err != nil {
		return nil, fmt.Errorf("failed to convert JSON content to excelize file: %w", err)
	}

	return f, nil
}

func GetFile(storageType string, input types.Input) (*excelize.File, error) {
	fileBytes, err := downloadProxy(storageType, input)
	if err != nil {
		return nil, fmt.Errorf("failed to download file: %w", err)
	}

	f, err := file.InitFileFromBytes(fileBytes)
	if err != nil {
		return nil, err
	}

	return f, nil
}
