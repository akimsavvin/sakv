package config

import "github.com/ilyakaznacheev/cleanenv"

type Config struct {
	Logging Logging `yaml:"logging"`
	Engine  Engine  `yaml:"engine"`
	Network Network `yaml:"network"`
}

type Logging struct {
	Level  string `yaml:"level" env-default:"info"`
	Output string `yaml:"output" env-default:"log/output.log"`
}

type Engine struct {
	Type string `yaml:"type" env-default:"in_memory"`
}

type Network struct {
	Addr        string `yaml:"address" env-default:"127.0.0.1:3223"`
	MaxConns    int    `yaml:"max_connections" env-default:"10"`
	MaxMsgSize  string `yaml:"max_message_size" env-default:"2kb"`
	IdleTimeout string `yaml:"idle_timeout" env-default:"5m"`
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
