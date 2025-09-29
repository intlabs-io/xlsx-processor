package types

type RequestBodyPaginate struct {
	Input  Input  `json:"input" validate:"required"`
	Output Output `json:"output" validate:"required"`
}

type PaginationResult struct {
	Message    string           `json:"message"`
	Attributes Attributes `json:"attributes"`
	TotalPages int              `json:"totalPages"`
}

type RequestBodyTransform struct {
	Input       Input    `json:"input" validate:"required"`
	Output      Output   `json:"output" validate:"required"`
	Rules       []Rule   `json:"rules" validate:"required"`
	Webhook     *Webhook `json:"webhook,omitempty"`
}
