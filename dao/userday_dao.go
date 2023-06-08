package dao

import (
	"Lucky_Wheel/models"
	"github.com/go-xorm/xorm"
)

type UserDayDao struct {
	engine *xorm.Engine
}

func NewUserDayDao(engine *xorm.Engine) *UserDayDao {
	return &UserDayDao{
		engine: engine,
	}
}

func (d *UserDayDao) Get(id int) *models.LwUserday {
	data := &models.LwUserday{Id: id}
	ok, err := d.engine.Get(data)
	if ok && err == nil {
		return data
	} else {
		data.Id = 0
		return data
	}
}

func (d *UserDayDao) GetAll(page, size int) []models.LwUserday {
	offset := (page - 1) * size
	dataList := make([]models.LwUserday, 0)
	err := d.engine.
		Desc("id").
		Limit(size, offset).
		Find(&dataList)
	if err != nil {
		return dataList
	}
	return dataList
}

func (d *UserDayDao) CountAll() int64 {
	count, err := d.engine.Count(&models.LwUserday{})
	if err != nil {
		return 0
	} else {
		return count
	}
}

func (d *UserDayDao) Search(uid, day int) []models.LwUserday {
	datalist := make([]models.LwUserday, 0)
	err := d.engine.
		Where("uid=?", uid).
		Where("day=?", day).
		Desc("id").
		Find(&datalist)
	if err != nil {
		return datalist
	} else {
		return datalist
	}
}

func (d *UserDayDao) Count(uid, day int) int {
	info := &models.LwUserday{}
	ok, err := d.engine.
		Where("uid=?", uid).
		Where("day=?", day).
		Get(info)
	if !ok || err != nil {
		return 0
	} else {
		return info.Num
	}
}

func (d *UserDayDao) Delete(id int) error {
	_, err := d.engine.Delete(&models.LwUserday{Id: id})
	return err
}

func (d *UserDayDao) Update(data *models.LwUserday, columns []string) error {
	_, err := d.engine.Id(data.Id).MustCols(columns...).Update(data)
	return err
}

func (d *UserDayDao) Create(data *models.LwUserday) error {
	_, err := d.engine.Insert(data)
	return err
}
