package dao

import (
	"Lucky_Wheel/models"
	"github.com/go-xorm/xorm"
)

type UserDao struct {
	engine *xorm.Engine
}

func NewUserDao(engine *xorm.Engine) *UserDao {
	return &UserDao{
		engine: engine,
	}
}

func (d *UserDao) Get(id int) *models.LwUser {
	data := &models.LwUser{Id: id}
	ok, err := d.engine.Get(data)
	if ok && err == nil {
		return data
	} else {
		data.Id = 0
		return data
	}
}

func (d *UserDao) GetAll(page, size int) []models.LwUser {
	offset := (page - 1) * size
	dataList := make([]models.LwUser, 0)
	err := d.engine.
		Desc("id").
		Limit(size, offset).
		Find(&dataList)
	if err != nil {
		return dataList
	}
	return dataList
}

func (d *UserDao) CountAll() int {
	count, err := d.engine.Count(&models.LwUser{})
	if err != nil {
		return 0
	} else {
		return int(count)
	}
}

func (d *UserDao) Delete(id int) error {
	_, err := d.engine.Delete(&models.LwUser{Id: id})
	return err
}

func (d *UserDao) Update(data *models.LwUser, columns []string) error {
	_, err := d.engine.Id(data.Id).MustCols(columns...).Update(data)
	return err
}

func (d *UserDao) Create(data *models.LwUser) error {
	_, err := d.engine.Insert(data)
	return err
}
