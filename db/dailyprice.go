package db

import (
	"fmt"
	"stockcrawler/dto"
)

func SaveDailyPrice(dp *dto.DailyPriceInfo) {

	sql := "INSERT INTO stocks.dailyprice(symbol, openprice, closeprice, maxprice, minprice, volume, turnover, increaserate, increaseprice) VALUES ('%s', %v, %v, %v, %v, %v, %v, %v, %v)"
	sql = fmt.Sprintf(sql, dp.Symbol, dp.OpenPrice, dp.ClosePrice, dp.MaxPrice, dp.MinPrice, dp.Volume, dp.Turnover, dp.IncreaseRate, dp.IncreasePrice)
	Execute(sql)
}
