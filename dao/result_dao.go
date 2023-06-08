package dao

import (
	"Lucky_Wheel/models"
	"github.com/go-xorm/xorm"
)

type ResultDao struct {
	engine *xorm.Engine
}

func NewResultDao(engine *xorm.Engine) *ResultDao {
	return &ResultDao{
		engine: engine,
	}
}

func (d *ResultDao) Get(id int) *models.LwResult {
	data := &models.LwResult{Id: id}
	ok, err := d.engine.Get(data)
	if ok && err == nil {
		return data
	} else {
		data.Id = 0
		return data
	}
}

func (d *ResultDao) GetAll(page, size int) []models.LwResult {
	offset := (page - 1) * size
	dataList := make([]models.LwResult, 0)
	err := d.engine.
		Desc("id").
		Limit(size, offset).
		Find(&dataList)
	if err != nil {
		return dataList
	}
	return dataList
}

func (d *ResultDao) CountAll() int64 {
	count, err := d.engine.Count(&models.LwResult{})
	if err != nil {
		return 0
	} else {
		return count
	}
}

func (d *ResultDao) GetNewPrize(size int, giftIds []int) []models.LwResult {
	datalist := make([]models.LwResult, 0)
	err := d.engine.
		In("gift_id", giftIds).
		Desc("id").
		Limit(size).
		Find(&datalist)
	if err != nil {
		return datalist
	} else {
		return datalist
	}
}

func (d *ResultDao) SearchByGift(giftId, page, size int) []models.LwResult {
	offset := (page - 1) * size
	datalist := make([]models.LwResult, 0)
	err := d.engine.
		Where("gift_id=?", giftId).
		Desc("id").
		Limit(size, offset).
		Find(&datalist)
	if err != nil {
		return datalist
	} else {
		return datalist
	}
}

func (d *ResultDao) SearchByUser(uid, page, size int) []models.LwResult {
	offset := (page - 1) * size
	datalist := make([]models.LwResult, 0)
	err := d.engine.
		Where("uid=?", uid).
		Desc("id").
		Limit(size, offset).
		Find(&datalist)
	if err != nil {
		return datalist
	} else {
		return datalist
	}
}

func (d *ResultDao) CountByGift(giftId int) int64 {
	num, err := d.engine.
		Where("gift_id=?", giftId).
		Count(&models.LwResult{})
	if err != nil {
		return 0
	} else {
		return num
	}
}

func (d *ResultDao) CountByUser(uid int) int64 {
	num, err := d.engine.
		Where("uid=?", uid).
		Count(&models.LwResult{})
	if err != nil {
		return 0
	} else {
		return num
	}
}

func (d *ResultDao) Delete(id int) error {
	_, err := d.engine.Delete(&models.LwResult{Id: id})
	return err
}

func (d *ResultDao) Update(data *models.LwResult, columns []string) error {
	_, err := d.engine.Id(data.Id).MustCols(columns...).Update(data)
	return err
}

func (d *ResultDao) Create(data *models.LwResult) error {
	_, err := d.engine.Insert(data)
	return err
}
