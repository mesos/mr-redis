# MrRedis  

Mesos runs Redis.

<img src="./logo.jpg" width="20%" height="20%"> 

A minimalistic framework for Redis workload on Apache Mesos

This framework supports the following features

 * Creates/Maintains Single Redis-instance
 * Creates/Maintains Redis-Instances with Master-Slave setup 
 * Vertical Auto/Manual scaling of a running redis-instance in terms of memory
 * Provides a cli to manage/monitor the redis instances those are being created 
 * A centralized persistance layer currently enabled by etcd

For example

```
$mrr create --name=app1-cache --mem=2G 
OK: Job Submitted to the framework
```

The cli itself will be async in nature as it does not wait for the operation to complete

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

Please substitute appropriate values with respect to your enviroment in the above config file for MasterIP/Port, ExecutorPath, DBEndPoint and IP adddres of this scheduler's VM that is accessible from the slaves for artifactIP

If you have a complicated Redis requirement then a simple http comamnd like below 
```
$curl -X "POST" http://10.11.12.17:8080/v1/CREATE/1Master21Slaves/1024/1/21
```
will result in creating 1 master with 21 Slaves in less than a min, Simples :-)

```
# Replication
role:master
connected_slaves:21
slave0:ip=10.11.12.20,port=6381,state=online,offset=323,lag=0
slave1:ip=10.11.12.21,port=6382,state=online,offset=323,lag=0
slave2:ip=10.11.12.20,port=6382,state=online,offset=323,lag=0
slave3:ip=10.11.12.21,port=6383,state=online,offset=323,lag=0
slave4:ip=10.11.12.21,port=6384,state=online,offset=323,lag=0
slave5:ip=10.11.12.20,port=6383,state=online,offset=323,lag=0
slave6:ip=10.11.12.21,port=6385,state=online,offset=323,lag=0
slave7:ip=10.11.12.20,port=6384,state=online,offset=323,lag=0
slave8:ip=10.11.12.21,port=6386,state=online,offset=323,lag=0
slave9:ip=10.11.12.20,port=6385,state=online,offset=323,lag=0
slave10:ip=10.11.12.21,port=6387,state=online,offset=323,lag=0
slave11:ip=10.11.12.20,port=6386,state=online,offset=323,lag=0
slave12:ip=10.11.12.21,port=6388,state=online,offset=323,lag=0
slave13:ip=10.11.12.20,port=6387,state=online,offset=323,lag=0
slave14:ip=10.11.12.21,port=6389,state=online,offset=323,lag=0
slave15:ip=10.11.12.21,port=6390,state=online,offset=323,lag=0
slave16:ip=10.11.12.21,port=6391,state=online,offset=323,lag=0
slave17:ip=10.11.12.21,port=6392,state=online,offset=323,lag=0
slave18:ip=10.11.12.21,port=6393,state=online,offset=323,lag=0
slave19:ip=10.11.12.21,port=6394,state=online,offset=323,lag=0
slave20:ip=10.11.12.21,port=6395,state=online,offset=323,lag=0
master_repl_offset:323
repl_backlog_active:1
repl_backlog_size:1048576
repl_backlog_first_byte_offset:2
repl_backlog_histlen:322

```

### Installation Instruction
Please Note the pkg dependency management will be done by godep, but we will hold integrating it as the next destination (Transfer this project) of this project from current location is still not clear.
TODO

### Contribution Guidlines
We have ourselves fallen into pitfalls to arrive at working code faster, some simple rules more to inculcate in our own future work and for reference to contributors
Go already provides a nice documentation on coding conventions and guidelines; just try to adhere to that [Effective Go](https://golang.org/doc/effective_go.html) :-) 

Specifically 
- Format code using go fmt, if an already prebuilt auto formatter is not their in your editor
- We suggest using extensive comments, as this code base is still evolving
- Try to stress on error handling as per [Effective error handling in Go](https://golang.org/doc/effective_go.html#errors) (which we ourselves have probably missed at places)
- Please use this framework; We are looking forward for issues, and nothing greater then an issue and a fix. Nonetheless, if interested in contributing something specific, please raise an issue outright to let us know that you are doing "this"
- We have not set up tests and test code yet, this is one obvious area to contribute without saying 

### Documentation 
TODO

### Future Plans

- [ ] Support REDIS 3.0 Cluster 
- [ ] Support a Proxy mechanism to expose Redis Instance Endpoint
- [ ] Build a UI for Create/Maintain/Monitor the entier redis framework
- [ ] Benchmakr Utility for testing the RedisFramework 
