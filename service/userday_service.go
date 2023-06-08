package service

import (
	"Lucky_Wheel/dao"
	"Lucky_Wheel/datasource"
	"Lucky_Wheel/models"
)

type UserDayService interface {
	GetAll(page, size int) []models.LwUserday
	CountAll() int64
	Search(uid, day int) []models.LwUserday
	Count(uid, day int) int
	Get(id int) *models.LwUserday
	Delete(id int) error
	Update(data *models.LwUserday, columns []string) error
	Create(data *models.LwUserday) error
	GetUserToday(uid int) *models.LwUserday
}

type userDayService struct {
	dao *dao.UserDayDao
}

func NewUserDayService() UserDayService {
	return &userDayService{
		dao: dao.NewUserDayDao(datasource.InstanceDbMaster()),
	}
}

func (s *userDayService) GetAll(page, size int) []models.LwUserday {
	return s.dao.GetAll(page, size)
}

func (s *userDayService) CountAll() int64 {
	return s.dao.CountAll()
}

func (s *userDayService) Get(id int) *models.LwUserday {
	return s.dao.Get(id)
}

func (s *userDayService) Search(uid, day int) []models.LwUserday {
	return s.dao.Search(uid, day)
}

func (s *userDayService) Count(uid, day int) int {
	return s.dao.Count(uid, day)
}

func (s *userDayService) Delete(id int) error {
	return s.dao.Delete(id)
}

func (s *userDayService) Update(data *models.LwUserday, columns []string) error {
	return s.dao.Update(data, columns)
}

func (s *userDayService) Create(data *models.LwUserday) error {
	return s.dao.Create(data)
}

func (s *userDayService) GetUserToday(uid int) *models.LwUserday {
	return &models.LwUserday{}
	//y, m, d := time.Now().Date()
	//strDay := fmt.Sprintf("%d%02d%02d", y, m, d)
	//day, _ := strconv.Atoi(strDay)
	//list := s.dao.Search(uid, day)
	//if list != nil && len(list) > 0 {
	//	return &list[0]
	//} else {
	//	return nil
	//}
}
