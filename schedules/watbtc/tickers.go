package watbtc

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/shopspring/decimal"
	"github.com/streadway/amqp"

	"quote/initializers"
	. "quote/models"
	"quote/utils"
)

func GetWatbtcTickers() {

	ctx, cancelFun := context.WithTimeout(context.Background(), time.Second*10)
	defer cancelFun()
	req, _ := http.NewRequest(http.MethodGet, "https://watbtc.vip/api/v2/tickers", nil)
	resp, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil {
		log.Println(err)
		return
	}
	body, err := ioutil.ReadAll(resp.Body)

	var result struct {
		Body []struct {
			At     int64  `json:"at"`
			Name   string `json:"name"`
			Code   string `json:"code"`
			Quote  string `json:"quote"`
			Ticker struct {
				Last decimal.Decimal `json:"last"`
			} `json:"ticker"`
		} `json:"body"`
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Println(err)
		return
	}
	db := utils.DbBegin()
	defer db.DbRollback()
	for _, m := range result.Body {
		var market Market
		db.Where("source in (?)", []string{"watbtc", "local"}).FirstOrInit(&market, map[string]interface{}{"name": strings.ToUpper(m.Name), "symbol": strings.ToLower(m.Name)})
		var base, quoteCurrency Currency
		db.Where("source in (?)", []string{"watbtc", "local"}).FirstOrInit(&base, map[string]interface{}{"symbol": strings.ToLower(m.Code), "key": strings.ToUpper(m.Code), "visible": true})
		db.Where("source in (?)", []string{"watbtc", "local"}).FirstOrInit(&quoteCurrency, map[string]interface{}{"symbol": strings.ToLower(m.Quote), "key": strings.ToUpper(m.Quote), "visible": true})
		db.Save(&base)
		db.Save(&quoteCurrency)
		market.BaseId = base.Id
		market.QuoteId = quoteCurrency.Id
		market.Visible = true
		market.Source = "watbtc"
		db.Save(&market)
		var quote Quote
		db.Where("source in (?)", []string{"watbtc", "local"}).FirstOrInit(&quote, map[string]interface{}{"type": "Quotes::Watbtc", "base_id": base.Id, "quote_id": quoteCurrency.Id, "market_id": market.Id})
		quote.Price = m.Ticker.Last
		quote.Timestamp = m.At * 1000
		db.Save(&quote)
		createSubQuote(&quote)
	}
	db.DbCommit()
}

func createSubQuote(quote *Quote) {
	b, err := json.Marshal(map[string]int{"id": quote.Id})
	if err != nil {
		log.Println(err)
	}
	err = initializers.PublishMessageWithRouteKey("quote.default", "quote.sub.build", "text/plain", &b, amqp.Table{}, amqp.Persistent)
	if err != nil {
		log.Println(err)
	}
}
