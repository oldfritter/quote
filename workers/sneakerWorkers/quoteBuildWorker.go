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
		Id    int `json:"id"`
		Level int `json:"level"`
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
		if payload.Level < 5 && (sub.IsLegal() || sub.IsAnchored()) {
			createSubQuote(&sub, payload.Level+1)
		}
	}
	worker.LogInfo(" payload: ", payload, ", time:", (time.Now().UnixNano()-start)/1000000, " ms")
	return
}

func subQuote(origin, q *Quote) (Quote, error) {
	subQuote := Quote{
		Type:          origin.Type,
		BaseId:        origin.BaseId,
		MarketId:      origin.MarketId,
		Source:        origin.Source,
		QuoteId:       q.QuoteId,
		Price:         origin.Price.Mul(q.Price),
		Timestamp:     origin.Timestamp,
		QuoteCurrency: q.QuoteCurrency,
		Price:         origin.Price.Mul(q.Price),
	}
	if subQuote.AlreadyHave() {
		return subQuote, fmt.Errorf("Already have.")
	}
	subQuote.SaveToRedis()
	subQuote.NotifyQuote()
	return subQuote, nil
}

func createSubQuote(quote *Quote, level int) {
	b, err := json.Marshal(struct {
		Id    int `json:"id"`
		Level int `json:"level"`
	}{
		Id:    quote.Id,
		Level: level,
	})
	if err != nil {
		log.Println(err)
	}
	err = initializers.PublishMessageWithRouteKey("quote.default", "quote.sub.build", "text/plain", &b, amqp.Table{}, amqp.Persistent)
	if err != nil {
		log.Println(err)
	}
}
