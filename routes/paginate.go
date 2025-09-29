package routes

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"xlsx-processor/pkg/sheet"
	"xlsx-processor/pkg/types"
	"xlsx-processor/storage"
)

/*
Paginate a multi sheet file into multiple single sheet json files
*/
func Paginate(c *gin.Context) {
	/*
		Request Body
	*/
	var requestData types.RequestBodyPaginate // Reusing the same request structure
	err := bindAndValidate(c, &requestData)
	if err != nil {
		sendError(c, http.StatusBadRequest, err, nil)
		return
	}

	input := requestData.Input
	output := requestData.Output

	/*
		Downloading the file from the input storage type
	*/
	f, err := storage.GetFile(input.StorageType, input)
	if err != nil {
		sendError(c, http.StatusInternalServerError, err, nil)
		return
	}

	/*
		Get all sheet names from the original document
	*/
	sheetNames := f.GetSheetList()
	sheetCount := len(sheetNames)

	if sheetCount == 0 {
		sendError(c, http.StatusBadRequest, fmt.Errorf("Excel file has no sheets"), nil)
		return
	}

	if sheetCount > 1000 { // Reasonable limit to prevent excessive processing
		sendError(c, http.StatusBadRequest, fmt.Errorf("Excel file has too many sheets (%d). Maximum supported is 1,000 sheets", sheetCount), nil)
		return
	}

	/*
		Extract and upload each sheet individually
	*/
	uploadedSheets := 0

	// Get global attributes once for the entire file
	_, attributes, err := sheet.ParseSheetsToCsvAndAttributes(f)
	if err != nil {
		sendError(c, http.StatusInternalServerError, err, nil)
		return
	}

	for sheetIndex, sheetName := range sheetNames {
		// Convert 0-indexed to 1-indexed for display
		displaySheetNum := sheetIndex + 1

		// Parse individual sheet
		singleSheet, err := sheet.ParseSheetToCsv(f, &sheetName)
		if err != nil {
			sendError(c, http.StatusInternalServerError, err, nil)
			return
		}

		// Create output reference for this specific sheet
		sheetOutput := output
		sheetOutput.Reference.Prefix = fmt.Sprintf("%s/pages/%d.json", output.Reference.Prefix, displaySheetNum)

		jsonBytes, err := json.Marshal(singleSheet)
		if err != nil {
			sendError(c, http.StatusInternalServerError, err, nil)
			return
		}
		// Store the individual sheet file
		err = storage.StoreFileJson(jsonBytes, sheetOutput, nil)
		if err != nil {
			sendError(c, http.StatusInternalServerError, fmt.Errorf("failed to upload sheet %d (%s): %v", displaySheetNum, sheetName, err), nil)
			return
		}

		uploadedSheets++
	}

	/*
		Return pagination result
	*/
	result := types.PaginationResult{
		Message:    fmt.Sprintf("Success: %d sheets paginated", uploadedSheets),
		Attributes: *attributes,
		TotalPages: uploadedSheets,
	}

	c.JSON(http.StatusAccepted, result)
}
