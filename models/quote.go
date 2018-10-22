package models

import (
	"github.com/shopspring/decimal"
)

type Quote struct {
	CommonModel
	Type         string          `json:"type"`
	CoinId       int             `json:"-"`
	Timestamp    int64           `json:"timestamp"`
	CurrencyName string          `json:"currency_name"`
	Price        decimal.Decimal `json:"price"`
}
