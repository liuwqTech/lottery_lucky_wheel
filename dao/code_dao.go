package dao

import (
	"Lucky_Wheel/models"
	"github.com/go-xorm/xorm"
)

type CodeDao struct {
	engine *xorm.Engine
}

func NewCodeDao(engine *xorm.Engine) *CodeDao {
	return &CodeDao{
		engine: engine,
	}
}

func (d *CodeDao) Get(id int) *models.LwCode {
	data := &models.LwCode{Id: id}
	ok, err := d.engine.Get(data)
	if ok && err == nil {
		return data
	} else {
		data.Id = 0
		return data
	}
}

func (d *CodeDao) Search(giftId int) []models.LwCode {
	datalist := make([]models.LwCode, 0)
	err := d.engine.
		Desc("id").
		Where("gift_id = ?", giftId).
		Find(&datalist)
	if err != nil {
		return datalist
	}
	return datalist
}

func (d *CodeDao) CountByGift(giftId int) int64 {
	count, err := d.engine.
		Desc("id").
		Where("gift_id = ?", giftId).
		Count()
	if err != nil {
		return 0
	}
	return count
}

func (d *CodeDao) GetAll(page, size int) []models.LwCode {
	offset := (page - 1) * size
	dataList := make([]models.LwCode, 0)
	err := d.engine.
		Desc("id").
		Limit(size, offset).
		Find(&dataList)
	if err != nil {
		return dataList
	}
	return dataList
}

func (d *CodeDao) CountAll() int64 {
	count, err := d.engine.Count(&models.LwCode{})
	if err != nil {
		return 0
	} else {
		return count
	}
}

func (d *CodeDao) Delete(id int) error {
	data := &models.LwCode{Id: id, SysStatus: 1}
	_, err := d.engine.Id(data.Id).Update(data)
	return err
}

func (d *CodeDao) Update(data *models.LwCode, columns []string) error {
	_, err := d.engine.Id(data.Id).MustCols(columns...).Update(data)
	return err
}

func (d *CodeDao) Create(data *models.LwCode) error {
	_, err := d.engine.Insert(data)
	return err
}

func (d *CodeDao) NextUsingCode(giftId, codeId int) *models.LwCode {
	datalist := make([]models.LwCode, 0)
	err := d.engine.Where("gift_id=?", giftId).
		Where("sys_status=?", 0).
		Where("id>?", codeId).
		Asc("id").
		Limit(1).
		Find(&datalist)
	if err != nil || len(datalist) == 0 {
		return nil
	}
	return &datalist[0]
}
