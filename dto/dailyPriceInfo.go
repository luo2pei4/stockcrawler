package dto

import "time"

type DailyPriceInfo struct {
	Id            int64
	Symbol        string
	OpenPrice     float32
	ClosePrice    float32
	MaxPrice      float32
	MinPrice      float32
	Volume        int32
	Turnover      int32
	IncreaseRate  float32
	IncreasePrice float32
	TradeDate     time.Time
}
