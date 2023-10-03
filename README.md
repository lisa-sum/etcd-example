# ETCD example

## Install
1. click [github releases](https://github.com/etcd-io/etcd/releases/) to download your os etcd package
2. unzip etcd
3. create etcd config file and data dir
    ```shell
    mkdir -p /home/data/etcd/etcd.d
    ```
4. run default etcd server
    ```shell
    nohup etcd &
    ```
   
OR  config etcd file

1. config etcd file
    ```shell
    vi /home/data/etcd/etcd.d/etcd.conf.yml
    ```
2. run etcd server
    ```shell
    nohup etcd --config-file=/home/data/etcd/etcd.d/etcd.conf.yml &
    ``` 
   
## Use 

### CLI

#### Put

```shell
./etcdctl put /myapp/config 1
```

#### Delete

```shell
./etcdctl del /myapp/config
```

#### Get

```shell
./etcdctl get /myapp/config
```

#### Watch

```shell
./etcdctl watch /myapp/config
```

#### List

```shell
./etcdctl ls /myapp/config
```

### Golang
CRUD example:
check [main.go](main.go)

Watch example:
check [watch.go](watch_test.go)
