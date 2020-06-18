package quote

import (
	"encoding/json"
	"fmt"
	// "log"

	"github.com/gomodule/redigo/redis"

	. "quote/models"
	"quote/utils"
)

func SaveDataFromRedis() {
	m := utils.DbBegin()
	defer m.DbRollback()
	var markets []Market
	m.Where("visible = ?", true).Find(&markets)

	dataRedis := utils.GetRedisConn("data")
	defer dataRedis.Close()

	for _, market := range markets {
		keys, err := redis.Strings(dataRedis.Do("KEYS", fmt.Sprintf("Quotes:%v:*", market.Id)))
		if err != nil {
		}
		for _, key := range keys {
			qByte, err := redis.Bytes(dataRedis.Do("GET", key))
			if err != nil {
			}
			var quote Quote
			json.Unmarshal(qByte, &quote)
			m.Save(&quote)
			// log.Println(quote)
		}
	}
	m.DbCommit()
}
