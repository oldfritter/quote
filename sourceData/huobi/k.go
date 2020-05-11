package huobi

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"quote/models"
	"quote/sourceData/redis"
	"quote/utils"
	"github.com/shopspring/decimal"
)

type HuobiPayload struct {
	Ch   string `json:"ch"`
	Ts   int64  `json:"ts"`
	Tick struct {
		Id     int             `json:"id"`
		Open   decimal.Decimal `json:"open"`
		Close  decimal.Decimal `json:"close"`
		Low    decimal.Decimal `json:"low"`
		High   decimal.Decimal `json:"high"`
		Amount decimal.Decimal `json:"amount"`
		Vol    decimal.Decimal `json:"vol"`
		Count  int             `json:"count"`
	} `json:"tick"`
}

func GetHuobiPrice() error {
	u := url.URL{Scheme: "wss", Host: "api.huobi.pro", Path: "/ws"}
	fmt.Println("connecting to ", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		fmt.Println("dial:", err)
		return err
	}
	huobiMarkets := models.FindMarketsBySource("huobi")
	for _, m := range huobiMarkets {
		err = c.WriteMessage(websocket.TextMessage, []byte("{\"sub\":\"market."+m.Symbol+".kline.1min\",\"id\":\"gorilla\"}"))
		if err != nil {
			fmt.Println("write:", err)
			return err
		}
	}
	defer c.Close()

	done := make(chan struct{})
	errChan := make(chan error)

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				fmt.Println("read:", err)
				errChan <- err
				return
			}
			strByte, _ := utils.GzipDecode(message)
			var payload HuobiPayload
			json.Unmarshal(strByte, &payload)
			if payload.Ch != "" {
				kLine := models.KLine{
					Timestamp: payload.Ts,
					Open:      payload.Tick.Open,
					Close:     payload.Tick.Close,
					Low:       payload.Tick.Low,
					High:      payload.Tick.High,
					Vol:       payload.Tick.Vol,
				}
				symbol := strings.Replace(payload.Ch, "market.", "", 1)
				symbol = strings.Replace(symbol, ".kline.1min", "", 1)
				for _, m := range huobiMarkets {
					if m.Symbol == symbol {
						redis.Save(&m, &kLine)
					}
				}
			}
		}
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return nil
		case <-errChan:
			return fmt.Errorf("connect is closed!")
		case t := <-ticker.C:
			err := c.WriteMessage(websocket.TextMessage, []byte(t.String()))
			if err != nil {
				fmt.Println("write:", err)
				return err
			}
		}
	}
	return nil
}
