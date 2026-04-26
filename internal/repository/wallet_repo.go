package repository

// redisWalletRepo implements WalletRepository using Redis hashes.
// Key: "wallet:{wallet_id}:stocks"  →  Hash{ stock_name: quantity }
type redisWalletRepo struct{}
