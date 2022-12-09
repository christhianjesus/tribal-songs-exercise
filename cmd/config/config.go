package config

type Config struct {
	Host string `env:"HOST,default=0.0.0.0"`
	Port string `env:"PORT,default=8080"`
}

func (c Config) Addr() string {
	return c.Host + ":" + c.Port
}
