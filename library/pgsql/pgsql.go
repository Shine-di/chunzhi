package pgsql

import (
	"context"
	"fmt"
	"game-test/library/config"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"go.elastic.co/apm/module/apmgorm"
	"go.uber.org/zap"
	"os"
	"regexp"
	"strconv"
	"sync"
)

var ConfigMap = make(map[string]*config.DbConfig, 0)
var DEFAULT_DB = "common"
var ORIGIN_DB = "game-origin"

type PgDbInfo struct {
	ServiceName string
	Env         string
	DbConfig    *config.DbConfig
	Conn        *gorm.DB
}

var pgInstanceMap = make(map[string]*PgDbInfo, 1)
var lock = &sync.Mutex{}
var onceLoadDb sync.Once

func init() {
	var dbConf config.DbConfig
	config.GetPgDbConfig(&dbConf, ORIGIN_DB)
	fmt.Println()
	ConfigMap[ORIGIN_DB] = &dbConf
}

func LoadDbConn(ctx context.Context) *gorm.DB {

	conn := GetDb()
	if conn == nil {
		return nil
	}

	onceLoadDb.Do(func() {
		apmgorm.RegisterCallbacks(conn)
	})

	conn = apmgorm.WithContext(ctx, conn)

	return conn
}

func GetDb() *gorm.DB {
	pgInstance := LoadPgDb(ORIGIN_DB, config.GetEnv())
	if pgInstance == nil {
		return nil
	}
	return pgInstance.CheckAndReturnConn()
}

func LoadPgDb(serviceName, env string) *PgDbInfo {
	_, ok := pgInstanceMap[serviceName]
	if !ok {

		lock.Lock()
		defer lock.Unlock()

		_, recheck := pgInstanceMap[serviceName]
		if !recheck {
			var dbConf *config.DbConfig
			if val, ok := ConfigMap[serviceName]; ok {
				dbConf = val
			} else {
				dbConf = ConfigMap[ORIGIN_DB]
			}
			fmt.Println(dbConf, ConfigMap, serviceName)
			PgDbInfo := new(PgDbInfo)
			PgDbInfo.DbConfig = dbConf
			PgDbInfo.Env = env
			PgDbInfo.ServiceName = serviceName

			if len(os.Getenv("DB")) > 0 {
				PgDbInfo.DbConfig.Database = os.Getenv("DB")
			}

			var convErr error
			if os.Getenv("OPEN") != "" {
				open := os.Getenv("OPEN")
				PgDbInfo.DbConfig.Open, convErr = strconv.ParseInt(open, 10, 64)
				if convErr != nil {
					PgDbInfo.DbConfig.Open = 0
				}
			}

			if os.Getenv("IDLE") != "" {
				open := os.Getenv("IDLE")
				PgDbInfo.DbConfig.Idle, convErr = strconv.ParseInt(open, 10, 64)
				if convErr != nil {
					PgDbInfo.DbConfig.Idle = 0
				}
			}

			PgDbInfo.InitConnect()
			pgInstanceMap[serviceName] = PgDbInfo
		}
	}

	return pgInstanceMap[serviceName]
}

func (PgDbInfo *PgDbInfo) InitConnect() {
	dbConf := PgDbInfo.DbConfig
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbConf.User, dbConf.Password, dbConf.Host, dbConf.Port, dbConf.Database)
	if db, e := gorm.Open("postgres", connStr); e != nil {
		log.Warn("load postage fail", zap.Error(e), zap.String("conn", connStr))
		PgDbInfo.Conn = nil
	} else {
		log.Warn("load postage success", zap.String("conn", connStr))
		if dbConf.Idle == 0 {
			dbConf.Idle = 5
		}
		if dbConf.Open == 0 {
			dbConf.Open = 20
		}
		db.LogMode(true)
		db.DB().SetMaxIdleConns(int(dbConf.Idle))
		db.DB().SetMaxOpenConns(int(dbConf.Open))
		PgDbInfo.Conn = db
	}
}

func (PgDbInfo *PgDbInfo) CheckAndReturnConn() *gorm.DB {
	if PgDbInfo.Conn == nil {
		lock.Lock()
		defer lock.Unlock()
		if PgDbInfo.Conn == nil {
			PgDbInfo.InitConnect()
		}
	}

	if err := PgDbInfo.Conn.DB().Ping(); err != nil {
		log.Warn("load postage fail", zap.Error(err))
		PgDbInfo.Clean()
		return nil
	}

	return PgDbInfo.Conn
}

func (PgDbInfo *PgDbInfo) Clean() {
	if PgDbInfo.Conn != nil {
		log.Warn("close postage conn")
		errClean := PgDbInfo.Conn.Close()
		if errClean != nil {
			log.Warn("close postage conn fail", zap.Error(errClean))
		}
	}

	PgDbInfo.Conn = nil
}

func (PgDbInfo *PgDbInfo) EscapeNotFound(errGet error) error {
	if errGet != nil {
		reg := regexp.MustCompile(`not found`)
		finder := reg.FindAllString(errGet.Error(), -1)
		if len(finder) > 0 {
			return nil
		}
	}

	return errGet
}
