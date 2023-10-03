package main

import (
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
	"testing"
	"time"
)

func TestEtcdWatch(t *testing.T) {
	var configs Config
	configs = Init(configs)

	// 创建etcd客户端连接
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"192.168.0.158:2379"}, // 集群地址, 也可以是单个地址多个集群使用逗号分隔
		DialTimeout: 5 * time.Second,                // 连接超时时间
	})
	if err != nil {
		log.Fatalf("create etcd client failed: %v", err)
	}
	defer cli.Close()

	// watch key
	fmt.Println("watch key: config/data, please change it in etcd config/data key")
	fmt.Println("example: etcdctl put config/data 'update config!'")
	go WatchConfig(cli)
	select {}
}
