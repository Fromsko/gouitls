package db

import (
	"fmt"
	"time"

	"github.com/Fromsko/gouitls/logs"
	"github.com/go-redis/redis/v9"
	"gorm.io/gorm"
	"gorm.io/driver/mysql"
	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

func init() {
	log = logs.InitLogger()
}

type DBManage struct {
	Rdb *redis.Client
	Mdb *gorm.DB
}

type Conf struct {
	Redis struct {
		Addr         string
		PassWord     string
		Db           int
		PoolSize     int
		MaxRetries   int
		MinIdleConns int
	}
	MySQL struct {
		Addr           string
		UserName       string
		PassWord       string
		Database       string
		MinIdleConns   int
		MaxOpenConns   int
		ConMaxLeftTime int
	}
}

func NewDBConfig() (config *Conf) {
	config = &Conf{
		Redis: struct {
			Addr         string
			PassWord     string
			Db           int
			PoolSize     int
			MaxRetries   int
			MinIdleConns int
		}{},
		MySQL: struct {
			Addr           string
			UserName       string
			PassWord       string
			Database       string
			MinIdleConns   int
			MaxOpenConns   int
			ConMaxLeftTime int
		}{},
	}
	return config
}

func NewMysql(c *Conf) (*gorm.DB, error) {
	// recover
	defer func() {
		if err := recover(); err != nil {
			log.Error(err)
		}
	}()

	// mysql数据库连接
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.MySQL.UserName,
		c.MySQL.PassWord,
		c.MySQL.Addr,
		c.MySQL.Database,
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
	sqlDB.SetMaxIdleConns(c.MySQL.MinIdleConns)
	sqlDB.SetMaxOpenConns(c.MySQL.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Hour * time.Duration(c.MySQL.ConMaxLeftTime))

	return db, nil
}

func NewRedis(c *Conf) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr: c.Redis.Addr,
		//Password:     c.Redis.Password,
		DB:           int(c.Redis.Db),
		PoolSize:     int(c.Redis.PoolSize),     // 连接池数量
		MinIdleConns: int(c.Redis.MinIdleConns), // 好比最小连接数
		MaxRetries:   int(c.Redis.MaxRetries),   // 命令执行失败时，最多重试多少次，默认为0即不重试
	})
	log.Error("Failed to connect to redis server.")
	return rdb
}
