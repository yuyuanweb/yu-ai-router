package dto

type RequestLogQueryRequest struct {
	PageNum      int64          `json:"pageNum"`
	PageSize     int64          `json:"pageSize"`
	SortField    string         `json:"sortField"`
	SortOrder    string         `json:"sortOrder"`
	UserID       *FlexibleInt64 `json:"userId"`
	RequestModel string         `json:"requestModel"`
	RequestType  string         `json:"requestType"`
	Source       string         `json:"source"`
	Status       string         `json:"status"`
	StartDate    string         `json:"startDate"`
	EndDate      string         `json:"endDate"`
}
