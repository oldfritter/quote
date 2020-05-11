package binance

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	. "quote/models"
	"quote/utils"
)

type BinanceResponse struct {
	Symbols []struct {
		Symbol string `json:"symbol"`
		Status string `json:"status"`
	} `json:"symbols"`
}

func GetBinanceMarkets() {
	ctx, cancelFun := context.WithTimeout(context.Background(), time.Second*10)
	defer cancelFun()
	req, err := http.NewRequest(http.MethodGet, "https://api.binance.com/api/v1/exchangeInfo", nil)
	if err != nil {
		log.Println(err)
		return
	}
	resp, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil {
		log.Println(err)
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}
	var binanceResponse BinanceResponse
	err = json.Unmarshal(body, &binanceResponse)
	if err != nil {
		log.Println(err)
		return
	}
	db := utils.DbBegin()
	defer db.DbRollback()
	var markets []Market
	db.Model(Market{}).Where("source = ?", "binance").Updates(Market{Visible: false})
	for _, symbol := range binanceResponse.Symbols {
		if symbol.Status == "TRADING" {
			db.Where("source = ?", "binance").Find(&markets)
			var market Market
			db.FirstOrInit(&market, map[string]interface{}{"name": symbol.Symbol, "symbol": strings.ToLower(symbol.Symbol), "source": "binance"})
			market.Visible = true
			db.Save(&market)
		}
	}
	db.DbCommit()
	return
}
