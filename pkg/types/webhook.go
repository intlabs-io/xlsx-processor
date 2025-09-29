package types

type Webhook struct {
	Url           string         `json:"url"`
	Payload       WebhookPayload `json:"payload"`
	ResponseToken string         `json:"responseToken"`
}

type WebhookPayload struct {
	Msg          string `json:"msg"`
	BrowserTabID string `json:"browserTabID,omitempty"`
	Uuid         string `json:"uuid"`
	UserId       string `json:"userId"`
	S3Bucket     string `json:"s3Bucket"`
	S3Key        string `json:"s3Key"`
	SourceId     string `json:"sourceId"`
	Status       string `json:"status"`
}