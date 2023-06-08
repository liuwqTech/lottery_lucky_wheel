package dao

import (
	comm "Lucky_Wheel/common"
	"Lucky_Wheel/models"
	"github.com/go-xorm/xorm"
	"log"
)

type GiftDao struct {
	engine *xorm.Engine
}

func NewGiftDao(engine *xorm.Engine) *GiftDao {
	return &GiftDao{
		engine: engine,
	}
}

func (d *GiftDao) Get(id int) *models.LwGift {
	data := &models.LwGift{Id: id}
	ok, err := d.engine.Get(data)
	if ok && err == nil {
		return data
	} else {
		data.Id = 0
		return data
	}
}

func (d *GiftDao) GetAll() []models.LwGift {
	dataList := make([]models.LwGift, 0)
	err := d.engine.Asc("sys_status").Asc("displayorder").Find(&dataList)
	if err != nil {
		log.Println("gift_data.GetAll error=", err)
		return dataList
	}
	return dataList
}

func (d *GiftDao) CountAll() int64 {
	count, err := d.engine.Count(&models.LwGift{})
	if err != nil {
		return 0
	} else {
		return count
	}
}

func (d *GiftDao) Delete(id int) error {
	data := &models.LwGift{Id: id, SysStatus: 1}
	_, err := d.engine.Id(data.Id).Update(data)
	return err
}

func (d *GiftDao) Update(data *models.LwGift, columns []string) error {
	_, err := d.engine.Id(data.Id).MustCols(columns...).Update(data)
	return err
}

func (d *GiftDao) Create(data *models.LwGift) error {
	_, err := d.engine.Insert(data)
	return err
}

func (d *GiftDao) GetAllUse() []models.LwGift {
	now := comm.NowUnix()
	dataList := make([]models.LwGift, 0)
	err := d.engine.Cols("id", "title", "prize_num", "left_num", "prize_code", "prize_time", "img", "displayorder", "gtype", "gdata").
		Desc("gtype").
		Asc("displayorder").
		Where("prize_num>=?", 0).
		Where("sys_status=?", 0).
		Where("time_begin<=?", now).
		Where("time_end>=?", now).
		Find(&dataList)
	if err != nil {
		log.Println("gift_dao.GetAllUse error=", err)
	}
	return dataList
}

func (d *GiftDao) DecrLeftNum(id, num int) (int64, error) {
	r, err := d.engine.Id(id).
		Decr("left_num", num).
		Where("left_num >=?", num).
		Update(&models.LwGift{Id: id})
	return r, err
}

func (d *GiftDao) IncrLeftNum(id, num int) (int64, error) {
	r, err := d.engine.Id(id).
		Incr("left_num", num).
		Where("left_num >=?", num).
		Update(&models.LwGift{Id: id})
	return r, err
}
