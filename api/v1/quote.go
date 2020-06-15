package v1

import (
	"github.com/labstack/echo"
	"net/http"
	"strings"

	. "quote/models"
	"quote/utils"
)

type Result struct {
	Data []Currency `json:"data"`
}

func GetApiQuotes(context echo.Context) error {
	db := utils.DbBegin()
	defer db.DbRollback()
	var coins, result []Currency
	var symbols, currencies, sources []string
	if context.QueryParam("symbols") != "" {
		symbols = strings.Split(context.QueryParam("symbols"), ",")
	}
	if context.QueryParam("currencies") != "" {
		currencies = strings.Split(context.QueryParam("currencies"), ",")
	}
	if context.QueryParam("currency") != "" {
		currencies = strings.Split(context.QueryParam("currency"), ",")
	}
	for _, source := range strings.Split(context.QueryParam("sources"), ",") {
		if source != "" {
			sources = append(sources, "Quotes::"+strings.Title(source))
		}
	}
	// for _, source := range strings.Split(context.QueryParam("source"), ",") {
	//   if source != "" {
	//     sources = append(sources, "Quotes::"+strings.Title(source))
	//   }
	// }
	if db.Where("symbol in (?)", symbols).Find(&coins).RecordNotFound() {
		return utils.BuildError("1020")
	}
	for i, coin := range coins {
		conditions := db.Where("base_id = ?", coin.Id)
		if len(sources) > 0 {
			conditions = conditions.Where("type in (?)", sources)
		}
		if len(currencies) > 0 {
			conditions = conditions.Joins("INNER JOIN (currencies) ON (currencies.id = quotes.quote_id)").Where("currencies.symbol in (?)", currencies)
		}
		conditions.Find(&coins[i].Quotes)
		if len(coins[i].Quotes) > 0 {
			result = append(result, coins[i])
		}
	}
	response := utils.ArrayResponse
	response.Body = &Result{Data: result}
	return context.JSON(http.StatusOK, response)
}
