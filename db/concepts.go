package db

import (
	"fmt"
	"stockcrawler/dto"
)

func CountConceptBySymbol(symbol string) (counter int, err error) {
	sql := "select count(id) as counter from stocks.concept where symbol = '%s'"
	sql = fmt.Sprintf(sql, symbol)
	rows, err := Select(sql)

	if err != nil {
		return 0, err
	}

	if rows.Next() {
		rows.Scan(&counter)
	}

	return
}

func DeleteConceptBySymbol(symbol string) error {
	sql := "delete from stocks.concept where symbol = '%s'"
	sql = fmt.Sprintf(sql, symbol)
	_, _, err := Execute(sql)

	if err != nil {
		return err
	}

	return nil
}

// SaveConceptInfo 保存股票概念基础数据信息
func SaveConceptInfo(conceptInfo *dto.ConceptInfo) error {
	sql := "INSERT INTO stocks.concept(symbol, conceptname) VALUES ('%s', '%s')"
	sql = fmt.Sprintf(sql, conceptInfo.Symbol, conceptInfo.ConceptName)
	_, _, err := Execute(sql)
	return err
}
