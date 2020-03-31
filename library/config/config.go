package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"

	"gitee.com/risewinter/data-lol/library/log"
	"gitee.com/risewinter/data-lol/library/site-var"
	"go.uber.org/zap"
)

var ENV string
var once sync.Once
var GrpcConfMap = &sync.Map{}
var grpcConfOnce = sync.Once{}
var watcher sync.Map

type RedisConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Password string `json:"password"`
}

type DbConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`
	Open     int64  `json:"open"`
	Idle     int64  `json:"idle"`
}

func GetEnv() string {
	once.Do(func() {
		env := os.Getenv("ENVIRON")
		if len(env) == 0 {
			log.Error("load environ empty")
			os.Exit(1)
		}
		ENV = env
	})
	return ENV
}

func ListAllDefaultGrpcHost() *sync.Map {
	grpcConfOnce.Do(func() {
		grpcConfPath := fmt.Sprintf("/%s/config/grpc/", GetEnv())
		items, err := site_var.GetDefaultEtcdService().GetList(grpcConfPath)
		if err != nil {
			log.Error("load grpc host list fail", zap.Error(err))
			return
		}
		for _, item := range items {
			grpcName := GetLastName(item.Path)
			GrpcConfMap.Store(grpcName, item.Value)
		}
	})

	return GrpcConfMap
}

func GetPgDbConfig(config *DbConfig, dbConfPath string) {
	//dbConfPath := os.Getenv("DB_CONFIG")
	if len(dbConfPath) == 0 {
		log.Info("load db config as common")
		dbConfPath = "common"
	}
	dbConfigPath := fmt.Sprintf("/%s/config/pgsql/%s", GetEnv(), dbConfPath)
	fmt.Println()
	resp, errGet := site_var.GetDefaultEtcdService().Get(dbConfigPath)
	if errGet != nil {
		log.Fatal("load db config fail", zap.Error(errGet))
	}
	errDecode := json.Unmarshal([]byte(resp), &config)
	if errDecode != nil {
		log.Fatal("get db config decode error", zap.String("resp", resp))
	}
}

func GetLastName(path string) string {
	var splicedList = strings.Split(path, "/")
	if len(splicedList) == 0 {
		return ""
	}
	return splicedList[len(splicedList)-1]
}

func GetRedisConfig(env string) (*RedisConfig, error) {
	s, err := site_var.GetDefaultEtcdService().Get(fmt.Sprintf("/%s/config/cache/redis", env))
	if err != nil {
		return nil, err
	}

	var redisConfig RedisConfig
	err = json.Unmarshal([]byte(s), &redisConfig)
	return &redisConfig, err
}

func GetRedisDbConfig(name string) int64 {
	s, err := site_var.GetDefaultEtcdService().Get(fmt.Sprintf("/%s/config/cache/redis-db", GetEnv()))
	if err != nil {
		return 0
	}
	confMap := make(map[string]int64)
	err = json.Unmarshal([]byte(s), &confMap)
	if err != nil {
		return 0
	}
	if c, ok := confMap[name]; ok {
		return c
	}
	if defaultValue, ok := confMap["default"]; ok {
		return defaultValue
	}
	return 0
}
