##A Sample redis-proxy config and Run

Create an instnace like below using the mrr cli
```
$mrr create -n TestInstance -m 200 -s 2
Attempting to Creating a Redis Instance (TestInstance) with mem=200 slaves=2
Instance Creation accepted..
Check $mrr status -n TestInstance for status
```
Check the status like below
```
$mrr status -n TestInstance
Status = RUNNING
Type = MS
Capacity = 200
Master = 10.11.12.125:6380
        Slave0 = 10.11.12.123:6380
        Slave1 = 10.11.12.123:6381
```

To build the proxy its relatively simple, its a plain go program with no external dependencies.

```
$go build redis_proxy.go
```

It takes a json config file like below FROM ip:port to TO ip:port pair, the below one would actually mean bind to port 6677 the current system and proxy all the tcp to 10.11.12.125:6380 (which is a redis master)

```
$cat TestInstance_proxy.json
	{
			"HTTPPort": "7979",
			"Entries": [{
					"Name": "Master",
					"Pair": {
							"From": "0.0.0.0:6677",
							"To": "10.11.12.125:6380"
					}
			}]
	}

```
Now start the proxy
```
$./redis_proxy --config ./TestInstance_proxy.json
2016/07/06 00:55:03 The config file name is ./TestInstance_proxy.json
2016/07/06 00:55:03 Configuration file is = {7979 [{Master {0.0.0.0:6677 10.11.12.125:6380}}]}
2016/07/06 00:55:03 HandleConnection() {Master {0.0.0.0:6677 10.11.12.125:6380}}
```

Lets try to connec the instance via proxy. Note the redis-server itself is running @ a remote server only the proxy is running @ localhost.
```
redis-cli -h localhost -p 6677
localhost:6677> set foo bar
OK
localhost:6677> get foo
"bar"
localhost:6677> exit
```

if the master dies a new slave is promoted as a master now, lets verify that via mrr cli 

```
$mrr status -n TestInstance
Status = RUNNING
Type = MS
Capacity = 200
Master = 10.11.12.123:6380
        Slave0 = 10.11.12.123:6381
        Slave1 = 10.11.12.125:6380
```

Now lets update the proxy about the new master. get the current configuration of the proxy
```
$curl http://localhost:7979/Get/
Current Config: {
   "Master": {
     "Name": "Master",
     "Pair": {
       "From": "0.0.0.0:6677",
       "To": "10.11.12.125:6380"
     }
   }
 }
```

Update it via http rest as the new master is now at 10.11.12.123:6380

```
$curl http://localhost:7979/Update/ -X "PUT" -d '{"Name":"Master", "Addr":"10.11.12.123:6380"}'
```

re-verify the proxy's config file
```
curl http://localhost:7979/Get/
Current Config: {
   "Master": {
     "Name": "Master",
     "Pair": {
       "From": "0.0.0.0:6677",
       "To": "10.11.12.123:6380"
     }
   }
 }
```

Now re-connect to the same redis-server endpoint and check if the message is available.
```
~/redis_src/redis-stable/src/redis-cli -h localhost -p 6677
localhost:6677> get foo
"bar"
localhost:6677>
```
