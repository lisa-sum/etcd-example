package main

import (
	"bytes"
	"context"
	"fmt"
	"github.com/spf13/viper"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
	"time"
)

type Config struct {
	Server struct {
		HTTP struct {
			Addr    string        `yaml:"addr"`
			Timeout time.Duration `yaml:"timeout"`
		} `yaml:"http"`
		GRPC struct {
			Addr    string        `yaml:"addr"`
			Timeout time.Duration `yaml:"timeout"`
		} `yaml:"grpc"`
	} `yaml:"server"`
	Data struct {
		Postgres struct {
			Host            string `yaml:"host"`
			Username        string `yaml:"username"`
			Password        string `yaml:"password"`
			Port            string `yaml:"port"`
			DBName          string `yaml:"dbname"`
			SSLMode         string `yaml:"ssl_mode"`
			TimeZone        string `yaml:"time_zone"`
			MaxIdleConn     int    `yaml:"max_idle_conn"`
			MaxOpenConn     int    `yaml:"max_open_conn"`
			ConnMaxLifetime int    `yaml:"conn_max_lifetime"`
		} `yaml:"postgres"`
		Redis struct {
			Addr         string        `yaml:"addr"`
			Username     string        `yaml:"username"`
			Password     string        `yaml:"password"`
			DB           int           `yaml:"db"`
			PoolSize     int           `yaml:"pool_size"`
			ReadTimeout  time.Duration `yaml:"read_timeout"`
			WriteTimeout time.Duration `yaml:"write_timeout"`
			DialTimeout  time.Duration `yaml:"dial_timeout"`
		} `yaml:"redis"`
		Minio struct {
			Endpoint        string `yaml:"endpoint"`
			AccessKeyID     string `yaml:"accessKeyID"`
			SecretAccessKey string `yaml:"secretAccessKey"`
			UseSSL          bool   `yaml:"useSSL"`
		} `yaml:"minio"`
	} `yaml:"data"`
	Auth struct {
		SigningKey string `yaml:"signing_key"`
	} `yaml:"auth"`
}

func main() {
	var configs Config
	configs = Init(configs)

	// 创建etcd客户端连接
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"192.168.0.158:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatalf("create etcd client failed: %v", err)
	}
	defer cli.Close()

	WatchConfig(cli)

}
func Init(config Config) Config {
	// 实例化配置结构体
	config = Config{}

	// 设置服务器配置
	config.Server.HTTP.Addr = "0.0.0.0:30001"
	config.Server.HTTP.Timeout = 1 * time.Second
	config.Server.GRPC.Addr = "0.0.0.0:30002"
	config.Server.GRPC.Timeout = 1 * time.Second

	// 设置数据库配置
	config.Data.Postgres.Host = "192.168.0.158"
	config.Data.Postgres.Username = "postgres"
	config.Data.Postgres.Password = "263393"
	config.Data.Postgres.Port = "5432"
	config.Data.Postgres.DBName = "tiktok"
	config.Data.Postgres.SSLMode = "disable"
	config.Data.Postgres.TimeZone = "Asia/Shanghai"
	config.Data.Postgres.MaxIdleConn = 10
	config.Data.Postgres.MaxOpenConn = 100
	config.Data.Postgres.ConnMaxLifetime = 30

	// 设置 Redis 配置
	config.Data.Redis.Addr = "192.168.0.158:36379"
	config.Data.Redis.Username = "default"
	config.Data.Redis.Password = "$2a$14$qMybjd8FNiFWA8Z6zzpHhu0f5zH86CFysXWkaHq8ZibOLakFj/xbi"
	config.Data.Redis.DB = 0
	config.Data.Redis.PoolSize = 10
	config.Data.Redis.ReadTimeout = 10 * time.Second
	config.Data.Redis.WriteTimeout = 10 * time.Second
	config.Data.Redis.DialTimeout = 10 * time.Second

	// 设置 Minio 配置
	config.Data.Minio.Endpoint = "192.168.0.158:5000"
	config.Data.Minio.AccessKeyID = "Y2IGHXF948egBlce"
	config.Data.Minio.SecretAccessKey = "SODtldt1Myaqwr6PEw1C9CzFtVmJyIeO"
	config.Data.Minio.UseSSL = false

	// 设置认证配置
	config.Auth.SigningKey = "jwt_signing"

	// 打印生成的配置
	// fmt.Printf("%+v\n", config)

	return config
}

func WatchConfig(cli *clientv3.Client) {
	fmt.Printf("WatchConfig")
	// rch := cli.Watch(context.Background(), "/myapp/config/", clientv3.WithPrefix())
	rch := cli.Watch(context.Background(), "/myapp/config")
	for wresp := range rch {
		fmt.Printf("WatchConfig")
		for _, ev := range wresp.Events {
			switch ev.Type {
			case clientv3.EventTypePut:
				fmt.Printf("Config Putd!")
				fmt.Printf("[%s] %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
				updateConfig(string(ev.Kv.Value))
			case clientv3.EventTypeDelete:
				fmt.Printf("Config deleted!")
			}
		}
	}
}

func updateConfig(configStr string) {
	viper.SetConfigType("yaml")
	if err := viper.ReadConfig(bytes.NewBufferString(configStr)); err != nil {
		fmt.Printf("Error reading config with Viper: %v\n", err)
		return
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		fmt.Printf("Unable to decode into struct, %v\n", err)
		return
	}

	// Now 'config' is updated with the new values
	fmt.Printf("Updated Config: %+v\n", config)
}
