package service

import (
	"errors"
	"fmt"
	"net/http"
	"stockcrawler/db"
	"stockcrawler/dto"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// GetStocksInfoFromIFeng 从凤凰网财经频道获取股票信息
func GetStocksInfoFromIFeng() {

	pageNo := 0
	hasNextPage := true

	for hasNextPage {

		pageNo++
		path := "https://app.finance.ifeng.com/list/stock.php?t=hs&f=symbol&o=asc&p="

		doc, err := getPageDoc(path + strconv.FormatInt(int64(pageNo), 10))

		if err != nil {
			fmt.Println(err.Error())
			return
		}

		hasNextPage = parseListData(doc)
	}
}

// getPageDoc 获取网页的Document
func getPageDoc(pageURL string) (doc *goquery.Document, err error) {

	retry := true
	counter := 1
	var response *http.Response

	for retry {

		response, err = http.Get(pageURL)

		if err != nil {

			if counter == 30 {
				return nil, err
			}

			retry = true
			time.Sleep(1e6)
			fmt.Println("retry...", counter)
			counter++

		} else {

			retry = false
		}
	}

	statusCode := response.StatusCode

	if statusCode != 200 {
		err = errors.New("Access error, status code is " + strconv.FormatInt(int64(statusCode), 10))
		return nil, err
	}

	doc, err = goquery.NewDocumentFromReader(response.Body)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	return doc, nil
}

func parseListData(doc *goquery.Document) (hasNextPage bool) {

	selector := "body > div.main > div > div.block02 > div > table > tbody"

	doc.Find(selector).Each(func(i int, table *goquery.Selection) {

		rowSize := table.Find("tr").Size()

		table.Find("tr").Each(func(j int, tr *goquery.Selection) {

			// 忽略第一行的title
			if j != 0 {

				si := &dto.StocksInfo{}
				dp := &dto.DailyPriceInfo{}
				var cellValue string

				tr.Find("td").Each(func(k int, td *goquery.Selection) {

					cellValue = td.Text()

					switch k {
					case 0:
						si.Symbol = cellValue

						// 通过代码前三位设置市场和股票类别信息
						substr := cellValue[0:3]
						market, category := getMarketAndCategory(substr)
						si.Market = market
						si.Category = category
						dp.Symbol = cellValue

					case 1:
						si.Name = cellValue
						detailUrl, exist := td.Find("a").Attr("href")

						if exist {

							detailDoc, err := getPageDoc(detailUrl)

							if err == nil {
								// 获取股票所属行业和概念
								business, concepts := getBusinessAndConcept(detailDoc)
								si.Business = business
								si.Concepts = concepts
							} else {
								fmt.Printf("Access detail failed. error:%v\n", err.Error())
							}
						}

					case 2:
						closePrice, _ := strconv.ParseFloat(cellValue, 32)
						dp.ClosePrice = float32(closePrice)

					case 3:
						value := strings.Replace(cellValue, "%", "", 1)
						increaseRate, _ := strconv.ParseFloat(value, 32)
						dp.IncreaseRate = float32(increaseRate)

					case 4:
						increasePrice, _ := strconv.ParseFloat(cellValue, 32)
						dp.IncreasePrice = float32(increasePrice)

					case 5:
						value := strings.Replace(cellValue, "手", "", 1)
						volume, _ := strconv.ParseInt(value, 10, 32)
						dp.Volume = int32(volume)

					case 6:
						value := strings.Replace(cellValue, "万", "", 1)
						turnover, _ := strconv.ParseInt(value, 10, 32)
						dp.Turnover = int32(turnover)

					case 7:
						openPrice, _ := strconv.ParseFloat(cellValue, 32)
						dp.OpenPrice = float32(openPrice)

					case 9:
						minPrice, _ := strconv.ParseFloat(cellValue, 32)
						dp.MinPrice = float32(minPrice)

					case 10:
						maxPrice, _ := strconv.ParseFloat(cellValue, 32)
						dp.MaxPrice = float32(maxPrice)
					}
				})

				fmt.Printf("Stock info: %v\n", si)

				// 最后一行的描述文字中包含“下一页”的情况，返回true，否者返回false
				if j == rowSize-1 {
					if strings.Contains(cellValue, "下一页") {
						hasNextPage = true
					} else {
						hasNextPage = false
					}
				}

				// 新增股票信息入库
				if db.StocksMap[si.Symbol] == nil {
					db.SaveStockInfo(si)
					saveConcepts(si.Symbol, si.Concepts)
					db.StocksMap[si.Symbol] = si
				} else {
					// 判断概念信息是否有变更
					oldValue := db.StocksMap[si.Symbol].Concepts
					newValue := si.Concepts

					// 长度有变化或者新旧值不相等的情况，更新数据
					if (len(oldValue) != len(newValue)) || (newValue != oldValue) {
						fmt.Printf("Concept info of stock %s(%v) changed.", si.Symbol, si.Name)
						db.UpdateStockInfo(si)
						updateConcepts(si.Symbol, newValue)
						db.StocksMap[si.Symbol].Concepts = newValue
					}
				}

				db.SaveDailyPrice(dp)
			}
		})
	})

	return hasNextPage
}

// getMarketAndCategory 根据股票代码前三位判断股票所属市场的分类
func getMarketAndCategory(substr string) (market, category string) {

	switch substr {
	case "300":
		market = ""
		category = "创业板"
	case "600", "601", "603", "605":
		market = "沪市"
		category = "A股"
	case "900":
		market = "沪市"
		category = "B股"
	case "688":
		market = ""
		category = "科创板"
	case "000", "001":
		market = "深市"
		category = "A股"
	case "200":
		market = "深市"
		category = "B股"
	case "002":
		market = ""
		category = "中小板"
	case "730":
		market = "沪市"
		category = "新股申购"
	case "700":
		market = "沪市"
		category = "配股"
	case "080":
		market = "深市"
		category = "配股"
	case "580":
		market = "沪市"
		category = "权证"
	case "031":
		market = "沪市"
		category = "权证"
	default:
		market = ""
		category = ""
	}

	return
}

// getBusinessAndConcept 获取详细页面的股票行业和概念数据
func getBusinessAndConcept(doc *goquery.Document) (business string, concepts string) {

	businessSelector := "body > div.col.clearfix.bg > div.Right770.clearfix > div:nth-child(1) > div.box549 > div.picForme > table > tbody > tr:nth-child(3) > td > span:nth-child(1)"
	business = doc.Find(businessSelector).Text()
	if business != "" {
		business = strings.Replace(business, "所属行业：", "", 1)
	}

	conceptSelector := "body > div.col.clearfix.bg > div.Right770.clearfix > div:nth-child(1) > div.box549 > div.picForme > table > tbody > tr:nth-child(3) > td > span:nth-child(2)"
	concepts = doc.Find(conceptSelector).Text()
	concepts = strings.TrimSpace(concepts)
	if concepts != "" {
		concepts = strings.Replace(concepts, "所属概念：", "", 1)
		concepts = strings.Replace(concepts, " ", ",", -1)
	}

	return
}

// saveConcepts 保存股票概念数据
func saveConcepts(symbol, concepts string) {

	arr := strings.Split(concepts, ",")
	conceptInfo := &dto.ConceptInfo{}

	db.TxBegin()

	for _, conceptName := range arr {

		// 部分概念数据中存在两个空格的情况，替换为逗号后会存在空字符串的情况
		if conceptName == "" {
			continue
		}

		conceptInfo.Symbol = symbol
		conceptInfo.ConceptName = conceptName
		err := db.SaveConceptInfo(conceptInfo)

		if err != nil {
			db.Rollback()
			return
		}
	}

	db.TxCommit()
}

// updateConcepts 更新股票概念数据
func updateConcepts(symbol, concepts string) {

	counter, err := db.CountConceptBySymbol(symbol)

	if err != nil {
		fmt.Printf("Concept count error: %v", err.Error())
		return
	}

	if counter != 0 {
		db.DeleteConceptBySymbol(symbol)
	}

	saveConcepts(symbol, concepts)
}
