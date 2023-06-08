package cron

import (
	comm "Lucky_Wheel/common"
	"Lucky_Wheel/service"
	"Lucky_Wheel/web/utils"
	"log"
	"time"
)

func ConfigureAppOne() {
	go resetAllGiftPrizeData()
}

func resetAllGiftPrizeData() {
	giftService := service.NewGiftService()
	nowTime := comm.NowUnix()
	list := giftService.GetAll(false)
	for _, giftInfo := range list {
		if giftInfo.PrizeTime > 0 && (giftInfo.PrizeData == "" || giftInfo.PrizeEnd <= nowTime) {
			log.Println("crontab start utils.ResetGiftPrizeData giftInfo=", giftInfo)
			utils.ResetGiftPrizeData(&giftInfo, giftService)
			giftService.GetAll(true)
			log.Println("crontab end utils.ResetGiftPrizeData giftInfo=", giftInfo)
		}
	}
	// 每5min执行一次
	time.AfterFunc(5*time.Minute, resetAllGiftPrizeData)
}

func distributionAllGiftPool() {
	log.Println("crontab start utils.DistributionAllGiftPool")
	utils.DistributionGiftPool()
	log.Println("crontab end utils.DistributionAllGiftPool")
	// 每5min执行一次
	time.AfterFunc(time.Minute, distributionAllGiftPool)
}
