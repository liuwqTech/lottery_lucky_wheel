package service

import (
	"Lucky_Wheel/dao"
	"Lucky_Wheel/datasource"
	"Lucky_Wheel/models"
)

type ResultService interface {
	GetAll(page, size int) []models.LwResult
	CountAll() int64
	GetNewPrize(size int, giftIds []int) []models.LwResult
	SearchByGift(giftId, page, size int) []models.LwResult
	SearchByUser(uid, page, size int) []models.LwResult
	CountByGift(giftId int) int64
	CountByUser(uid int) int64
	Get(id int) *models.LwResult
	Delete(id int) error
	Update(data *models.LwResult, columns []string) error
	Create(data *models.LwResult) error
}

type resultService struct {
	dao *dao.ResultDao
}

func NewResultService() ResultService {
	return &resultService{
		dao: dao.NewResultDao(datasource.InstanceDbMaster()),
	}
}

func (s *resultService) GetAll(page, size int) []models.LwResult {
	return s.dao.GetAll(page, size)
}

func (s *resultService) CountAll() int64 {
	return s.dao.CountAll()
}

func (s *resultService) GetNewPrize(size int, giftIds []int) []models.LwResult {
	return s.dao.GetNewPrize(size, giftIds)
}

func (s *resultService) SearchByGift(giftId, page, size int) []models.LwResult {
	return s.dao.SearchByGift(giftId, page, size)
}

func (s *resultService) SearchByUser(uid, page, size int) []models.LwResult {
	return s.dao.SearchByUser(uid, page, size)
}

func (s *resultService) CountByGift(giftId int) int64 {
	return s.dao.CountByGift(giftId)
}

func (s *resultService) CountByUser(uid int) int64 {
	return s.dao.CountByUser(uid)
}

func (s *resultService) Get(id int) *models.LwResult {
	return s.dao.Get(id)
}

func (s *resultService) Delete(id int) error {
	return s.dao.Delete(id)
}

func (s *resultService) Update(data *models.LwResult, columns []string) error {
	return s.dao.Update(data, columns)
}

func (s *resultService) Create(data *models.LwResult) error {
	return s.dao.Create(data)
}
