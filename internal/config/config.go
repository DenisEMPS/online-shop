package config

import (
	"flag"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env    string `yaml:"env" env-default:"local"`
	DB     DB     `yaml:"db"`
	Server Server `yaml:"server"`
	Redis  Redis  `yaml:"redis"`
}

type DB struct {
	Username string `yaml:"username" env-required:"true"`
	Password string `yaml:"password" env-required:"true"`
	Host     string `yaml:"host" env-default:"localhost"`
	Port     string `yaml:"port" env-default:"5436"`
	DBname   string `yaml:"dbname" env-default:"postgres"`
	SSLmode  string `yaml:"sslmode" env-default:"disabled"`
}

type Server struct {
	Port string `yaml:"port" env-default:"8080"`
}

type Redis struct {
	Host     string `yaml:"host" env-default:"cache"`
	Port     string `yaml:"port" env-default:"6379"`
	DB       int    `yaml:"DB" env-default:"0"`
	Password string `yaml:"password" env-required:"true"`
}

func MustLoad() *Config {
	path := fetchConfig()
	if path == "" {
		panic("config path is empty")
	}

	return MustLoadByPath(path)
}

func MustLoadByPath(path string) *Config {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config file does not exists:" + path)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic("failed to read config: " + err.Error())
	}

	return &cfg
}

func fetchConfig() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
