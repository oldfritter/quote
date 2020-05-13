package sneakerWorkers

import (
	"encoding/json"
	"log"
	"time"

	"github.com/streadway/amqp"

	"quote/initializers"
	. "quote/models"
	"quote/utils"
)

func (worker Worker) QuoteBuildWorker(payloadJson *[]byte) (err error) {
	var payload struct {
		Id int `json:"id"`
	}
	json.Unmarshal([]byte(*payloadJson), &payload)
	db := utils.DbBegin()
	defer db.DbRollback()
	var origin Quote
	if db.Where("id = ?", payload.Id).First(&origin).RecordNotFound() {
		return
	}
	var quotes []Quote
	if db.Where("`source` in (?)", []string{"local", origin.Source}).Where("base_id = ?", origin.QuoteId).Find(&quotes).RecordNotFound() {
		return
	}
	var subQuotes []Quote
	for _, q := range quotes {
		m := utils.DbBegin()
		defer m.DbRollback()
		var subQuote Quote
		m.FirstOrInit(&subQuote, map[string]interface{}{
			"type":      origin.Type,
			"base_id":   origin.BaseId,
			"quote_id":  q.QuoteId,
			"market_id": origin.MarketId,
			"source":    origin.Source,
		})
		subQuote.QuoteCurrency = q.QuoteCurrency
		subQuote.Price = origin.Price.Mul(q.Price)
		subQuote.Timestamp = time.Now().UnixNano() / 1000000
		m.Save(&subQuote)
		m.DbCommit()
		subQuotes = append(subQuotes, subQuote)
	}
	db.DbCommit()
	if origin.QuoteCurrency.Source != "local" {
		for _, q := range subQuotes {
			if q.QuoteCurrency.Source != "local" {
				createSubQuote(&q)
			}
		}
	}
	return
}

func createSubQuote(quote *Quote) {
	b, err := json.Marshal(map[string]int{"id": quote.Id})
	if err != nil {
		log.Println(err)
	}
	err = initializers.PublishMessageWithRouteKey("quote.default", "quote.build", "text/plain", &b, amqp.Table{}, amqp.Persistent)
	if err != nil {
		log.Println(err)
	}
}
