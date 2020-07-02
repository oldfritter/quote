package quote

import (
	"encoding/json"
	"fmt"

	"github.com/gomodule/redigo/redis"

	. "quote/models"
	"quote/utils"
)

func SaveDataFromRedis() {
	m := utils.DbBegin()
	defer m.DbRollback()
	var markets []Market
	m.Where("visible = ?", true).Find(&markets)

	dataRedis := utils.GetRedisConn("data")
	defer dataRedis.Close()

	for _, market := range markets {
		keys, err := redis.Strings(dataRedis.Do("KEYS", fmt.Sprintf("Quotes:%v:*", market.Id)))
		if err != nil {
		}
		for _, key := range keys {
			qByte, err := redis.Bytes(dataRedis.Do("GET", key))
			if err != nil {
			}
			var quote Quote
			var market Market
			var simple SimpleQuote
			var baseC, quoteC Currency
			json.Unmarshal(qByte, &simple)
			if m.Where("source = ?", simple.Source).Where("symbol = ?", simple.Market).First(&market).RecordNotFound() {
				return
			}
			if m.Where("source = ?", simple.Source).Where("symbol = ?", simple.Base).First(&baseC).RecordNotFound() {
				return
			}
			if m.Where("source = ?", simple.Source).Where("symbol = ?", simple.Quote).First(&quoteC).RecordNotFound() {
				return
			}
			m.Where("source = ?", simple.Source).Where("base_id = ?", baseC.Id).Where("quote_id = ?", quoteC.Id).Where("market_id = ?", market.Id).FirstOrInit(&quote)
			quote.Price = simple.Price
			quote.Timestamp = simple.Timestamp
			m.Debug().Save(&quote)
		}
	}
	m.DbCommit()
}
