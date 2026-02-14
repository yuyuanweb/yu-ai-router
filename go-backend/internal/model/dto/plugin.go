package dto

type PluginExecuteRequest struct {
	PluginKey string         `json:"pluginKey"`
	Input     string         `json:"input"`
	FileURL   string         `json:"fileUrl"`
	FileBytes []byte         `json:"-"`
	FileType  string         `json:"fileType"`
	Params    map[string]any `json:"params"`
}

type PluginUpdateRequest struct {
	ID          *FlexibleInt64 `json:"id"`
	PluginName  *string        `json:"pluginName"`
	Description *string        `json:"description"`
	Config      *string        `json:"config"`
	Status      *string        `json:"status"`
	Priority    *int           `json:"priority"`
}
