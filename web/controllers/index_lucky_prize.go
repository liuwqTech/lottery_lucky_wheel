package controllers

import (
	"Lucky_Wheel/conf"
	"Lucky_Wheel/models"
)

func (c *IndexController) prize(prizeCode int, limitBlack bool) *models.ObjGiftPrize {
	var prizeGift *models.ObjGiftPrize
	giftList := c.ServiceGift.GetAllUse(true)
	for _, gift := range giftList {
		if gift.PrizeCodeA <= prizeCode && gift.PrizeCodeB >= prizeCode {
			// 满足中奖条件
			if !limitBlack || gift.Gtype < conf.GtypeGiftSmall {
				prizeGift = &gift
				break
			}
		}
	}
	return prizeGift
}
