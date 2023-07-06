package mysql

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"xorm.io/xorm"
)

type Config struct {
	Addr     string `yaml:"addr"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Db       string `yaml:"database"`
	Timeout  string `yaml:"timeout"`
}

var Client *xorm.Engine

func InitDB(config *Config) error {
	var err error
	if config.Timeout == "" {
		config.Timeout = "10s"
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local&timeout=%s", config.Username,
		config.Password, config.Addr, config.Db, config.Timeout)
	// Open 连接
	Client, err = xorm.NewEngine("mysql", dsn)
	if err != nil {
		return err
	}

	return nil
}
