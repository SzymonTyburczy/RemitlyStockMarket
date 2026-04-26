package repository

// redisBankRepo implements BankRepository using Redis hashes.
// Key: "bank:stocks"  →  Hash{ stock_name: quantity }
type redisBankRepo struct{}
