package discovery

import (
	"context"
	"errors"
	"main/pkg/util"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	clientv3 "go.etcd.io/etcd/client/v3"
)

// ServiceDiscovery 服务发现
type ServiceDiscovery struct {
	//etcd 客户端
	cli *clientv3.Client
	// 服务列表
	serverList map[string]string
	//读写锁
	lock sync.RWMutex
}

func DiscoveryService(endpoint []string) {
	s := NewServiceDiscovery(endpoint)

	//传入user服务的前缀
	if err := s.DiscoveryService(util.PM["user"]); err != nil {
		logrus.Fatalln("err : ", err)
		return
	}

	//传入task服务的前缀
	if err := s.DiscoveryService(util.PM["user"]); err != nil {
		logrus.Fatalln("err : ", err)
		return
	}

	for k, v := range s.serverList {
		logrus.Println("已发现服务:", k, " :: ", v)
	}
}

// NewServiceDiscovery 新建服务发现
func NewServiceDiscovery(endpoint []string) *ServiceDiscovery {
	// 初始化etcd client
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoint,
		DialTimeout: time.Duration(5) * time.Second,
	})
	if err != nil {
		logrus.Fatalln(err)
	}

	return &ServiceDiscovery{
		cli:        cli,
		serverList: make(map[string]string),
		lock:       sync.RWMutex{},
	}
}

// 服务发现
func (s *ServiceDiscovery) DiscoveryService(servicePrefix string) (err error) {
	kvCli := clientv3.NewKV(s.cli)
	// 先通过前缀查询所有val
	valResp, err := kvCli.Get(context.Background(), servicePrefix, clientv3.WithPrefix())
	if err != nil {
		logrus.Fatal("etcd error")
		return err
	}

	// 服务未找到
	if len(valResp.Kvs) == 0 {
		return errors.New("服务未找到")
	}

	// 将所有服务遍历进list
	for _, kv := range valResp.Kvs {
		s.PutServiceInList(string(kv.Key), string(kv.Value))
	}

	// 开启协程启动watcher
	go s.watcher(servicePrefix)
	return nil
}

// watcher 监听Key的前缀
func (s *ServiceDiscovery) watcher(servicePrefix string) {
	// 获取etcd的watch
	watcher := clientv3.NewWatcher(s.cli)

	// 通过watch服务的前缀 ,检测该前缀下的所有变化
	watchChan := watcher.Watch(context.Background(), servicePrefix, clientv3.WithPrefix())

	for watchResponse := range watchChan {
		for _, e := range watchResponse.Events {
			switch e.Type {
			case clientv3.EventTypePut: // 新增或修改
				s.PutServiceInList(string(e.Kv.Key), string(e.Kv.Value))
			case clientv3.EventTypeDelete: // 删除
				s.RemoveServiceInList(string(e.Kv.Key))
			}
		}
	}
}
func (s *ServiceDiscovery) PutServiceInList(k, v string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.serverList[k] = v
}

func (s *ServiceDiscovery) RemoveServiceInList(k string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	delete(s.serverList, k)
}
