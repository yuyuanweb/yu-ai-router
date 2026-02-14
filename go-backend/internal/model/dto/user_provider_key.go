package dto

type UserProviderKeyAddRequest struct {
	ProviderID *FlexibleInt64 `json:"providerId"`
	APIKey     string         `json:"apiKey"`
}

type UserProviderKeyUpdateRequest struct {
	ID     *FlexibleInt64 `json:"id"`
	APIKey string         `json:"apiKey"`
	Status string         `json:"status"`
}
