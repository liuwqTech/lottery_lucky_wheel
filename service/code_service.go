package service

import (
	"Lucky_Wheel/dao"
	"Lucky_Wheel/datasource"
	"Lucky_Wheel/models"
)

type CodeService interface {
	GetAll(page, size int) []models.LwCode
	CountByGift(giftId int) int64
	CountAll() int64
	Search(giftId int) []models.LwCode
	Get(id int) *models.LwCode
	Delete(id int) error
	Update(data *models.LwCode, columns []string) error
	Create(data *models.LwCode) error
	NextUsingCode(giftId, codeId int) *models.LwCode
}

type codeService struct {
	dao *dao.CodeDao
}

func NewCodeService() CodeService {
	return &codeService{
		dao: dao.NewCodeDao(datasource.InstanceDbMaster()),
	}
}

func (s *codeService) GetAll(page, size int) []models.LwCode {
	return s.dao.GetAll(page, size)
}

func (s *codeService) Search(giftId int) []models.LwCode {
	return s.dao.Search(giftId)
}
func (s *codeService) CountByGift(giftId int) int64 {
	return s.dao.CountByGift(giftId)
}

func (s *codeService) CountAll() int64 {
	return s.dao.CountAll()
}

func (s *codeService) Get(id int) *models.LwCode {
	return s.dao.Get(id)
}

func (s *codeService) Delete(id int) error {
	return s.dao.Delete(id)
}

func (s *codeService) Update(data *models.LwCode, columns []string) error {
	return s.dao.Update(data, columns)
}

func (s *codeService) Create(data *models.LwCode) error {
	return s.dao.Create(data)
}

func (s *codeService) NextUsingCode(giftId, codeId int) *models.LwCode {
	return s.dao.NextUsingCode(giftId, codeId)
}
