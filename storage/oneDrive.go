package storage

import (
	"fmt"
	"io"
	"net/http"
	"xlsx-processor/pkg/types"
)

// Downloading the file as bytes from OneDrive
func DownloadFromOneDrive(input types.Input) (fileContents []byte, err error) {
	// Create a new HTTP client
	client := &http.Client{}

	// Create a new request with bearer token authorization header
	req, err := http.NewRequest("GET", "https://graph.microsoft.com/v1.0/me/drive/items/"+input.Reference.Id+"/content", nil)
	if err != nil {
		return nil, err
	}
	if input.Credential.Secrets.AccessToken != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", input.Credential.Secrets.AccessToken))
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	fileContents, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return fileContents, nil
}
