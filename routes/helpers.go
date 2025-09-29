package routes

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"xlsx-processor/pkg/httphelper"
	"xlsx-processor/pkg/types"
)

func sendError(c *gin.Context, status int, err error, webhook *types.Webhook) {
	fmt.Println("Error transforming data", err.Error())
	if webhook != nil && webhook.Url != "" {
		payload := webhook.Payload
		payload.Status = "ERROR"
		payload.Msg = err.Error()
		err = httphelper.SendPostRequest(payload, webhook.Url, webhook.ResponseToken)
		if err != nil {
			fmt.Println("Error sending webhook", err.Error())
		}
	}
	c.JSON(status, gin.H{"message": err.Error()})
}

func sendTransformError(c *gin.Context, status int, transformErr *types.TransformError, webhook *types.Webhook) {
	fmt.Println("Error transforming data", transformErr.Message)
	if webhook != nil && webhook.Url != "" {
		payload := webhook.Payload
		payload.Status = "ERROR"
		payload.Msg = transformErr.Message
		err := httphelper.SendPostRequest(payload, webhook.Url, webhook.ResponseToken)
		if err != nil {
			fmt.Println("Error sending webhook", err.Error())
			transformErr = &types.TransformError{
				Message: err.Error(),
			}
		}
	}
	c.JSON(status, gin.H{"transformError": transformErr})
}

func bindAndValidate(c *gin.Context, requestData any) error {
	validate := validator.New()
	err := c.ShouldBindJSON(requestData)
	if err != nil {
		return err
	}
	err = validate.Struct(requestData)
	if err != nil {
		return err
	}
	return nil
}
