package binance

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/shopspring/decimal"
	"quote/models"
	"quote/sourceData/redis"
)

type BinancePayload struct {
	Stream string `json:"stream"`
	Data   struct {
		E int64  `json:"E"`
		S string `json:"s"`
		K struct {
			O decimal.Decimal `json:"o"`
			C decimal.Decimal `json:"c"`
			H decimal.Decimal `json:"h"`
			L decimal.Decimal `json:"l"`
			V decimal.Decimal `json:"v"`
		} `json:"k"`
	} `json:"data"`
}

func GetBinancePrice() error {

	u := url.URL{Scheme: "wss", Host: "stream.binance.com:9443", Path: "/stream"}
	log.Println("connecting to ", u.String())
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Println("dial:", err)
		return err
	}
	stream := struct {
		Method string   `json:"method"`
		Params []string `json:"params"`
		Id     int64    `json:"id"`
	}{
		Method: "SUBSCRIBE",
		Id:     704483774,
	}
	binanceMarkets := models.FindMarketsBySource("binance")
	for _, m := range binanceMarkets {
		stream.Params = append(stream.Params, m.Symbol+"@kline_1m")
	}
	b, err := json.Marshal(stream)
	if err != nil {
		log.Println("error:", err)
	}
	err = c.WriteMessage(websocket.TextMessage, b)
	if err != nil {
		log.Println("write:", err)
		return err
	}
	defer c.Close()

	done := make(chan struct{})
	errChan := make(chan error)

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				errChan <- err
				return
			}
			var payload BinancePayload
			json.Unmarshal(message, &payload)
			kLine := models.KLine{
				Timestamp: payload.Data.E,
				Open:      payload.Data.K.O,
				Close:     payload.Data.K.C,
				Low:       payload.Data.K.L,
				High:      payload.Data.K.H,
				Vol:       payload.Data.K.V,
			}
			symbol := strings.ToLower(payload.Data.S)
			for _, m := range binanceMarkets {
				if m.Symbol == symbol {
					redis.Save(&m, &kLine)
				}
			}

		}
	}()

	for {
		select {
		case <-done:
			return nil
		case <-errChan:
			return fmt.Errorf("connect is closed!")
		}
	}
	return nil
}
