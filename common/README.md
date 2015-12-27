# Store
Storage layer of redis framework, this will give MrRedis the ability to remember things permanently.  

There are two types of objects this framework will persist.
* (Redis) Proc
* (Redis) Service Instance


### Redis Proc
A `redis-server` process running in any of the `mesos-slave` is called a **Redis Proc** (Redis Process). 

A Redis **Proc** has the following properties
* In the datacenter each Proc is identified by a UID
* It belongs to a Service Instance
* It binds to a particular port
* It can either be a Redis Master or a Redis Slave
* It has a **PID**
* It is monitored by Redis Monitor **(REDMON)** 
* **REDMON** updates the statistically information about this **PROC** periodically



### Redis Service Instance

A logical representation of the service instance that encapsulates one or more Redis **Proc** 

A Service Instance can be of the following type
* **Single Instance**:	Contains one redis **Proc** exposes an IP/Port
* **Master-Slave**:	Contains one `Proc` as master and rest of the redis `Proc`s as slaves
* **Cluster**: `Future Work when support of Redis 3.0 Cluster is added`

This is identical to **POD** terminology in K8s, we could group one or more **Proc**s as one unit, they are created and monitored together. 

***PS:*** *For convenience we will loose the obvious prefix 'Redis' and simple call `Service Instance` and `Proc` in the rest of the project*

### Considerations

It has been decided to use `etcd` as data store backend initially.  More support of other DB to be added later.

<img src="common/Store.jpg" width="100%" height="100%">
