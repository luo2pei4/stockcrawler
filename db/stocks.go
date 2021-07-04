package db

import (
	"fmt"
	"stockcrawler/dto"
)

var StocksMap map[string]*dto.StocksInfo

// SelectAllStocksInfo 查询stocks表的所有数据
func SelectAllStocksInfo() {

	sql := "SELECT symbol, name, market, category, business, concepts, createtime, lastupdatetime FROM stocks.stocks"
	rows, err := Select(sql)

	if err != nil {
		fmt.Printf("stocks.SelectAll error. %v", err.Error())
	}

	StocksMap = make(map[string]*dto.StocksInfo)

	for rows.Next() {
		dto := dto.StocksInfo{}
		rows.Scan(&dto.Symbol, &dto.Name, &dto.Market, &dto.Category, &dto.Business, &dto.Concepts, &dto.CreateTime, &dto.LastUpdateTime)
		StocksMap[dto.Symbol] = &dto
	}
}

// SaveStockInfo 保存股票基础数据信息
func SaveStockInfo(stocksInfo *dto.StocksInfo) error {

	sql := "INSERT INTO stocks.stocks(symbol, name, market, category, business, concepts) VALUES ('%s', '%s', '%s', '%s', '%s', '%s')"
	sql = fmt.Sprintf(sql, stocksInfo.Symbol, stocksInfo.Name, stocksInfo.Market, stocksInfo.Category, stocksInfo.Business, stocksInfo.Concepts)
	_, _, err := Execute(sql)
	return err
}

// UpdateStockInfo 更新股票基础信息表
func UpdateStockInfo(stocksInfo *dto.StocksInfo) error {
	sql := "UPDATE stocks.stocks SET name='%s', market='%s', category='%s', business='%s', concepts='%s', lastupdatetime=current_timestamp() WHERE symbol='%s'"
	sql = fmt.Sprintf(sql, stocksInfo.Name, stocksInfo.Market, stocksInfo.Category, stocksInfo.Business, stocksInfo.Concepts, stocksInfo.Symbol)
	_, _, err := Execute(sql)
	return err
}
