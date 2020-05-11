package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"strconv"

	"github.com/robfig/cron/v3"

	envConfig "quote/config"
	"quote/initializers"
	"quote/schedules/backup/tasks"
	"quote/schedules/sourceData"
	"quote/sourceData/binance"
	"quote/utils"
	"quote/workers/sneakerWorkers"
)

func main() {
	initialize()

	setLog()
	err := ioutil.WriteFile("pids/schedule.pid", []byte(strconv.Itoa(os.Getpid())), 0644)
	if err != nil {
		fmt.Println(err)
	}

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	closeResource()
}

func initialize() {
	envConfig.InitEnv()
	utils.InitDB()
	utils.InitAwsS3Config()
	utils.InitRedisPools()
	initializers.InitializeAmqpConfig()
	initializers.LoadCacheData()
	InitSchedule()

	sneakerWorkers.InitWorkers()
}

func closeResource() {
	initializers.CloseAmqpConnection()
	utils.CloseRedisPools()
	utils.CloseDB()
}

func InitSchedule() {
	c := cron.New(cron.WithSeconds())

	// 日志备份
	c.AddFunc("0 55 23 * * *", tasks.BackupLogFiles)
	c.AddFunc("0 56 23 * * *", tasks.UploadLogFileToS3)
	c.AddFunc("0 59 23 * * *", tasks.CleanLogs)

	// 清理过期价格
	c.AddFunc("0 59 23 * * *", sourceData.CleanMarketExpiredKLine)
	c.AddFunc("0 59 23 * * *", sourceData.CleanMarketExpiredPrice)

	c.AddFunc("0 0 * * * *", binance.GetBinanceMarkets)

	c.Start()
}

func setLog() {
	err := os.Mkdir("logs", 0755)
	if err != nil {
		if !os.IsExist(err) {
			log.Fatalf("create folder error: %v", err)
		}
	}
	file, err := os.OpenFile("logs/schedule.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("open file error: %v", err)
	}
	log.SetOutput(file)
}
