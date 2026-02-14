package vo

type QuotaVO struct {
	TokenQuota     int64 `json:"tokenQuota,string"`
	UsedTokens     int64 `json:"usedTokens,string"`
	RemainingQuota int64 `json:"remainingQuota,string"`
}

type UserAnalysisVO struct {
	UserID          int64   `json:"userId,string"`
	UserAccount     string  `json:"userAccount"`
	UserName        string  `json:"userName"`
	UserStatus      string  `json:"userStatus"`
	UserRole        string  `json:"userRole"`
	TokenQuota      int64   `json:"tokenQuota,string"`
	UsedTokens      int64   `json:"usedTokens,string"`
	RemainingQuota  int64   `json:"remainingQuota,string"`
	TotalRequests   int64   `json:"totalRequests,string"`
	SuccessRequests int64   `json:"successRequests,string"`
	TotalTokens     int64   `json:"totalTokens,string"`
	TotalCost       float64 `json:"totalCost"`
	TodayCost       float64 `json:"todayCost"`
}
