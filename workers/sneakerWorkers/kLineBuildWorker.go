package sneakerWorkers

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/streadway/amqp"

	"quote/initializers"
	. "quote/models"
	"quote/utils"
)

func (worker *Worker) KLineBuildWorker(payloadJson *[]byte) (err error) {
	fmt.Println("--KLineBuild payload: ", string(*payloadJson))
	worker.LogInfo("--KLineBuild payload: ", string(*payloadJson))
	// start := time.Now().UnixNano()
	var payload struct {
		MarketId int   `json:"market_id"`
		Period   int64 `json:"period"`
	}
	json.Unmarshal([]byte(*payloadJson), &payload)
	buildKLine(worker, payload.MarketId, payload.Period)
	worker.LogInfo("--KLineBuild payload: ", payload)
	return
}

func buildKLine(worker *Worker, marketId int, period int64) {
	dataRedis := utils.GetRedisConn("data")
	defer dataRedis.Close()
	market, err := FindMarketById(marketId)
	if err != nil {
		worker.LogError(err)
		return
	}
	now := time.Now().Unix()
	start := now / (period * 60) * (period * 60)
	kJsons, _ := redis.Values(dataRedis.Do("ZRANGEBYSCORE", market.KLine(1), start, now))
	var k, total KLine
	for i, kJson := range kJsons {
		json.Unmarshal(kJson.([]byte), &k)
		if i == 0 {
			total.Open = k.Open
		}
		total.Close = k.Close
		if total.Low.IsZero() || total.Low.GreaterThan(k.Low) {
			total.Low = k.Low
		}
		if total.High.LessThan(k.High) {
			total.High = k.High
		}
		total.Vol = total.Vol.Add(k.Vol)
	}
	total.Timestamp = start
	b, err := json.Marshal(total)
	if err != nil {
		worker.LogError(err)
	}
	dataRedis.Send("zremrangebyscore", market.KLine(period), start, now)
	dataRedis.Send("zadd", market.KLine(period), start, b)
	dataRedis.Send("DEL", market.UpdateLock(period))
	dataRedis.Do("")

	b, err = json.Marshal(NotifyKLine{
		KLine:    total,
		MarketId: marketId,
		Period:   period,
	})
	if err != nil {
		worker.LogError(err)
	}
	err = initializers.PublishMessageWithRouteKey(initializers.AmqpGlobalConfig.Exchange["fanout"]["k"], "#", "text/plain", &b, amqp.Table{}, amqp.Persistent)
	if err != nil {
		worker.LogError(err)
	}
}
