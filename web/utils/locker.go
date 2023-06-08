package utils

import (
	"Lucky_Wheel/datasource"
	"fmt"
)

func getLuckLockKey(uid int) string {
	return fmt.Sprintf("lucky_lock_%d", uid)
}

func LockLucky(uid int) bool {
	key := getLuckLockKey(uid)
	cacheObj := datasource.InstanceCache()
	rs, _ := cacheObj.Do("SET", key, 1, "EX", 3, "NX")
	return rs == "ok"
}

func UnLockLucky(uid int) bool {
	key := getLuckLockKey(uid)
	cacheObj := datasource.InstanceCache()
	rs, _ := cacheObj.Do("DEL", key)
	return rs == "ok"
}
