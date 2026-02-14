package vo

type TokenStatsVO struct {
	TotalTokens int64 `json:"totalTokens,string"`
}

type CostStatsVO struct {
	TotalCost float64 `json:"totalCost"`
	TodayCost float64 `json:"todayCost"`
}

type UserSummaryStatsVO struct {
	TotalTokens     int64   `json:"totalTokens,string"`
	TokenQuota      int64   `json:"tokenQuota,string"`
	UsedTokens      int64   `json:"usedTokens,string"`
	RemainingQuota  int64   `json:"remainingQuota,string"`
	TotalCost       float64 `json:"totalCost"`
	TodayCost       float64 `json:"todayCost"`
	TotalRequests   int64   `json:"totalRequests,string"`
	SuccessRequests int64   `json:"successRequests,string"`
}
