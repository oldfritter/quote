package models

import (
	"strings"

	"github.com/shopspring/decimal"
)

type Quote struct {
	CommonModel
	Type       string          `json:"-"`
	CurrencyId int             `json:"-"`
	QuoteId    int             `json:"-"`
	MarketId   int             `json:"-"`
	Timestamp  int64           `json:"timestamp"`
	Price      decimal.Decimal `json:"price" gorm:"type:decimal(32,16);default:0;"`
	Source     string          `json:"source"`
	Currency   string          `json:"currency" sql:"-"`
}

func (quote *Quote) AfterFind() {
	quote.Source = strings.ToLower(strings.Split(quote.Type, "Quotes::")[1])
	// quote.Currency = quote.CurrencyName
}
