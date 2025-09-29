package routes

import (
	"net/http"

	"xlsx-processor/pkg/types"
	"xlsx-processor/storage"
	"xlsx-processor/transform"

	"github.com/gin-gonic/gin"
)

func Transform(c *gin.Context) {
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
		Downloading the file from the input storage type
	*/
	f, err := storage.GetFile(input.StorageType, input)
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

	/*
		Storing the file in the output storage type
	*/
	err = storage.StoreFile(f, output, webhook)
	if err != nil {
		sendError(c, http.StatusInternalServerError, err, webhook)
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"message": "File transformed successfully"})
	return
}