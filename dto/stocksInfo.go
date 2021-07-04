package dto

import "time"

// StocksInfo 股票信息
type StocksInfo struct {
	Symbol         string
	Name           string
	Market         string
	Category       string
	Business       string
	Concepts       string
	CreateTime     time.Time
	LastUpdateTime time.Time
}
