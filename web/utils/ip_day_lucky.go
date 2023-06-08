package utils

import (
	comm "Lucky_Wheel/common"
	"Lucky_Wheel/datasource"
	"fmt"
	"log"
	"math"
)

const ipFrameSize = 2

//func init() {
//	resetGroupIpList()
//}
//
//func resetGroupIpList() {
//	log.Println("ip_day_luck resetGroupIpList start")
//	cacheObj := datasource.InstanceCache()
//	for i := 0; i < ipFrameSize; i++ {
//		key := fmt.Sprintf("day_ips_%d", i)
//		cacheObj.Do("DEL", key)
//	}
//	log.Println("ip_day_luck resetGroupIpList stop")
//	// ip 当天的统计数 0点时归0 设置定时器
//	duration := comm.NextDayDuration()
//	time.AfterFunc(duration, resetGroupIpList)
//}

func IncrIpLuckyNum(strIp string) int64 {
	ip := comm.Ip4toInt(strIp)
	i := ip % ipFrameSize
	key := fmt.Sprintf("day_ips_%d", i)
	cacheObj := datasource.InstanceCache()
	rs, err := cacheObj.Do("HINCRBY", key, ip, 1)
	if err != nil {
		log.Println("ip_day_lucky redis HINCRBY error", err)
		return math.MaxInt32
	}
	return rs.(int64)
}
