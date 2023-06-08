package conf

type RedisConfig struct {
	Host      string
	Port      int
	User      string
	Pwd       string
	IsRunning bool
}

var RedisCacheList = []RedisConfig{
	{
		Host:      "127.0.0.1",
		Port:      6379,
		User:      "",
		Pwd:       "",
		IsRunning: true,
	},
}

var RedisCache RedisConfig = RedisCacheList[0]
