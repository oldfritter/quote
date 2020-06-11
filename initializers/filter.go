package initializers

import (
	"github.com/gomodule/redigo/redis"
	"github.com/labstack/echo"

	"quote/utils"
)

func Filter(next echo.HandlerFunc) echo.HandlerFunc {
	return func(context echo.Context) (err error) {
		limitRedis := utils.GetRedisConn("limit")
		key := context.RealIP()
		times, _ := redis.Int(limitRedis.Do("GET", key))
		if times == 0 {
			limitRedis.Do("Set", key, 1)
			limitRedis.Do("EXPIRE", key, 60)
		} else {
			if times > 100 {
				return utils.BuildError("1021")
			} else {
				limitRedis.Do("INCR", key)
			}
		}
		return
	}
}
