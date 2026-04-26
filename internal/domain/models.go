package domain

// Stock represents a named stock with a quantity.
type Stock struct {
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
}

// Wallet represents a wallet identified by ID, holding a set of stocks.
type Wallet struct {
	ID     string  `json:"id"`
	Stocks []Stock `json:"stocks"`
}

// OperationType is either "buy" or "sell".
type OperationType string

const (
	Buy  OperationType = "buy"
	Sell OperationType = "sell"
)

// LogEntry represents a single audit log record.
type LogEntry struct {
	Type      OperationType `json:"type"`
	WalletID  string        `json:"wallet_id"`
	StockName string        `json:"stock_name"`
}

// TradeRequest is the body for POST /wallets/{id}/stocks/{stock}.
type TradeRequest struct {
	Type OperationType `json:"type"`
}
