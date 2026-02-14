package dto

type ProviderAddRequest struct {
	ProviderName string `json:"providerName"`
	DisplayName  string `json:"displayName"`
	BaseURL      string `json:"baseUrl"`
	APIKey       string `json:"apiKey"`
	Priority     *int   `json:"priority"`
	Config       string `json:"config"`
}

type ProviderUpdateRequest struct {
	ID          *FlexibleInt64 `json:"id"`
	DisplayName *string        `json:"displayName"`
	BaseURL     *string        `json:"baseUrl"`
	APIKey      *string        `json:"apiKey"`
	Status      *string        `json:"status"`
	Priority    *int           `json:"priority"`
	Config      *string        `json:"config"`
}

type ProviderQueryRequest struct {
	PageNum      int64  `json:"pageNum"`
	PageSize     int64  `json:"pageSize"`
	SortField    string `json:"sortField"`
	SortOrder    string `json:"sortOrder"`
	ProviderName string `json:"providerName"`
	DisplayName  string `json:"displayName"`
	Status       string `json:"status"`
	HealthStatus string `json:"healthStatus"`
}
