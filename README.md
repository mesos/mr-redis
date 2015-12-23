# MrRedis

Mesos runs Redis

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

### Installation Instruction
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
