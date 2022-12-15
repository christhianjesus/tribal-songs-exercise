package config

type redis struct {
	Host string `env:"REDIS_HOST,default=localhost"`
	Port string `env:"REDIS_PORT,default=6379"`
}
