package db

import (
	"fmt"
	"time"

	"github.com/Fromsko/gouitls/logs"
	"github.com/go-redis/redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var log *logrus.Logger

func init() {
	log = logs.InitLogger()
}

type Redis struct {
	Addr         string
	PassWord     string
	Db           int
	PoolSize     int
	MaxRetries   int
	MinIdleCons int
}

type MySQL struct {
	Addr           string
	UserName       string
	PassWord       string
	Database       string
	MinIdleCons   int
	MaxOpenCons   int
	ConMaxLeftTime int
}

func NewMysql(c *MySQL) (*gorm.DB, error) {
	// recover
	defer func() {
		if err := recover(); err != nil {
			log.Error(err)
		}
	}()

	// mysql数据库连接
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.UserName,
		c.PassWord,
		c.Addr,
		c.Database,
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// 连接池参数设置
	sqlDB.SetMaxIdleConns(c.MinIdleCons)
	sqlDB.SetMaxOpenConns(c.MaxOpenCons)
	sqlDB.SetConnMaxLifetime(time.Hour * time.Duration(c.ConMaxLeftTime))

	return db, nil
}

func NewRedis(c *Redis) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr: c.Addr,
		//Password:     c.Redis.Password,
		DB:           int(c.Db),
		PoolSize:     int(c.PoolSize),     // 连接池数量
		MinIdleConns: int(c.MinIdleCons), // 好比最小连接数
		MaxRetries:   int(c.MaxRetries),   // 命令执行失败时，最多重试多少次，默认为0即不重试
	})
	log.Error("Failed to connect to redis server.")
	return rdb
}
