package main

import (
	"stockcrawler/db"
	"stockcrawler/service"
	"time"
)

func init() {

	// 创建数据库链接
	db.NewConnection("stocks", "mysql", "dbo:caecaodb@tcp(192.168.3.168:3306)/stocks?charset=utf8&parseTime=true&loc=Local")
	// 将所有股票基础信息加载到内存
	db.SelectAllStocksInfo()
}

func main() {

	timer := time.Tick(3600 * 1e9)

	for next := range timer {
		hour := next.Hour()
		if hour == 20 {
			service.GetStocksInfoFromIFeng()
		}
	}
}
