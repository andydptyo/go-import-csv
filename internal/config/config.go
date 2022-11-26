package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Database *Mysql `yaml:"database"`
}
type Mysql struct {
	DbName            string `yaml:"dbname"`
	Host              string `yaml:"host"`
	Port              string `yaml:"port"`
	User              string `yaml:"user"`
	Password          string `yaml:"password"`
	MaxLifeTime       int    `yaml:"maxLifeTime"`
	MaxOpenConnection int    `yaml:"maxOpenConnection"`
	MaxIdleConnection int    `yaml:"maxIdleConnection"`
}

func FromFile(file string) (*Config, error) {
	var config Config
	content, err := os.ReadFile(file)

	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(content, &config)

	if err != nil {
		return nil, err
	}

	return &config, err
}

func (ms *Mysql) GetDsn() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", ms.User, ms.Password, ms.Host, ms.Port, ms.DbName)
}
