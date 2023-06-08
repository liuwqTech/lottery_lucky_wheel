package service

import (
	comm "Lucky_Wheel/common"
	"Lucky_Wheel/dao"
	"Lucky_Wheel/datasource"
	"Lucky_Wheel/models"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"log"
	"sync"
)

// IP信息，可以缓存(本地或者redis)，有更新的时候，再根据具体情况更新缓存
var cachedBlackipList = make(map[string]*models.LwBlackip)
var cachedBlackipLock = sync.Mutex{}

type BlackIpService interface {
	GetAll(page, size int) []models.LwBlackip
	CountAll() int64
	Search(ip string) []models.LwBlackip
	Get(id int) *models.LwBlackip
	Delete(id int) error
	Update(data *models.LwBlackip, columns []string) error
	Create(data *models.LwBlackip) error
	GetByIp(ip string) *models.LwBlackip
}

type blackIpService struct {
	dao *dao.BlackIpDao
}

func NewBlackIpService() BlackIpService {
	return &blackIpService{
		dao: dao.NewBlackIpDao(datasource.InstanceDbMaster()),
	}
}

func (s *blackIpService) GetAll(page, size int) []models.LwBlackip {
	return s.dao.GetAll(page, size)
}

func (s *blackIpService) CountAll() int64 {
	return s.dao.CountAll()
}

func (s *blackIpService) Search(ip string) []models.LwBlackip {
	return s.dao.Search(ip)
}

func (s *blackIpService) Get(id int) *models.LwBlackip {
	return s.dao.Get(id)
}

func (s *blackIpService) Delete(id int) error {
	return s.dao.Delete(id)
}

func (s *blackIpService) Update(data *models.LwBlackip, columns []string) error {
	// 先更新缓存的数据
	s.updateByCache(data, columns)
	// 再更新数据的数据
	return s.dao.Update(data, columns)
}

func (s *blackIpService) Create(data *models.LwBlackip) error {
	return s.dao.Create(data)
}

// 根据IP读取IP的黑名单数据
func (s *blackIpService) GetByIp(ip string) *models.LwBlackip {
	// 先从缓存中读取数据
	data := s.getByCache(ip)
	if data == nil || data.Ip == "" {
		// 再从数据库中读取数据
		data = s.dao.GetByIp(ip)
		if data == nil || data.Ip == "" {
			data = &models.LwBlackip{Ip: ip}
		}
		s.setByCache(data)
	}
	return data
}

func (s *blackIpService) getByCache(ip string) *models.LwBlackip {
	// 集群模式，redis缓存
	key := fmt.Sprintf("info_blackip_%s", ip)
	rds := datasource.InstanceCache()
	dataMap, err := redis.StringMap(rds.Do("HGETALL", key))
	if err != nil {
		log.Println("blackip_service.getByCache HGETALL key=", key, ", error=", err)
		return nil
	}
	dataIp := comm.GetStringFromStringMap(dataMap, "Ip", "")
	if dataIp == "" {
		return nil
	}
	data := &models.LwBlackip{
		Id:         int(comm.GetInt64FromStringMap(dataMap, "Id", 0)),
		Ip:         dataIp,
		Blacktime:  int(comm.GetInt64FromStringMap(dataMap, "Blacktime", 0)),
		SysCreated: int(comm.GetInt64FromStringMap(dataMap, "SysCreated", 0)),
		SysUpdated: int(comm.GetInt64FromStringMap(dataMap, "SysUpdated", 0)),
	}
	return data
}

func (s *blackIpService) setByCache(data *models.LwBlackip) {
	if data == nil || data.Ip == "" {
		return
	}
	// 集群模式，redis缓存
	key := fmt.Sprintf("info_blackip_%s", data.Ip)
	rds := datasource.InstanceCache()
	// 数据更新到redis缓存
	params := []interface{}{key}
	params = append(params, "Ip", data.Ip)
	if data.Id > 0 {
		params = append(params, "Blacktime", data.Blacktime)
		params = append(params, "SysCreated", data.SysCreated)
		params = append(params, "SysUpdated", data.SysUpdated)
	}
	_, err := rds.Do("HMSET", params...)
	if err != nil {
		log.Println("blackip_service.setByCache HMSET params=", params, ", error=", err)
	}
}

// 数据更新了，直接清空缓存数据
func (s *blackIpService) updateByCache(data *models.LwBlackip, columns []string) {
	if data == nil || data.Ip == "" {
		return
	}
	// 集群模式，redis缓存
	key := fmt.Sprintf("info_blackip_%s", data.Ip)
	rds := datasource.InstanceCache()
	// 删除redis中的缓存
	rds.Do("DEL", key)
}
