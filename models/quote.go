package models

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/shopspring/decimal"

	"quote/utils"
)

const RedisNotify = "quote"

type Quote struct {
	CommonModel
	Type      string          `json:"-"`
	BaseId    int             `json:"-"`
	QuoteId   int             `json:"-"`
	MarketId  int             `json:"-"`
	Timestamp int64           `json:"timestamp"`
	Price     decimal.Decimal `json:"price" gorm:"type:decimal(32,16);default:0;"`
	Source    string          `json:"source"`
	Base      string          `sql:"-" json:"base"`
	Quote     string          `sql:"-" json:"quote"`
}

func (quote *Quote) AfterFind() {
	quote.Source = strings.ToLower(strings.Split(quote.Type, "Quotes::")[1])
	for _, currency := range AllCurrencies {
		if currency.Id == quote.BaseId {
			quote.Base = currency.Symbol
		}
		if currency.Id == quote.QuoteId {
			quote.Quote = currency.Symbol
		}
	}
}

func (quote *Quote) AfterUpdate() {
	quote.notifyQuote()
}

func (quote *Quote) notifyQuote() {
	b, err := json.Marshal(*quote)
	if err != nil {
		log.Println("error:", err)
	}
	utils.PublishToPubSubChannels(RedisNotify, &b)
}
