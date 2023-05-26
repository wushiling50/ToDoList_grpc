package login

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	clientv3 "go.etcd.io/etcd/client/v3"
)

// 实例化注册器
type Register struct {
	EtcdAddrs   []string //地址
	DialTimeout int      //超时时间

	closeCh     chan struct{}                           //是否关闭
	leasesID    clientv3.LeaseID                        //租约
	keepAliveCh <-chan *clientv3.LeaseKeepAliveResponse //检验心跳

	srvInfo Server           //服务
	srvTTL  int64            //存活时间
	cli     *clientv3.Client //客户端
	logger  *logrus.Logger   //日志
}

// 基于etcd创建注册器
func NewRegister(etcdAddrs []string, logger *logrus.Logger) *Register {
	return &Register{
		EtcdAddrs:   etcdAddrs,
		DialTimeout: 3,
		logger:      logger,
	}
}

// 初始化注册器
func (r *Register) Register(srvInfo Server, ttl int64) (chan<- struct{}, error) {
	var err error
	if strings.Split(srvInfo.Address, ":")[0] == "" {
		return nil, errors.New("无效的地址")
	}

	if r.cli, err = clientv3.New(clientv3.Config{
		Endpoints:   r.EtcdAddrs,
		DialTimeout: time.Duration(r.DialTimeout) * time.Second,
	}); err != nil {
		return nil, err
	}

	r.srvInfo = srvInfo
	r.srvTTL = ttl

	if err = r.register(); err != nil {
		return nil, err
	}

	r.closeCh = make(chan struct{})
	go r.KeepAlive()
	return r.closeCh, nil
}

// 构建etcd自带的实例
func (r *Register) register() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(r.DialTimeout)*time.Second)

	defer cancel()

	//申请租约
	leaseResp, err := r.cli.Grant(ctx, r.srvTTL)
	if err != nil {
		return err
	}

	//将租约赋值
	r.leasesID = leaseResp.ID

	//定义KeepAlive
	if r.keepAliveCh, err = r.cli.KeepAlive(context.Background(), r.leasesID); err != nil {
		return err
	}

	data, err := json.Marshal(r.srvInfo)
	if err != nil {
		return err
	}

	//push服务到服务注册中
	_, err = r.cli.Put(context.Background(), CreateRegisterPath(r.srvInfo), string(data), clientv3.WithLease(r.leasesID))

	return err

}

func (r *Register) KeepAlive() error {
	ticker := time.NewTicker(time.Duration(r.srvTTL) * time.Second)
	//关闭就注销，存活或者超时都注册
	for {
		select {
		case <-r.closeCh:
			if err := r.unregister(); err != nil {
				r.logger.Error("注销失败, error: ", err)
			}
			//废除租约
			if _, err := r.cli.Revoke(context.Background(), r.leasesID); err != nil {
				r.logger.Error("废除租约失败, error: ", err)
			}
		case res := <-r.keepAliveCh:
			if res == nil {
				if err := r.register(); err != nil {
					r.logger.Error("注册服务错误, error: ", err)
				}
			}
			//超时器
		case <-ticker.C:
			if r.keepAliveCh == nil {
				if err := r.register(); err != nil {
					r.logger.Error("注册服务错误, error: ", err)
				}
			}
		}
	}
}

// 在etcd中删除节点
func (r *Register) unregister() error {
	_, err := r.cli.Delete(context.Background(), CreateRegisterPath(r.srvInfo))
	return err
}
