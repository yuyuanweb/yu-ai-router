package dto

type ModelAddRequest struct {
	ProviderID     *FlexibleInt64 `json:"providerId"`
	ModelKey       string         `json:"modelKey"`
	ModelName      string         `json:"modelName"`
	ModelType      string         `json:"modelType"`
	Description    string         `json:"description"`
	ContextLength  *int           `json:"contextLength"`
	InputPrice     *float64       `json:"inputPrice"`
	OutputPrice    *float64       `json:"outputPrice"`
	Priority       *int           `json:"priority"`
	DefaultTimeout *int           `json:"defaultTimeout"`
	Capabilities   string         `json:"capabilities"`
}

type ModelUpdateRequest struct {
	ID             *FlexibleInt64 `json:"id"`
	ModelName      *string        `json:"modelName"`
	Description    *string        `json:"description"`
	ContextLength  *int           `json:"contextLength"`
	InputPrice     *float64       `json:"inputPrice"`
	OutputPrice    *float64       `json:"outputPrice"`
	Status         *string        `json:"status"`
	Priority       *int           `json:"priority"`
	DefaultTimeout *int           `json:"defaultTimeout"`
	Capabilities   *string        `json:"capabilities"`
}

type ModelQueryRequest struct {
	PageNum    int64          `json:"pageNum"`
	PageSize   int64          `json:"pageSize"`
	SortField  string         `json:"sortField"`
	SortOrder  string         `json:"sortOrder"`
	ProviderID *FlexibleInt64 `json:"providerId"`
	ModelKey   string         `json:"modelKey"`
	ModelName  string         `json:"modelName"`
	ModelType  string         `json:"modelType"`
	Status     string         `json:"status"`
}
