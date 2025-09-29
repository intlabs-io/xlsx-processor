package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"xlsx-processor/pkg/sheet"
	"xlsx-processor/pkg/types"
	"xlsx-processor/storage"
	"xlsx-processor/transform"

	"github.com/gin-gonic/gin"
)

/*
TransformJson is a function that transforms a JSON file into an excelize file and then executes the rules on it.
*/
func TransformJson(c *gin.Context) {
	/*
		Request Body
	*/
	var requestData types.RequestBodyTransform
	err := bindAndValidate(c, &requestData)
	if err != nil {
		sendError(c, http.StatusBadRequest, err, nil)
		return
	}

	rules := requestData.Rules
	input := requestData.Input
	output := requestData.Output
	webhook := requestData.Webhook

	/*
		Check if the file is a JSON file
	*/
	if !strings.HasSuffix(strings.ToLower(input.Reference.Prefix), ".json") {
		sendError(c, http.StatusBadRequest, fmt.Errorf("file is not a JSON file"), nil)
		return
	}

	/*
		Downloading the file from the input storage type
	*/
	f, err := storage.GetFileFromJson(input.StorageType, input)
	if err != nil {
		sendError(c, http.StatusInternalServerError, err, webhook)
		return
	}

	/*
		Executing the rules
	*/
	rulesExecutor := transform.MakeRulesExecutor(f, rules)
	transformErr := rulesExecutor.Execute()
	if transformErr != nil {
		sendTransformError(c, http.StatusInternalServerError, transformErr, webhook)
		return
	}

	sheetContent, err := sheet.ParseSheetToCsv(f, nil)
	if err != nil {
		sendError(c, http.StatusInternalServerError, err, webhook)
		return
	}

	/*
		Storing the file in the output storage type
	*/
	jsonBytes, err := json.Marshal(sheetContent)
	if err != nil {
		sendError(c, http.StatusInternalServerError, err, webhook)
		return
	}
	err = storage.StoreFileJson(jsonBytes, output, webhook)
	if err != nil {
		sendError(c, http.StatusInternalServerError, err, webhook)
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"message": "File transformed successfully"})
	return
}
