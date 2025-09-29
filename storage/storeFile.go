package storage

import (
	"fmt"

	"github.com/xuri/excelize/v2"

	"xlsx-processor/pkg/types"
)

func uploadProxy(storageType string, output types.Output, fileContents []byte) error {
	switch storageType {
	case "S3", "FILE":
		return UploadToS3Output(output, fileContents)
	case "GOOGLEDRIVE":
		return fmt.Errorf("not implemented")
	default:
		return fmt.Errorf("storage type not found")
	}
}

func StoreFile(f *excelize.File, output types.Output, webhook *types.Webhook) error {
	// Saving the modified Excel file to a byte buffer
	fileContentsBuffer, err := f.WriteToBuffer()
	if err != nil {
		return err
	}

	// Uploading the modified Excel file to the output storage type
	err = uploadProxy(output.StorageType, output, fileContentsBuffer.Bytes())
	if err != nil {
		return err
	}

	return nil
}

func StoreFileJson(jsonContent []byte, output types.Output, webhook *types.Webhook) error {
	// Uploading the modified Excel file to the output storage type
	err := uploadProxy(output.StorageType, output, jsonContent)
	if err != nil {
		return err
	}

	return nil
}