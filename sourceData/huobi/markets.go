package huobi

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

type HuobiResponse struct {
	Data []struct {
		Symbol string `json:"symbol"`
		Status bool   `json:"visit-enabled"`
	} `json:"data"`
}

func GetMarkets() {
	ctx, cancelFun := context.WithTimeout(context.Background(), time.Second*10)
	defer cancelFun()
	req, err := http.NewRequest(http.MethodGet, "https://api.huobi.com/v1/settings/symbols?r=o58hsjzm4m", nil)
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
	var huobiResponse HuobiResponse
	err = json.Unmarshal(body, &huobiResponse)
	if err != nil {
		log.Println(err)
		return
	}
	db := utils.DbBegin()
	defer db.DbRollback()
	var markets []Market
	db.Model(Market{}).Where("source = ?", "huobi").Updates(Market{Visible: false})
	for _, symbol := range huobiResponse.Data {
		if symbol.Status {
			db.Where("source = ?", "huobi").Find(&markets)
			var market Market
			db.FirstOrInit(&market, map[string]interface{}{"symbol": symbol.Symbol, "name": strings.ToUpper(symbol.Symbol), "source": "huobi"})
			market.Visible = true
			db.Save(&market)
		}
	}
	db.DbCommit()
	return
}
