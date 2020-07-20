package quote

import (
	"encoding/json"
	"fmt"

	"github.com/gomodule/redigo/redis"

	. "quote/models"
	"quote/utils"
)

func SaveDataFromRedis() {
	db := utils.Db
	var markets []Market
	db.Where("visible = ?", true).Find(&markets)

	dataRedis := utils.GetRedisConn("data")
	defer dataRedis.Close()

	for _, market := range markets {
		keys, err := redis.Strings(dataRedis.Do("KEYS", fmt.Sprintf("Quotes:%d:*", market.Id)))
		if err != nil {
		}
		for _, key := range keys {

			qByte, err := redis.Bytes(dataRedis.Do("GET", key))
			if err != nil {
			}
			var quote Quote
			var market Market
			var simple SimpleQuote
			var baseC Currency
			var quoteCs []Currency
			json.Unmarshal(qByte, &simple)
			if db.Where("source = ?", simple.Source).Where("symbol = ?", simple.Market).First(&market).RecordNotFound() {
				return
			}
			if db.Where("source = ?", simple.Source).Where("symbol = ?", simple.Base).First(&baseC).RecordNotFound() {
				return
			}
			if db.Where("source in (?)", []string{simple.Source, "local"}).Where("symbol = ?", simple.Quote).Find(&quoteCs).RecordNotFound() {
				return
			}
			for _, quoteC := range quoteCs {
				m := utils.DbBegin()
				defer m.DbRollback()
				if m.Where("source = ?", simple.Source).Where("base_id = ?", baseC.Id).Where("quote_id = ?", quoteC.Id).Where("market_id = ?", market.Id).First(&quote).RecordNotFound() {
					quote.Source = simple.Source
					quote.BaseId = baseC.Id
					quote.QuoteId = quoteC.Id
					quote.MarketId = market.Id
				}
				quote.Price = simple.Price
				quote.Timestamp = simple.Timestamp
				m.Save(&quote)
				m.DbCommit()
			}
		}
	}
}
