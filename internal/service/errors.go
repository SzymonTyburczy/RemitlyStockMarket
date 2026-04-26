package service

import "errors"

// Sentinel errors returned by services — handlers map these to HTTP status codes.
var (
	ErrStockNotFound      = errors.New("stock not found")
	ErrInsufficientBank   = errors.New("insufficient stock in bank")
	ErrInsufficientWallet = errors.New("insufficient stock in wallet")
	ErrInvalidOperation   = errors.New("invalid operation type, must be buy or sell")
)
