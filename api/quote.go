package api

import (
	"github.com/labstack/echo"
	"net/http"
	"strings"

	. "github.com/oldfritter/quote/models"
	"github.com/oldfritter/quote/utils"
)

type Result struct {
	Data []Coin `json:"data"`
}

func GetApiQuotes(context echo.Context) error {
	db := utils.DbBegin()
	defer db.DbRollback()
	var coins []Coin
	var symbols, currencies, sources []string
	if context.QueryParam("symbols") != "" {
		symbols = strings.Split(context.QueryParam("symbols"), ",")
	}
	if context.QueryParam("currency") != "" {
		currencies = strings.Split(context.QueryParam("currency"), ",")
	}
	for _, source := range strings.Split(context.QueryParam("source"), ",") {
		if source != "" {
			sources = append(sources, "Quotes::"+strings.Title(source))
		}
	}
	if db.Where("symbol in (?)", symbols).Find(&coins).RecordNotFound() {
		return utils.BuildError("1020")
	}
	for i, coin := range coins {
		conditions := db.Where("coin_id = ?", coin.Id)
		if len(sources) > 0 {
			conditions = conditions.Where("type in (?)", sources)
		}
		if len(currencies) > 0 {
			conditions = conditions.Where("currency_name in (?)", currencies)
		}
		conditions.Find(&coins[i].Quotes)
	}
	response := utils.ArrayResponse
	response.Body = &Result{Data: coins}
	return context.JSON(http.StatusOK, response)
}
