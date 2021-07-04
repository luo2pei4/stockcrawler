package dto

import "time"

// ConceptInfo 股票概念信息
type ConceptInfo struct {
	Id          int64
	Symbol      string
	ConceptName string
	CreateTime  time.Time
}
