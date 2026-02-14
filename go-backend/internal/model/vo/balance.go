package vo

type BalanceVO struct {
	Balance       float64 `json:"balance"`
	TotalSpending float64 `json:"totalSpending"`
	TotalRecharge float64 `json:"totalRecharge"`
}
