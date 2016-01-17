# MrRedis  

Mesos runs Redis.

<img src="./logo.jpg" width="20%" height="20%"> 

A minimalistic framework for Redis workload on Apache Mesos

This framework supports the following freatures

 * Creates/Maintains Single Redis-instance
 * Creates/Maintains Redis-Instances with Master-Slave setup 
 * Vertical Auto/Manual scaling of a running redis-instance in terms of memory
 * Provides a cli to manage/monitor the redis instances those are being created 
 * A centralized persistance layer enabled by etcd

For example

```
$mrr create --name=app1-cache --mem=2G 
OK: Job Submitted to the framework
```

The cli itslef is async in nature as it does not wait for the operation to complete

```
$mrr status --name=app1-cache 
Status		= RUNNING
IP:PORT		= 176.134.0.10:45001
MemoryAvail	= 2GB
MemoryUsed	= 100MB
Slaves		= None
```

### Sample Run
After cloning the project and setting up the GOPATH for dependent libraries (should use go version 1.5)
```
$cd exec
$go build -o MrRedisExecutor main.go
$cd ../sched
$go build main.go
$./main -config="./config.json"
2016/01/17 16:35:11 *****************************************************************
2016/01/17 16:35:11 *********************Starting MrRedis-Scheduler******************
2016/01/17 16:35:11 *****************************************************************
2016/01/17 16:35:11 Starting the HTTP server at port 8080
```

The configuration file should be of json format

```
$cat config.json
{
        "MasterIP":"10.11.12.13",
        "MasterPort":"5050",
        "ExecutorPath":"/home/ubuntu/MrRedis/exec/MrRedisExecutor",
        "DBType":"etcd",
        "DBEndPoint": "http://11.12.13.14:2379",
        "ArtifactIP": "12.13.14.15"
}

```

Please substitute appropriate values with respect to your enviroment in the above config file for MasterIP/Port, ExecutorPath, DBEndPoint and IP adddres of this scheduler's VM that is accessible from the lsaves for artifactIP

### Installation Instruction
Please Note the pkg dependency management will be done by godep, but we will hold integrating it as the next destination (Transfer this project) of this project from my username is still not clear.
TODO

### Contribution Guidlines
TODO

### Documentation 
TODO

### Future Plans

- [ ] Support REDIS 3.0 Cluster 
- [ ] Support a Proxy mechanism to expose Redis Instance Endpoint
- [ ] Build a UI for Create/Maintain/Monitor the entier redis framework
- [ ] Benchmakr Utility for testing the RedisFramework 
