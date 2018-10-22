package models

import (
	"strings"

	"github.com/shopspring/decimal"
)

type Quote struct {
	CommonModel
	Type         string          `json:"-"`
	CoinId       int             `json:"-"`
	Timestamp    int64           `json:"timestamp"`
	CurrencyName string          `json:"-"`
	Price        decimal.Decimal `json:"price"`

	Source   string `json:"source"`
	Currency string `json:"currency"`
}

func (quote *Quote) AfterFind() {
	quote.Source = strings.Split(quote.Type, "Quotes::")[0]
	quote.Currency = quote.CurrencyName
}
