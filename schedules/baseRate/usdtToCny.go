package baseRate

import (
	"strings"
	"time"

	"github.com/shopspring/decimal"

	. "quote/models"
	"quote/utils"
)

func UsdtToCny() {
	price := decimal.NewFromFloat(1)
	db := utils.DbBegin()
	defer db.DbRollback()
	var cny Currency
	db.FirstOrInit(&cny, map[string]interface{}{
		"key":     "CNY",
		"symbol":  "cny",
		"source":  "local",
		"visible": true,
	})
	db.Save(&cny)
	var usdts []Currency
	db.Where("symbol = ?", "usdt").Find(&usdts)
	for _, u := range usdts {
		var quote Quote
		db.FirstOrInit(&quote, map[string]interface{}{
			"base_id":  u.Id,
			"quote_id": cny.Id,
			"type":     "Quotes::" + strings.Title(u.Source),
		})
		quote.Timestamp = time.Now().UnixNano() / 1000000
		quote.Price = price
		db.Save(&quote)
	}
	db.DbCommit()
}
