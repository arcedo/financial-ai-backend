package types

type PositionSummary struct {
	TotalBuys  float64 `json:"total_buys"`
	TotalSells float64 `json:"total_sells"`
	TotalSaves float64 `json:"total_saves"`
	TotalEntry float64 `json:"total_entry"`
	NetMarket  float64 `json:"net_market"`  // TotalBuys - TotalSells
	NetBalance float64 `json:"net_balance"` // NetMarket + TotalSaves + TotalEntry
}

type Insights struct {
	User         User            `json:"user"`
	Transactions []Transaction   `json:"transactions"`
	Position     PositionSummary `json:"position"`
}
