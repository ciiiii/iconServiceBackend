package config

import (
	"sync"
	"os"
	"strings"
	"github.com/BurntSushi/toml"
	"github.com/caarlos0/env/v6"
)

type config struct {
	Cos struct {
		BucketName string `env:"BUCKET"`
		Region     string `env:"REGION"`
		SecretID   string `env:"ID"`
		SecretKey  string `env:"KEY"`
	}
	Config struct {
		Mode string `env:"MODE"`
		Port string `env:"PORT"`
	}
}

var (
	c    config
	once sync.Once
)

func Parser() config {
	once.Do(func() {
		if os.Getenv("PORT") != "" {
			if err := env.Parse(&c); err != nil {
				panic(err)
				panic("[config] parse env error")
			}
		} else {
			rootPath, _ := os.Getwd()
			confPath := strings.Join([]string{rootPath, "conf.toml"}, "/")
			if _, err := toml.DecodeFile(confPath, &c); err != nil {
				panic("[config] need conf.toml")
			}
		}
	})
	return c
}
