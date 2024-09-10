package server

type Config struct {
	Host string `env:"HOST,required"`
	Port string `env:"PORT,required"`
}
