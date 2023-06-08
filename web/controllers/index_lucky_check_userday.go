package controllers

import (
	"Lucky_Wheel/conf"
	"Lucky_Wheel/models"
	"fmt"
	"log"
	"strconv"
	"time"
)

func (c *IndexController) checkUserDay(uid int) bool {
	userDayInfo := c.ServiceUserDay.GetUserToday(uid)
	log.Println(userDayInfo)
	if userDayInfo != nil && userDayInfo.Uid == uid {
		if userDayInfo.Num >= conf.UserPrizeMax {
			return false
		} else {
			userDayInfo.Num++
			err := c.ServiceUserDay.Update(userDayInfo, nil)
			if err != nil {
				log.Println("index_luck_check_user_day ServiceUserDay Update error", err)
			}
		}
	} else {
		y, m, d := time.Now().Date()
		strDay := fmt.Sprintf("%d%02d%02d", y, m, d)
		day, _ := strconv.Atoi(strDay)
		userDayInfo = &models.LwUserday{
			Uid:        uid,
			Day:        day,
			Num:        1,
			SysCreated: int(time.Now().Unix()),
		}
		err := c.ServiceUserDay.Create(userDayInfo)
		if err != nil {
			log.Println("index_luck_check_user_day ServiceUserDay Update error", err)
		}
	}
	return true
}
