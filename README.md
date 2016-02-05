# MrRedis  

Mesos runs Redis.

<img src="./logo.jpg" width="20%" height="20%"> 

A minimalistic framework for Redis workload on Apache Mesos

This framework supports the following features

 * Creates/Maintains Single Redis-instance
 * Creates/Maintains Redis-Instances with Master-Slave setup 
 * A centralized persistance layer currently enabled by etcd

## Why MrRedis?
At [Huawei] (http://www.huawei.com/en/) we foresee creating, running and maintaing huge number of redis instances on our datacenters.  We intially evaluated few cluster managers for this job, but due to the specific requirements of 'redis' itslef those solutions did not satisfy most of our needs.  We quickly did a POC by writing a framework exclusively for Redis on Apache Mesos. Based on the outcome we decided to initate this project and work with the opensource community to build a robust custom framework for Redis which will be usefull for Huawei as well as rest of the world.

## Who should use MrRedis
* If your organization has a requirement of creating and maintaing huge number of redis service instances.
* If you are is planning to host a 'redis' as a Service 
* If redis instances need to be created in seconds and not in minutes
* If you are already using Apache Mesos as a Resource Manager for your Datacenter and want to add Redis workload to it


For example
```
$mrr help
NAME:
   mrr - MrRedis Command Line Interface

USAGE:
   mrr [global options] command [command options] [arguments...]

VERSION:
   0.0.0

COMMANDS:
   init, i      $mrr init <http://MrRedisEndPoint>
   create, c    Create a Redis Instance
   status, s    Status of a Redis Instance
   delete, d    Delete a Redis Instance
   help, h      Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h   show help
```
Help on a specific command
```
$mrr help create
NAME:
   mrr create - Create a Redis Instance

USAGE:
   mrr create [command options] [arguments...]

OPTIONS:
   --name, -n           Name of the Redis Instance
   --memory, -m "0"     Memory in MB
   --slaves, -s "0"     Number of Slaves
   --wait, -w           Wait for the Instnace to be create (by default the command is async)
   
```

The cli itself will be async in nature (by default) as it does not wait for the operation to complete

```
$mrr create -n testInst -m 200 -s 1
Attempting to Creating a Redis Instance (testInst) with mem=200 slaves=1
Instance Creation accepted..
Check $mrr status -n testInst for status
```
```
$mrr status -n testInst
Status = RUNNING
Type = MS
Capacity = 200
Master = 10.11.12.33:6389
        Slave0 = 10.11.12.32:6380
```

### Sample Build and Run
After cloning the project and setting up the GOPATH for dependent libraries (should use go version 1.5 or atlesat go 1.4)
Perform go get in all the three directories namely (sched, exec and cli) then follow this build order
```
$cd exec
$go build -o MrRedisExecutor main.go
$cd ../cli
$go build -o mrr
$cd ../sched
$go build main.go
```
to start the scheduler

```
$./main -config="./config.json"
2016/01/17 16:35:11 *****************************************************************
2016/01/17 16:35:11 *********************Starting MrRedis-Scheduler******************
2016/01/17 16:35:11 *****************************************************************
2016/01/17 16:35:11 Starting the HTTP server at port 8080
.
.
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

If you have a complicated Redis requirement then a simple comamnd like below 
```
$time mrr create -n hello50 -m 100 -s 50 -w true
Attempting to Creating a Redis Instance (hello50) with mem=100 slaves=50
Instance Creation accepted................
Instance Created.

real    0m14.269s
user    0m0.033s
sys     0m0.037s
```
will result in creating one redis instance with 1 master and 50 Slaves in less than 15 secs, Simples :-)

```
$mrr status -n hello50
Status = RUNNING
Type = MS
Capacity = 100
Master = 10.11.12.21:6380
        Slave0 = 10.11.12.31:6381
        Slave1 = 10.11.12.31:6383
        Slave2 = 10.11.12.31:6384
        Slave3 = 10.11.12.31:6385
        Slave4 = 10.11.12.31:6382
        Slave5 = 10.11.12.31:6386
        Slave6 = 10.11.12.31:6387
        Slave7 = 10.11.12.31:6388
        Slave8 = 10.11.12.31:6391
        Slave9 = 10.11.12.31:6392
        Slave10 = 10.11.12.31:6390
        Slave11 = 10.11.12.31:6389
        Slave12 = 10.11.12.31:6393
        Slave13 = 10.11.12.31:6394
        Slave14 = 10.11.12.31:6395
        Slave15 = 10.11.12.20:6380
        Slave16 = 10.11.12.20:6381
        Slave17 = 10.11.12.20:6383
        Slave18 = 10.11.12.20:6384
        Slave19 = 10.11.12.20:6387
        Slave20 = 10.11.12.20:6385
        Slave21 = 10.11.12.20:6386
        Slave22 = 10.11.12.20:6382
        Slave23 = 10.11.12.29:6380
        Slave24 = 10.11.12.29:6381
        Slave25 = 10.11.12.29:6382
        Slave26 = 10.11.12.29:6384
        Slave27 = 10.11.12.29:6385
        Slave28 = 10.11.12.29:6383
        Slave29 = 10.11.12.29:6387
        Slave30 = 10.11.12.29:6386
        Slave31 = 10.11.12.29:6389
        Slave32 = 10.11.12.29:6391
        Slave33 = 10.11.12.29:6392
        Slave34 = 10.11.12.29:6388
        Slave35 = 10.11.12.29:6390
        Slave36 = 10.11.12.29:6394
        Slave37 = 10.11.12.29:6395
        Slave38 = 10.11.12.29:6393
        Slave39 = 10.11.12.21:6383
        Slave40 = 10.11.12.21:6384
        Slave41 = 10.11.12.21:6386
        Slave42 = 10.11.12.21:6385
        Slave43 = 10.11.12.21:6387
        Slave44 = 10.11.12.21:6388
        Slave45 = 10.11.12.21:6390
        Slave46 = 10.11.12.21:6389
        Slave47 = 10.11.12.21:6391
        Slave48 = 10.11.12.21:6381
        Slave49 = 10.11.12.21:6382
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
