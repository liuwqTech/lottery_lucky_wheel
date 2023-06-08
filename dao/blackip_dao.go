package dao

import (
	"Lucky_Wheel/models"
	"github.com/go-xorm/xorm"
)

type BlackIpDao struct {
	engine *xorm.Engine
}

func NewBlackIpDao(engine *xorm.Engine) *BlackIpDao {
	return &BlackIpDao{
		engine: engine,
	}
}

func (d *BlackIpDao) Get(id int) *models.LwBlackip {
	data := &models.LwBlackip{Id: id}
	ok, err := d.engine.Get(data)
	if ok && err == nil {
		return data
	} else {
		data.Id = 0
		return data
	}
}

func (d *BlackIpDao) GetAll(page, size int) []models.LwBlackip {
	offset := (page - 1) * size
	dataList := make([]models.LwBlackip, 0)
	err := d.engine.
		Desc("id").
		Limit(size, offset).
		Find(&dataList)
	if err != nil {
		return dataList
	}
	return dataList
}

func (d *BlackIpDao) CountAll() int64 {
	count, err := d.engine.Count(&models.LwBlackip{})
	if err != nil {
		return 0
	} else {
		return count
	}
}

func (d *BlackIpDao) Delete(id int) error {
	_, err := d.engine.Delete(&models.LwBlackip{Id: id})
	return err
}

func (d *BlackIpDao) Update(data *models.LwBlackip, columns []string) error {
	_, err := d.engine.Id(data.Id).MustCols(columns...).Update(data)
	return err
}

func (d *BlackIpDao) Create(data *models.LwBlackip) error {
	_, err := d.engine.Insert(data)
	return err
}

func (d *BlackIpDao) GetByIp(ip string) *models.LwBlackip {
	dataList := make([]models.LwBlackip, 0)
	err := d.engine.
		Where("ip=?", ip).
		Desc("id").
		Limit(1).
		Find(&dataList)
	if err != nil || len(dataList) < 1 {
		return nil
	} else {
		return &dataList[0]
	}
}

func (d *BlackIpDao) Search(ip string) []models.LwBlackip {
	datalist := make([]models.LwBlackip, 0)
	err := d.engine.
		Where("ip=?", ip).
		Desc("id").
		Find(&datalist)
	if err != nil {
		return datalist
	} else {
		return datalist
	}
}
