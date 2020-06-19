package sneakerWorkers

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/streadway/amqp"

	"quote/initializers"
	. "quote/models"
	"quote/utils"
)

func (worker Worker) SubQuoteBuildWorker(payloadJson *[]byte) (err error) {
	start := time.Now().UnixNano()
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
	if db.Where("`source` in (?)", []string{origin.Source, "local"}).Where("base_id = ?", origin.QuoteId).Group("quote_id").Find(&quotes).RecordNotFound() {
		return
	}
	db.DbRollback()
	for _, q := range quotes {
		sub, err := subQuote(&origin, &q)
		if err != nil {
			continue
		}
		if !q.IsLegal() && (sub.IsLegal() || sub.IsAnchored()) {
			createSubQuote(&sub)
		}
	}
	worker.LogInfo(" payload: ", payload, ", time:", (time.Now().UnixNano()-start)/1000000, " ms")
	return
}

func subQuote(origin, q *Quote) (Quote, error) {
	var subQuote Quote
	m := utils.DbBegin()
	defer m.DbRollback()
	if m.Where("type = ?", origin.Type).
		Where("base_id = ?", origin.BaseId).
		Where("quote_id = ?", q.QuoteId).
		Where("market_id = ?", origin.MarketId).
		Where("source = ?", origin.Source).
		First(&subQuote).RecordNotFound() {
		subQuote.Type = origin.Type
		subQuote.BaseId = origin.BaseId
		subQuote.MarketId = origin.MarketId
		subQuote.Source = origin.Source
		subQuote.QuoteId = q.QuoteId
		subQuote.Price = origin.Price.Mul(q.Price)
		subQuote.Timestamp = origin.Timestamp
		m.Save(&subQuote)
	} else {
		if subQuote.Timestamp >= origin.Timestamp {
			return subQuote, fmt.Errorf("Already have.")
		}
	}
	subQuote.QuoteCurrency = q.QuoteCurrency
	if subQuote.Price.Equal(origin.Price.Mul(q.Price)) {
		return subQuote, fmt.Errorf("Already have.")
	}
	subQuote.Price = origin.Price.Mul(q.Price)
	subQuote.Timestamp = origin.Timestamp
	subQuote.SaveToRedis()
	subQuote.NotifyQuote()
	return subQuote, nil
}

func createSubQuote(quote *Quote) {
	b, err := json.Marshal(struct {
		Id int `json:"id"`
	}{
		Id: quote.Id,
	})
	if err != nil {
		log.Println(err)
	}
	err = initializers.PublishMessageWithRouteKey("quote.default", "quote.sub.build", "text/plain", &b, amqp.Table{}, amqp.Persistent)
	if err != nil {
		log.Println(err)
	}
}
