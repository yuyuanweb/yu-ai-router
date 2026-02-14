package vo

import "time"

type PluginConfigVO struct {
	ID          int64     `json:"id,string"`
	PluginKey   string    `json:"pluginKey"`
	PluginName  string    `json:"pluginName"`
	PluginType  string    `json:"pluginType"`
	Description string    `json:"description"`
	Config      string    `json:"config"`
	Status      string    `json:"status"`
	Priority    int       `json:"priority"`
	CreateTime  time.Time `json:"createTime"`
	UpdateTime  time.Time `json:"updateTime"`
}

type PluginExecuteVO struct {
	Success      bool           `json:"success"`
	PluginKey    string         `json:"pluginKey"`
	Content      string         `json:"content"`
	ErrorMessage string         `json:"errorMessage"`
	Duration     int64          `json:"duration"`
	Data         map[string]any `json:"data"`
}
