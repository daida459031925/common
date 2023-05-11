package mysql

import (
	"github.com/daida459031925/common/file"
	"github.com/daida459031925/common/fmt"
	"github.com/daida459031925/common/sql"
	"github.com/daida459031925/common/sql/redis"
	"github.com/daida459031925/common/time"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type config struct {
	Mysql struct {
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		DbName   string `yaml:"dbName"`
		MaxConn  int    `yaml:"maxConn"`
		MaxOpen  int    `yaml:"maxOpen"`
		MaxTime  int    `yaml:"maxTime"`
	} `yaml:"mysql"`
}

func NewDbConfig(filePath string) (*config, error) {
	return file.NewConfig[config](filePath)
}

func (c *config) Connect() (*sql.BaseDb, error) {
	//创建sql连接内容
	mysqlUser := c.Mysql.User
	mysqlPassword := c.Mysql.Password
	mysqlHost := c.Mysql.Host
	mysqlPort := c.Mysql.Port
	mysqlDBName := c.Mysql.DbName
	//连接池
	MaxConn := c.Mysql.MaxConn
	MaxOpen := c.Mysql.MaxOpen
	MaxTime := c.Mysql.MaxTime

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", mysqlUser, mysqlPassword, mysqlHost, mysqlPort, mysqlDBName)
	//dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Asia%2FShanghai", mysqlUser, mysqlPassword, mysqlHost, mysqlPort, mysqlDBName)
	db, e := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if e != nil {
		return nil, e
	}

	sqlDB, e := db.DB()
	if e != nil {
		return nil, e
	}

	sqlDB.SetMaxIdleConns(MaxConn)
	sqlDB.SetMaxOpenConns(MaxOpen)
	sqlDB.SetConnMaxLifetime(time.GetHour(MaxTime))

	return &sql.BaseDb{
		Db:    c,
		Gdb:   *db,
		Redis: nil,
	}, nil
}

func (c *config) ConnectCheck(checkFilePath string) (*sql.BaseDb, error) {
	db, e := c.Connect()
	if e != nil {
		return nil, e
	}
	// 设置缓存
	rc, e := redis.NewRedisConfig(checkFilePath)
	if e != nil {
		return nil, e
	}
	cacheStore := rc.NewRedis()
	db.Redis = cacheStore
	return db, nil
}
