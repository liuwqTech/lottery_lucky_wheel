package utils

import (
	comm "Lucky_Wheel/common"
	"Lucky_Wheel/conf"
	"Lucky_Wheel/datasource"
	"Lucky_Wheel/models"
	"Lucky_Wheel/service"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gomodule/redigo/redis"
)

func PrizeGift(id int, leftNum int) bool {
	giftService := service.NewGiftService()
	rows, err := giftService.DecrLeftNum(id, 1)
	if rows < 1 || err != nil {
		log.Println("prizedata.PrizeGift giftService DecrLeftNum error=", err, ", rows=", rows)
		return false
	}
	return true
}

func PrizeCodeDiff(id int, codeService service.CodeService) string {
	return prizeServCodeDiff(id, codeService)
}

func PrizeLocalCodeDiff(id int, codeService service.CodeService) string {
	lockUid := 0 - id - 1000000000
	LockLucky(lockUid)
	defer UnLockLucky(lockUid)
	codeId := 0
	codeInfo := codeService.NextUsingCode(id, codeId)
	if codeInfo != nil && codeInfo.Id > 0 {
		codeInfo.SysStatus = 2
		codeInfo.SysUpdated = comm.NowUnix()
		err := codeService.Update(codeInfo, nil)
		if err != nil {
			log.Println("prizedata.prizeCodeDiff codeService.Update error=", err)
			return ""
		}
	} else {
		log.Println("prizedata prizeCodeDiff num codeInfo, gift_id=", id)
		return ""
	}
	return codeInfo.Code
}

func ImportCacheCodes(id int, code string) bool {
	key := fmt.Sprintf("gift_code_%d", id)
	cacheObj := datasource.InstanceCache()
	_, err := cacheObj.Do("SADD", key, code)
	if err != nil {
		log.Println("prizedata.ImportCacheCodes SADD error=", err)
		return false
	} else {
		return true
	}
}

func RecacheCodes(id int, codeService service.CodeService) (sucNum, errNum int) {
	list := codeService.Search(id)
	if list == nil || len(list) <= 0 {
		return 0, 0
	}
	key := fmt.Sprintf("gift_code_%d", id)
	cacheObj := datasource.InstanceCache()
	tmpKey := "tmp_" + key
	for _, data := range list {
		if data.SysStatus == 0 {
			code := data.Code
			_, err := cacheObj.Do("SADD", tmpKey, code)
			if err != nil {
				log.Println("prizedata.RecacheCodes SADD error=", err)
				errNum++
			} else {
				sucNum++
			}
		}
	}
	_, err := cacheObj.Do("RENAME", tmpKey, key)
	if err != nil {
		log.Println("prizedata.RecacheCodes RENAME error=", err)
	}
	return sucNum, errNum
}

func GetCacheCodeNum(id int, codeService service.CodeService) (int, int) {
	num := 0
	cacheNum := 0
	list := codeService.Search(id)
	if len(list) > 0 {
		for _, data := range list {
			if data.SysStatus == 0 {
				num++
			}
		}
	}
	// redis中缓存的key
	key := fmt.Sprintf("gift_code_%d", id)
	cacheObj := datasource.InstanceCache()
	rs, err := cacheObj.Do("SCARD", key)
	if err != nil {
		log.Println("prizedata.GetCacheCodeNum SCARD error=", err)
	} else {
		cacheNum = int(comm.GetInt64(rs, 0))
	}
	return num, cacheNum
}

func prizeServCodeDiff(id int, codeService service.CodeService) string {
	key := fmt.Sprintf("gift_code_%d", id)
	cacheObj := datasource.InstanceCache()
	rs, err := cacheObj.Do("SPOP", key)
	if err != nil {
		log.Println("prizedata.prizeServCodeDiff SPOP error=", err)
		return ""
	}
	code := comm.GetString(rs, "")
	if code == "" {
		log.Println("prizedata.prizeServCodeDiff rs=", rs)
		return ""
	}
	err = codeService.Update(&models.LwCode{
		Code:       code,
		SysStatus:  2,
		SysUpdated: comm.NowUnix(),
	}, nil)
	if err != nil {
		return ""
	}
	return code
}

func ResetGiftPrizeData(giftInfo *models.LwGift, giftService service.GiftService) {
	if giftInfo == nil || giftInfo.Id < 1 {
		return
	}
	id := giftInfo.Id
	nowTime := comm.NowUnix()
	if giftInfo.SysStatus == 1 || giftInfo.TimeBegin >= nowTime || giftInfo.TimeEnd <= nowTime || giftInfo.PrizeNum <= 0 || giftInfo.LeftNum <= 0 {
		if giftInfo.PrizeData != "" {
			//todo:清空旧的发奖计划数据
			clearGiftPrizeData(giftInfo, giftService)
		}
		return
	}
	// 没有设置发奖周期
	dayNum := giftInfo.PrizeTime
	if dayNum <= 0 {
		setGiftPool(id, giftInfo.LeftNum)
		return
	}
	// 重置发奖计划数据
	setGiftPool(id, 0)
	// 实际的奖品计划分布运算
	prizeNum := giftInfo.PrizeNum
	avgNum := prizeNum / dayNum
	// 每天可以分配的奖品数
	dayPrizeNum := make(map[int]int)
	if avgNum >= 1 {
		for day := 0; day < dayNum; day++ {
			dayPrizeNum[day] = avgNum
		}
	}
	// 剩下的随机分配到任意哪天
	prizeNum -= dayNum * avgNum
	for prizeNum > 0 {
		prizeNum--
		day := comm.Random(dayNum)
		_, ok := dayPrizeNum[day]
		if !ok {
			dayPrizeNum[day] = 1
		} else {
			dayPrizeNum[day] += 1
		}
	}
	// 每天的map，每小时的map，60min数组，奖品数
	prizeData := make(map[int]map[int][60]int)
	for day, num := range dayPrizeNum {
		// 计算出来这一天的发奖计划
		dayPrizeData := getGiftPrizeDataOneDay(num)
		prizeData[day] = dayPrizeData
	}
	//将周期内每天、每小时、每分钟的数据 prizeData 格式化[时间：数量]
	dataList := formatGiftPrizeData(nowTime, dayNum, prizeData)
	str, err := json.Marshal(dataList)
	if err != nil {
		log.Println("prizedata.ResetGiftPrizeData json Marshal error=", err)
	} else {
		info := &models.LwGift{
			Id:         giftInfo.Id,
			LeftNum:    giftInfo.LeftNum,
			PrizeData:  string(str),
			PrizeBegin: nowTime,
			PrizeEnd:   nowTime + dayNum*86400,
			SysUpdated: nowTime,
		}
		err := giftService.Update(info, nil)
		if err != nil {
			log.Println("prizedata.ResetGiftPrizeData giftService.Update error=", err)
		}
	}
}

// 清空旧的发奖计划数据
func clearGiftPrizeData(giftInfo *models.LwGift, giftService service.GiftService) {
	info := &models.LwGift{
		Id:        giftInfo.Id,
		PrizeData: "",
	}
	err := giftService.Update(info, []string{"prize_data"})
	if err != nil {
		log.Println("prizedata.clearGiftPrizeData giftService.Update error=", err)
	}
	setGiftPool(giftInfo.Id, 0)
}

// 设置奖品池的库存数量
func setGiftPool(id int, num int) {
	key := "gift_pool"
	cacheObj := datasource.InstanceCache()
	_, err := cacheObj.Do("HSET", key, id, num)
	if err != nil {
		log.Println("prizedata.setGiftPool error=", err)
	}
}

// 计算出来一天的发奖计划
func getGiftPrizeDataOneDay(num int) map[int][60]int {
	rs := make(map[int][60]int)
	// 计算24h各自的奖品数
	hourData := [24]int{}
	if num > 100 {
		hourNum := 0
		for _, h := range conf.PrizeDataRandomDayTime {
			hourData[h]++
		}
		for h := 0; h < 24; h++ {
			d := hourData[h]
			n := num * d / 100
			hourData[h] = n
			hourNum += n
		}
		num -= hourNum
	}
	for num > 0 {
		num--
		hourIndex := comm.Random(100)
		h := conf.PrizeDataRandomDayTime[hourIndex]
		hourData[h]++
	}
	// 将每个小时内的奖品数量分配到60min
	for h, hnum := range hourData {
		if hnum <= 0 {
			continue
		}
		minuteData := [60]int{}
		if hnum >= 60 {
			avgMinute := hnum / 60
			for i := 0; i < 60; i++ {
				minuteData[i] = avgMinute
			}
			hnum -= avgMinute * 60
		}
		for hnum > 0 {
			hnum--
			m := comm.Random(60)
			minuteData[m]++
		}
		rs[h] = minuteData
	}
	return rs
}

// 将每天、每小时、每分钟的奖品数量，格式化为一个具体时间、数量的格式
// 结构为：prizeData [day][hour][minute]num
// result: [][时间，数量]
func formatGiftPrizeData(nowTime, dayNum int, prizeData map[int]map[int][60]int) [][2]int {
	rs := make([][2]int, 0)
	nowHour := time.Now().Hour()
	// 处理日期数据
	for dn := 0; dn < dayNum; dn++ {
		dayData, ok := prizeData[dn]
		if !ok {
			continue
		}
		dayTime := nowTime + dn*86400
		// 处理小时的数据
		for hn := 0; hn < 24; hn++ {
			hourData, ok := dayData[(hn+nowHour)%24]
			if !ok {
				continue
			}
			hourTime := dayTime + hn*3600
			//处理分钟的数据
			for mn := 0; mn < 60; mn++ {
				num := hourData[mn]
				if num <= 0 {
					continue
				}
				minuteTime := hourTime + mn*60
				rs = append(rs, [2]int{minuteTime, num})
			}
		}
	}
	return rs
}

func DistributionGiftPool() int {
	totalNum := 0
	now := comm.NowUnix()
	giftService := service.NewGiftService()
	list := giftService.GetAll(false)
	if list != nil && len(list) > 0 {
		for _, gift := range list {
			if gift.SysStatus != 0 {
				continue
			}
			if gift.PrizeNum < 1 {
				continue
			}
			if gift.TimeBegin > now || gift.TimeEnd < now {
				continue
			}
			if len(gift.PrizeData) <= 7 {
				continue
			}
			var cronData [][2]int
			err := json.Unmarshal([]byte(gift.PrizeData), &cronData)
			if err != nil {
				log.Println("prizedata.DistributionGiftPool Unmarshal error=", err)
			} else {
				index := 0
				giftNum := 0
				for i, data := range cronData {
					ct := data[0]
					num := data[1]
					if ct <= now {
						giftNum += num
						index = i + 1
					} else {
						break
					}
				}
				// 更新奖品池
				if giftNum > 0 {
					incrGiftPool(gift.Id, giftNum)
					totalNum += giftNum
				}
				// 更新奖品的发奖计划
				if index > 0 {
					if index >= len(cronData) {
						cronData = make([][2]int, 0)
					} else {
						cronData = cronData[index:]
					}
					str, err := json.Marshal(cronData)
					if err != nil {
						log.Println("prizedata.DistributionGiftPool Marshal", cronData, "，error=", err)
					}
					columns := []string{"prize_data"}
					err = giftService.Update(&models.LwGift{
						Id:        gift.Id,
						PrizeData: string(str),
					}, columns)
					if err != nil {
						log.Println("prizedata.DistributionGiftPool giftService.Update error=", err)
					}
				}
			}
		}
		if totalNum > 0 {
			giftService.GetAll(true)
		}
	}
	return totalNum
}

func incrGiftPool(id, num int) int {
	key := "gift_pool"
	cacheObj := datasource.InstanceCache()
	rtNum, err := redis.Int64(cacheObj.Do("HINCBY", key, id, num))
	if err != nil {
		log.Println("prizedata.incrGiftPool HINCBY error=", err)
		return 0
	}
	if int(rtNum) < num {
		// 递增少于预期值，补偿一次
		num2 := num - int(rtNum)
		rtNum, err = redis.Int64(cacheObj.Do("HINCBY", key, id, num2))
		if err != nil {
			log.Println("prizedata.incrGiftPool HINCBY2 error=", err)
			return 0
		}
	}
	return int(rtNum)
}
