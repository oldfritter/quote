package sneakerWorkers

// import (
//   "encoding/json"
//   "time"
//
//   "github.com/gomodule/redigo/redis"
//   "github.com/streadway/amqp"
//
//   "quote/initializers"
//   . "quote/models"
//   "quote/utils"
// )
//
// var payload struct {
//   MarketId int   `json:"market_id"`
//   Period   int64 `json:"period"`
// }
//
// func (worker Worker) QuoteBuildWorker(payloadJson *[]byte) (err error) {
//   worker.LogInfo(string(*payloadJson))
//   json.Unmarshal([]byte(*payloadJson), &payload)
//   buildQuote(&worker, payload.MarketId)
//   return
// }
//
// func buildQuote(worker *Worker, marketId int) {
//   dataRedis := utils.GetRedisConn("data")
//   defer dataRedis.Close()
//   market, err := FindMarketById(marketId)
//   if err != nil {
//     worker.LogError(err)
//     return
//   }
//
// }
