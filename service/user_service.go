package service

import (
	comm "Lucky_Wheel/common"
	"Lucky_Wheel/dao"
	"Lucky_Wheel/datasource"
	"Lucky_Wheel/models"
	"fmt"
	"log"
	"sync"

	"github.com/gomodule/redigo/redis"
)

// 用户信息，可以缓存(本地或者redis)，有更新的时候，可以直接清除缓存或者根据具体情况更新缓存
var cachedUserList = make(map[int]*models.LwUser)
var cachedUserLock = sync.Mutex{}

type UserService interface {
	GetAll(page, size int) []models.LwUser
	CountAll() int
	Get(id int) *models.LwUser
	Delete(id int) error
	Update(data *models.LwUser, columns []string) error
	Create(data *models.LwUser) error
}

type userService struct {
	dao *dao.UserDao
}

func NewUserService() UserService {
	return &userService{
		dao: dao.NewUserDao(datasource.InstanceDbMaster()),
	}
}

func (s *userService) GetAll(page, size int) []models.LwUser {
	return s.dao.GetAll(page, size)
}

func (s *userService) CountAll() int {
	return s.dao.CountAll()
}

func (s *userService) Get(id int) *models.LwUser {
	data := s.getByCache(id)
	if data == nil || data.Id <= 0 {
		data = s.dao.Get(id)
		if data == nil || data.Id <= 0 {
			data = &models.LwUser{Id: id}
		}
		s.setByCache(data)
	}
	return data
}

func (s *userService) Delete(id int) error {
	return s.dao.Delete(id)
}

func (s *userService) Update(data *models.LwUser, columns []string) error {
	// 先更新缓存
	s.updateByCache(data, columns)
	// 然后再更新数据
	return s.dao.Update(data, columns)
}

func (s *userService) Create(data *models.LwUser) error {
	err := s.dao.Create(data)
	if err == nil {
		s.updateByCache(data, nil)
	}
	return err
}

// 从缓存中得到信息
func (s *userService) getByCache(id int) *models.LwUser {
	// 集群模式，redis缓存
	key := fmt.Sprintf("info_user_%d", id)
	rds := datasource.InstanceCache()
	dataMap, err := redis.StringMap(rds.Do("HGETALL", key))
	if err != nil {
		log.Println("user_service.getByCache HGETALL key=", key, ", error=", err)
		return nil
	}
	dataId := comm.GetInt64FromStringMap(dataMap, "Id", 0)
	if dataId <= 0 {
		return nil
	}
	data := &models.LwUser{
		Id:         int(dataId),
		Username:   comm.GetStringFromStringMap(dataMap, "Username", ""),
		Blacktime:  int(comm.GetInt64FromStringMap(dataMap, "Blacktime", 0)),
		Realname:   comm.GetStringFromStringMap(dataMap, "Realname", ""),
		Mobile:     comm.GetStringFromStringMap(dataMap, "Mobile", ""),
		Address:    comm.GetStringFromStringMap(dataMap, "Address", ""),
		SysCreated: int(comm.GetInt64FromStringMap(dataMap, "SysCreated", 0)),
		SysUpdated: int(comm.GetInt64FromStringMap(dataMap, "SysUpdated", 0)),
		SysIp:      comm.GetStringFromStringMap(dataMap, "SysIp", ""),
	}
	return data
}

// 将信息更新到缓存
func (s *userService) setByCache(data *models.LwUser) {
	if data == nil || data.Id <= 0 {
		return
	}
	id := data.Id
	// 集群模式，redis缓存
	key := fmt.Sprintf("info_user_%d", id)
	rds := datasource.InstanceCache()
	// 数据更新到redis缓存
	params := []interface{}{key}
	params = append(params, "Id", id)
	if data.Username != "" {
		params = append(params, "Username", data.Username)
		params = append(params, "Blacktime", data.Blacktime)
		params = append(params, "Realname", data.Realname)
		params = append(params, "Mobile", data.Mobile)
		params = append(params, "Address", data.Address)
		params = append(params, "SysCreated", data.SysCreated)
		params = append(params, "SysUpdated", data.SysUpdated)
		params = append(params, "SysIp", data.SysIp)
	}
	_, err := rds.Do("HMSET", params...)
	if err != nil {
		log.Println("user_service.setByCache HMSET params=", params, ", error=", err)
	}
}

// 数据更新了，直接清空缓存数据
func (s *userService) updateByCache(data *models.LwUser, columns []string) {
	if data == nil || data.Id <= 0 {
		return
	}
	// 集群模式，redis缓存
	key := fmt.Sprintf("info_user_%d", data.Id)
	rds := datasource.InstanceCache()
	// 删除redis中的缓存
	rds.Do("DEL", key)
}
