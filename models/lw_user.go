package models

type LwUser struct {
	Id         int    `xorm:"not null pk autoincr comment('主键') INT(11)"`
	Username   string `xorm:"comment('用户名') VARCHAR(50)"`
	Blacktime  int    `xorm:"comment('黑名单限制到期时间') INT(11)"`
	Realname   string `xorm:"comment('联系人') VARCHAR(50)"`
	Mobile     string `xorm:"comment('手机号') VARCHAR(50)"`
	Address    string `xorm:"comment('联系地址') VARCHAR(255)"`
	SysCreated int    `xorm:"comment('创建时间') INT(11)"`
	SysUpdated int    `xorm:"comment('修改时间') INT(11)"`
	SysIp      string `xorm:"comment('IP地址') VARCHAR(50)"`
}
