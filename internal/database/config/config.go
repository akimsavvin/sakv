package config

import "github.com/ilyakaznacheev/cleanenv"

type Config struct {
	Logging Logging `yaml:"logging"`
	Engine  Engine  `yaml:"engine"`
	Network Network `yaml:"network"`
}

type Logging struct {
	Level  string `yaml:"level"`
	Output string `yaml:"output"`
}

type Engine struct {
	Type string `yaml:"type"`
}

type Network struct {
	Addr        string `yaml:"address"`
	MaxConns    int    `yaml:"max_connections"`
	MaxMsgSize  int    `yaml:"max_message_size"`
	IdleTimeout string `yaml:"idle_timeout"`
}

func New(filePath string) (*Config, error) {
	cfg := new(Config)
	return cfg, cleanenv.ReadConfig(filePath, cfg)
}

func MustNew(filePath string) *Config {
	cfg, err := New(filePath)
	if err != nil {
		panic(err)
	}
	return cfg
}
