package baseRate

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/shopspring/decimal"

	. "quote/models"
	"quote/utils"
)

func UsdtToCny() {
	url := "https://markets.money.cnn.com/common/modules/iframe/currencyConverter.asp?convert=1&amount=1&base=USD&quote=CNY"
	ctx, cancelFun := context.WithTimeout(context.Background(), time.Second*10)
	defer cancelFun()
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	resp, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil {
		log.Println(err)
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	p := strings.TrimLeft(string(body), "= ")
	p = strings.TrimRight(p, " CNY")
	price, _ := decimal.NewFromString(p)
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
