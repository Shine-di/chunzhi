package site_var

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"gitee.com/risewinter/data-lol/library/log"
	"go.etcd.io/etcd/clientv3"
	"go.uber.org/zap"
)

const TIMEOUT = 5 * time.Second

var once = sync.Once{}
var etcdInstance *Etcd
var env string

type EtcdConf struct {
	Urls     []string
	UserName string
	Password string
}

func getEtcdUrls() []string {
	etcdUrls := os.Getenv("ETCD_URLS")

	if len(etcdUrls) == 0 {
		log.Error("load ETCD_URLS blank")
		os.Exit(1)
	}

	urls := strings.Split(etcdUrls, ",")
	return urls
}

var getEtcdConf = &EtcdConf{
	Urls: getEtcdUrls(),
}

type Etcd struct {
	cli *clientv3.Client
}

type EtcdItem struct {
	Path  string `json:"path"`
	Value string `json:"value"`
}

func GetDefaultEtcdService() *Etcd {
	once.Do(func() {
		var etcdConfig *EtcdConf
		etcdConfig = getEtcdConf

		etcdInstance = initEtcd(etcdConfig)
	})
	return etcdInstance
}

func initEtcd(etcdConfig *EtcdConf) *Etcd {

	cli, err := clientv3.New(clientv3.Config{
		DialTimeout: TIMEOUT,
		Endpoints:   etcdConfig.Urls,
		Username:    etcdConfig.UserName,
		Password:    etcdConfig.Password,
	})
	if err != nil {
		log.Panic("connect etcd failed", zap.Error(err))
	}
	return &Etcd{cli: cli}
}

func (etcd *Etcd) Delete(keyPath string) error {
	kv := clientv3.KV(etcd.cli)
	ctx, _ := context.WithTimeout(context.Background(), TIMEOUT)
	res, err := kv.Delete(ctx, keyPath)
	log.Info("delete etcd value", zap.String("keyPath", keyPath), zap.Any("res", res), zap.Error(err))
	if err != nil {
		return err
	}
	return nil
}

func (etcd *Etcd) Put(keyPath, value string) error {
	kv := clientv3.KV(etcd.cli)
	ctx, _ := context.WithTimeout(context.Background(), TIMEOUT)
	res, err := kv.Put(ctx, keyPath, value)
	log.Info("put etcd value", zap.String("keyPath", keyPath), zap.Any("res", res), zap.Error(err))
	if err != nil {
		return err
	}
	return nil
}

func (etcd *Etcd) Get(keyPath string) (string, error) {
	kv := clientv3.KV(etcd.cli)
	ctx, _ := context.WithTimeout(context.Background(), TIMEOUT)
	res, err := kv.Get(ctx, keyPath)
	if err != nil {
		return "", err
	}

	for _, val := range res.Kvs {
		if string(val.Key[:]) == keyPath {
			return string(val.Value[:]), nil
		}
	}

	return "", errors.New("no value in etcd")
}

func (etcd *Etcd) GetBytes(keyPath string) ([]byte, error) {
	kv := clientv3.KV(etcd.cli)
	ctx, _ := context.WithTimeout(context.Background(), TIMEOUT)
	res, err := kv.Get(ctx, keyPath)
	if err != nil {
		return nil, err
	}

	for _, val := range res.Kvs {
		if string(val.Key[:]) == keyPath {
			return val.Value, nil
		}
	}

	return nil, errors.New("no value in etcd")
}

func (etcd *Etcd) GetList(keyPath string) ([]*EtcdItem, error) {
	kv := clientv3.KV(etcd.cli)
	ctx, _ := context.WithTimeout(context.Background(), TIMEOUT)
	res, err := kv.Get(ctx, keyPath, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}

	var items []*EtcdItem

	for _, val := range res.Kvs {
		item := new(EtcdItem)
		item.Path = string(val.Key[:])
		item.Value = string(val.Value[:])

		items = append(items, item)
	}

	return items, nil
}

// 监控 grpc host修改
func (etcd *Etcd) WatchGrpcHostModify() clientv3.WatchChan {
	grpcConfPath := fmt.Sprintf("/%s/config/grpc/", env)
	return etcd.cli.Watch(context.Background(), grpcConfPath, clientv3.WithPrefix())
}

func (etcd *Etcd) GetClient() *clientv3.Client {
	return etcd.cli
}

func init() {
	env = os.Getenv("ENVIRON")
	if env == "" {
		env = "develop"
	}
}
