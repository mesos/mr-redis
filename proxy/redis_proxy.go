package main

import (
	"encoding/json"
	"fmt"
	"github.com/curator-go/curator"
	"github.com/curator-go/curator/recipes/cache"
	"github.com/hhkbp2/go-logging"
	"github.com/samuel/go-zookeeper/zk"
	"io"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	//ConfigMap A map of name of the proxy vs its actually backend endpoint
	ConfigMap map[string]Entry

	//LocalPortsMap A map of local ports which fetch the data from zk once proxy daemon restarts
	LocalPortsMap map[string]string

	//Define logger name of program as redis_proxy
	//logger = logging.GetLogger("redis_proxy")

	logger logging.Logger
)

var startTime = time.Now()

//Add WMutex to ConfigMap avoid concurrent read and write error
var lock = sync.RWMutex{}

var bufferPool = sync.Pool{
	New: func() interface{} {
		// TODO maybe different buffer size?
		// benchmark pls
		return make([]byte, 1<<15)
	},
}

const (

	//RedisPortBaseNum Local redis listen port range from 6100
	RedisPortMinNum = 6100

	RedisPortMaxNum = 6400

	ProxyAddr = "127.0.0.1:7979"

	SyncZKIntervalSecs = 3

	RedisPath = "/MrRedis/Instances"

	RedisLocalPortsPath = "/MrRedisLocalPorts"

	LogFilePath = "/data/apps/log/MrRedis-local-proxy.log"

	LogFileMaxSize = 100 * 1024 * 1024 // megabytes

	LogFileMaxBackups = 10

	ProgrameStartTimeAtLeast = 30

	FetchRedisIpTimeOutSecs = 60
)

//Config json config structure for the proxy
type Config struct {
	HTTPPort string  //HTTPPort server Port number that we should bind to
	Entries  []Entry //Entries List of proxy entries
}

//Entry Representation of each entry in the proxy configå
type Entry struct {
	Name string
	Pair PorxyPair
}

//PorxyPair The actuall proxy pair from (bind port) to actual port
type PorxyPair struct {
	From string //IP:PORT pair
	To   string //IP:PORT pair
}

//HTTPUpdate This structure is used by the HTTP PUT request to change the IP address of the destination on the fly
type HTTPUpdate struct {
	Name string
	Addr string
}

func must(err error) {
	if err != nil {
		logger.Error("panic")
		panic(err)
	}
}

func connect() *zk.Conn {
	zksStr := os.Getenv("ZOOKEEPER_SERVERS")
	zks := strings.Split(zksStr, ",")
	conn, _, err := zk.Connect(zks, time.Second)
	must(err)
	return conn
}

func newTCPListener(addr string) (net.Listener, error) {
	conn, err := net.Listen("tcp", addr)
	if err != nil {
		return conn, err
	}

	return conn, nil
}

func RandInt64(min, max int) int {
	if min >= max || min == 0 || max == 0 {
		return max
	}
	return rand.Intn(max-min) + min
}

func PrepareLocalPorts(conn *zk.Conn, path string) {

	logger.Info("Begin to prepare redis_local_ports")

	found, _, err := conn.Exists(path)

	must(err)

	if found {
		logger.Infof(path + " already exist.")
	} else {
		logger.Infof(path + " doesn't exist, need to create it.")
		flags := int32(0)
		acl := zk.WorldACL(zk.PermAll)

		_, err := conn.Create(path, []byte("Mesos_local_ports_parent"), flags, acl)
		if err != nil {
			logger.Warnf("Failed to create parent node " + path)
		}
	}

	redis_local_ports, _, err := conn.Children(path)
	if err != nil {
		logger.Errorf("Error to get redis_local_ports, error is %s.", err)
		return
	}

	for _, name := range redis_local_ports {

		local_port, _, _ := conn.Get(path + "/" + name)

		_, ok := LocalPortsMap[name]

		if ok {
			logger.Infof("%s local port %s already exist in LocalPortsMap.\n", name, local_port)
		} else {
			LocalPortsMap[name] = string(local_port)
		}

	}

}


func getRedisMnameInfo(name string, conn *zk.Conn) (string, string) {

	logger.Infof("Get redis %v Mname info redis_ip and redis_port.", name)

	redis_id_path := RedisPath + "/" + name + "/Mname"

	var redis_id string

	idTimeCount := time.Now()

	for {
		redisId, _, idErr := conn.Get(redis_id_path)

		if idErr != nil {
			logger.Errorf("zk path /name/instance/Mname error: %v\n", RedisPath+"/"+name+"/Mname")
			return "",""
		}

		if redisId != nil && string(redisId) != "" {

			logger.Infof("Redis %s get the id from zk, the redis id is %s", name, string(redis_id))
			redis_id = string(redisId)
			break

		} else {
			elapsed := time.Since(idTimeCount).Seconds()

			logger.Infof("Fetch redis %s spends %d seconds already.", name, elapsed)

			if elapsed > FetchRedisIpTimeOutSecs {
				logger.Errorf("Failed to fetch redis %s id, and it's over %d secoonds. will ignore this request!", name, FetchRedisIpTimeOutSecs)
				break
			}

			time.Sleep(1 * time.Second)

			logger.Errorf("Redis %s failed to get new redis id, the id is %s. Will get fetch it again.", name, string(redisId))
		}


	}

    if redis_id == "" {
    	logger.Errorf("Get redis %s Mname id null, will return empty string!")
    	return "",""
	}

	redis_ip_path := RedisPath + "/" + name + "/Procs/" + redis_id + "/IP"

	logger.Infof("redis %s redis_ip_path is %s", name, redis_id_path)

	var redis_ip string

	timeCount := time.Now()

	for {
		redisIp, _, err := conn.Get(redis_ip_path)

		if err == nil {

			if redisIp != nil && string(redisIp) != "" {

				logger.Infof("Redis %s get the ip from zk, the redis ip is %v", name, redis_ip)
				redis_ip = string(redisIp)
				break

			} else {

				logger.Errorf("Redis %v failed to get new redis ip, the ip is %v. Will get fetch it again.", name, redis_ip)
			}

			elapsed := time.Since(timeCount).Seconds()

			logger.Infof("Fetch redis %s spends %d seconds already.", name, elapsed)

			if elapsed > FetchRedisIpTimeOutSecs {
				logger.Errorf("Failed to fetch redis %s ip, and it's over %d secoonds. will ignore this request!", name, FetchRedisIpTimeOutSecs)
				break
			}

			time.Sleep(1 * time.Second)

		} else {
			logger.Errorf("failed to get redis ip, error is %s", err)
			logger.Error("Failed to get redis %s ip, redis ip is %v, zk conection might have problem, eth error is %v", name, redis_ip, err)
			break
		}
	}

	if redis_ip == "" {
		logger.Errorf("Get redis %s IP as null, will return empty string!")
		return "",""
	}

	redis_port_path := RedisPath + "/" + name + "/Procs/" + redis_id + "/Port"
	redis_port, _, err := conn.Get(redis_port_path)

	if err != nil {
		logger.Errorf("zk path name/Pros/instance/Port error: %v\n", RedisPath+"/"+name+"/Procs/"+string(redis_id)+"/Port")
		return "", ""
	}

	return string(redis_ip), string(redis_port)
}

func InitializeProxy(conn *zk.Conn, path string) {

	logger.Infof("Run InitializeProxy at boot time %v", time.Now())

	redis_instance, _, err := conn.Children(path)

	if err != nil {
		logger.Error("Failed to load all redis instances from zk mr-redis path /MrRedis/Instances .")
		panic(err)
	}

	for _, name := range redis_instance {

		redis_status, _, _ := conn.Get(RedisPath + "/" + name + "/Status")

		if redis_status != nil && strings.EqualFold(string(redis_status), "RUNNING") {

			logger.Infof("redis instance %s status is running.", name)

			redis_mname, _, _ :=  conn.Get(RedisPath + "/" + name + "/Mname")

			if redis_mname == nil || string(redis_mname) == "" {
				logger.Errorf("redis %s Mname is empty. Will skip this redis instance.", name)
				continue
			}

			redis_ip, redis_port := getRedisMnameInfo(name, conn)

			if redis_ip == "" || redis_port == "" {
				logger.Errorf("redis %s Pairto ip %s or port %s is empty. Will skip this redis instance.", name, redis_ip, redis_port)
				continue
			}

			var redis_tcp_local_port string

			if CurrentE, ok := ConfigMap[name]; ok {

				logger.Infof("Redis instance %s is in the configMap. \n", name)

				if found, _, _ := conn.Exists(RedisLocalPortsPath + "/" + name); found {

					redis_port_byte, _, _ := conn.Get(RedisLocalPortsPath + "/" + name)

					redis_tcp_local_port = string(redis_port_byte[:])

					logger.Infof("InitializeProxy: Redis %s local port %s is already in the MrRedisLocalPort, sync with zk to keep it consistent . \n", name, redis_tcp_local_port)

					CurrentE.Pair.From = "127.0.0.1" + ":" + redis_tcp_local_port

					logger.Infof("Set redis instance %s Pair.From properties to %s", name, CurrentE.Pair.From)

				}

				CurrentE.Pair.To = string(redis_ip) + ":" + string(redis_port)

				ConfigMap[name] = CurrentE

			} else {

				logger.Infof("Redis name %s not found in the configMap \n", name)

				if found, _, _ := conn.Exists(RedisLocalPortsPath + "/" + name); found {

					redis_port_byte, _, _ := conn.Get(RedisLocalPortsPath + "/" + name)

					redis_tcp_local_port = string(redis_port_byte[:])

					logger.Infof("redis instance %s already exists, redis tcp local port is %s \n", name, redis_tcp_local_port)

				} else {

					//redis_listen_port := RedisPortBaseNum + len(redis_instance)

					redis_tcp_local_port = getLocalRedisPort()

					//redis_tcp_listen_port := strconv.Itoa(random_tcp_port)
					flags := int32(0)

					acl := zk.WorldACL(zk.PermAll)

					conn.Create(RedisLocalPortsPath+"/"+name, []byte(redis_tcp_local_port), flags, acl)

				}

				local_addr := "127.0.0.1" + ":" + redis_tcp_local_port

				local_tcp_addr, _ := net.ResolveTCPAddr("tcp4", local_addr)

				to_addr := string(redis_ip) + ":" + string(redis_port)

				to_tcp_addr, _ := net.ResolveTCPAddr("tcp4", to_addr)

				logger.Infof("Redis %s local_tcp_addr is %s, to_tcp_addr is %s \n", name, local_tcp_addr.String(), to_tcp_addr.String())

				currentProxyPair := PorxyPair{From: local_tcp_addr.String(), To: to_tcp_addr.String()}

				CurrentEntry := Entry{Name: name, Pair: currentProxyPair}

				ConfigMap[name] = CurrentEntry

				go HandleConnection(CurrentEntry)

			}
		}
	}

}

func getLocalRedisPort() string {

	var redis_port_found bool = false

	var redis_tcp_local_port string

	for {
		random_port := RandInt64(RedisPortMinNum, RedisPortMaxNum)

		redis_tcp_local_port = strconv.Itoa(random_port)

		logger.Infof("redis generate random local_ ort number is %s \n", redis_tcp_local_port)

		local_port_num := len(LocalPortsMap)

		logger.Infof("redis local port num is %d \n", local_port_num)

		if local_port_num > 0 {
			for _, value := range LocalPortsMap {
				if strings.EqualFold(redis_tcp_local_port, value) {
					redis_port_found = true
					logger.Infof("redis port %s is already assigned.\n", value)
					break
				}
			}

			if redis_port_found {
				logger.Infof("Local tcp port %s is duplicated, will generate a new one.\n", redis_tcp_local_port)
				continue
			} else {
				logger.Info("random_tcp_port not assigned in local, so it can be used, will skip this loop.")
				break
			}
		} else {
			logger.Warn("LocalPortsMap length is zero, so a random port can be choosen")
			break
		}

	}

	return redis_tcp_local_port
}

//HandleConnection Actuall proxy implementation per client. Untimatly this performs a implments a duplex io.Copy
func HandleConnection(E Entry) error {

	var CurrentE Entry //A Temp variable to get the latest Desination proxy value
	var OK bool

	logger.Info("HandleConnection() %v", E)
	//src, err := net.Listen("tcp", E.Pair.From)
	listener, err := newTCPListener(E.Pair.From)

	if err != nil {
		logger.Errorf("Error binding to the IP %v", err)
		return err
	}

	defer listener.Close()

	//Add this in the global Map so that it can be updated dynamically by HTTP apis
	ConfigMap[E.Name] = E

	for {
		conn, err := listener.Accept()
		if err != nil {
			logger.Errorf("Error accepting a new connection %v", err)
			continue
		}

		//Get the latest Entry from the MAP because it migh thave been updated on the fly.
		if CurrentE, OK = ConfigMap[E.Name]; !OK {
			logger.Errorf("Error Proxy entry is incorrect / empty for %s", E.Name)
			conn.Close()
			continue
		}

		//Start a Lamda for performing the proxy
		//F := From Connection
		//T := To Connection
		//This proxy will simply transfer everything from F to T net.Conn
		go func(E Entry, srcConn net.Conn) {

			destConn, err := net.Dial("tcp", E.Pair.To)

			if err != nil {
				logger.Errorf("Unable to connect to the Destination %s %v", E.Pair.To, err)
				return
			}

			first := make(chan<- struct{}, 1)
			var wg sync.WaitGroup
			cp := func(dst net.Conn, src net.Conn) {
				buf := bufferPool.Get().([]byte)
				// TODO use splice on linux
				// TODO needs some timeout to prevent torshammer ddos
				_, err := io.CopyBuffer(dst, src, buf)
				select {
				case first <- struct{}{}:
					if err != nil {
						logger.Errorf("Copy error is %v:", err)
					}
					_ = dst.Close()
					_ = src.Close()
				default:
				}
				bufferPool.Put(buf)
				wg.Done()
			}
			wg.Add(2)
			go cp(destConn, srcConn)
			go cp(srcConn, destConn)
			wg.Wait()

		}(CurrentE, conn)
	}
}

//HandleHTTPUpdate Call beack to handle /Update/ HTTP call back
func HandleHTTPUpdate(w http.ResponseWriter, r *http.Request) {
	//log.Printf(w, "Hi there, Going to Update %s! Method=%s\n", r.URL.Path[1:], r.Method)
	if r.Method == "PUT" {
		//This can be used for updating an existing variable
		content, err := ioutil.ReadAll(r.Body)
		r.Body.Close()
		if err != nil {
			//log.Printf(w, "Error understanding the Body %v", err)
			logger.Errorf("Error understanding the Body %v", err)
			return
		}

		var val HTTPUpdate
		var CurrentE Entry
		var OK bool
		err = json.Unmarshal(content, &val)
		if err != nil {
			//log.Printf(w, "Wrong json format %v", err)
			logger.Errorf("Wrong json format %v", err)
			return
		}
		if CurrentE, OK = ConfigMap[val.Name]; !OK {
			logger.Infof("Error Proxy entry is incorrect / empty for %s", val.Name)
			//log.Printf(w, "Error Proxy entry is incorrect / empty for %s", val.Name)
			return
		}
		logger.Info("Updating From porxy for %s From %s TO %s", val.Name, CurrentE.Pair.To, val.Addr)
		CurrentE.Pair.To = val.Addr
		ConfigMap[val.Name] = CurrentE
		return
	}
	return
}

//HandleHTTPGet call back to handle /Get/ HTTP callback
func HandleHTTPGet(w http.ResponseWriter, r *http.Request) {
	retBytes, err := json.MarshalIndent(ConfigMap, " ", "  ")
	if err != nil {
		logger.Errorf("Error Marshalling HandleHTTPGet() %v", err)
		//log.Printf(w, "Error Marshalling HandleHTTPGet() %v", err)
		return

	}
	fmt.Fprintf(w, string(retBytes))
	return
}

func addRedisProxy(name string, conn *zk.Conn) {

	var CurrentE Entry
	var OK bool

	if name == "" {
		logger.Errorf("Redis name is empty, will ingore this request.")
		return
	}
	if CurrentE, OK = ConfigMap[name]; OK {

		logger.Infof("Redis instance %v proxy already exist in configMap.", name)
		return

	} else {

		logger.Infof("Redis instance not exsit in configMap.", name)

		redis_ip, redis_port := getRedisMnameInfo(name, conn)

		//Add lock to ConfigMap in case of concurrent read and write on configMap. eg: create redis and existant redis failover happens at the same time, this might occur

		if redis_ip == "" || redis_port == "" {
			logger.Errorf("Failed to add redis instance %s, eigher redis_ip or redis_port values is empty. redis_ip is %v, redis_port is %v", name, redis_ip, redis_port)
			return
		}

		lock.Lock()

		defer lock.Unlock()

		CurrentE.Pair.To = redis_ip + ":" + redis_port

		redis_tcp_local_port := getLocalRedisPort()

		//redis_tcp_listen_port := strconv.Itoa(random_tcp_port)
		flags := int32(0)

		acl := zk.WorldACL(zk.PermAll)

		conn.Create(RedisLocalPortsPath+"/"+name, []byte(redis_tcp_local_port), flags, acl)

		CurrentE.Pair.From = "127.0.0.1" + ":" + redis_tcp_local_port

		ConfigMap[name] = CurrentE

		go HandleConnection(CurrentE)

	}
}

func updateRedisProxy(name string, conn *zk.Conn) {

	var CurrentE Entry
	var OK bool

	if CurrentE, OK = ConfigMap[name]; OK {

		logger.Infof("Redis %s exist in ConfigMap, and it might have failoevr occurred, will master ip.", name)

		redis_ip, redis_port := getRedisMnameInfo(name, conn)

		if redis_ip == "" || redis_port == "" {
			logger.Errorf("Failed to update redis, eigher redis_ip or redis_port values is empty. redis_ip is %v, redis_port is %v", redis_ip, redis_port)
			return
		}
		//Add lock to ConfigMap in case of concurrent read and write on configMap. eg: create redis and existant redis failover happens at the same time, this might occur
		lock.Lock()

		defer lock.Unlock()

		CurrentE.Pair.To = redis_ip + ":" + redis_port
		ConfigMap[name] = CurrentE
		logging.Warnf("Change Redis %v master address into %v", name, CurrentE.Pair.To)

		return

	} else {
		logger.Warnf("Redis %s not exit in ConfigMap, will return", name)
		return
	}

}

func watchRedisStatus(conn *zk.Conn) {

	zksStr := os.Getenv("ZOOKEEPER_SERVERS")

	retryPolicy := curator.NewExponentialBackoffRetry(time.Second, 3, 15*time.Second)

	client := curator.NewClient(zksStr, retryPolicy)

	client.Start()

	defer client.Close()

	treeCache := cache.NewTreeCache(client, RedisPath, cache.DefaultTreeCacheSelector)

	treeCache.Start()

	defer treeCache.Stop()

	cacheListener := cache.NewTreeCacheListener(func(client curator.CuratorFramework, event cache.TreeCacheEvent) error {

		switch event.Type.String() {

		case "NodeAdded":
			//fmt.Printf( event_path)
			event_path := event.Data.Path()
			logger.Infof("TreeCache event is: NodeAdded, zk path is %s \n", event_path)

			if strings.Contains(event_path, "Mname") {

				elapsed := time.Since(startTime).Seconds()

				if elapsed < ProgrameStartTimeAtLeast {

					logger.Infof("Program is just started in, will skip the InitializePorxy function ")

				} else {

					logger.Infof("New redis has been created, Will Sync the status")

					time.Sleep(2 * SyncZKIntervalSecs * time.Second)

					redisName := strings.Split(event_path, "/")[3]

					if redisName != "" {

						if _, ok := ConfigMap[redisName]; ok {

							logger.Infof("Redis %s has already been created!", redisName)

						} else {

							logger.Infof("Redis %s has not been created, will create it later.", redisName)
							//addRedisProxy(redisName, conn)

						}
					} else {
						logger.Errorf("Failed to get redis name ")
					}

				}
			}

		case "NodeUpdated":
			//	fmt.Printf( event_path)
			event_path := event.Data.Path()

			logger.Infof("TreeCache event is: NodeUpdated, zk path is %s \n", event_path)

			if strings.Contains(event_path, "Mname") {

				logger.Infof("Redis node instance has changed, will sync the updates to ConfigMap.")
				//time.Sleep(SyncZKIntervalSecs * time.Second)

				elapsed := time.Since(startTime).Seconds()

				if elapsed < ProgrameStartTimeAtLeast {

					logger.Infof("Program might be just started in very short time, will skip the InitializePorxy function ")

				} else {


					time.Sleep(SyncZKIntervalSecs * time.Second)

					redisName := strings.Split(event_path, "/")[3]

					if redisName != "" {

						if _, ok := ConfigMap[redisName]; ok {

							redis_status_path := RedisPath + "/" + redisName + "/Status"

							redis_status, _, err := conn.Get(redis_status_path)

							if err != nil {
								logger.Errorf("Failed to get redis %v status %v, error is %v", redisName, redis_status, err.Error())
							} else {

								logger.Infof("redis %s status is %v.", redisName, redis_status)
							}

							switch  string(redis_status) {

							case "RUNNING":
								logger.Infof("Redis %s status is %v, failover might have occurred, will try to update the master ip by running updateRedisProxy.!", redisName, redis_status)
								updateRedisProxy(redisName, conn)
							default:
								logger.Infof("Redis %s status is %s, failover might have occurred, or redis is deleted!", redisName, redis_status)

							}

						} else {

							logger.Infof("Redis %s is not in ConfigMap, will create it by running addRedisProxy function", redisName)
							addRedisProxy(redisName, conn)

						}

					} else {
						logger.Errorf("Failed to extract redis name from event_path %v, and redis name  %v is empty", event_path, redisName)
					}

				}

			}

			if strings.Contains(event_path, "/Status") {

				logger.Infof("Redis node instance status has changed, will sync the updates to ConfigMap.")

				redisName := strings.Split(event_path, "/")[3]

				if redisName != "" {

					if _, ok := ConfigMap[redisName]; ok {

						redis_status_path := RedisPath + "/" + redisName + "/Status"

						redis_status, _, err := conn.Get(redis_status_path)

						if err != nil {
							logger.Errorf("Failed to get redis %v status %v, error is %v", redisName, redis_status, err.Error())
						} else {

							logger.Infof("redis %s status is %v.", redisName, redis_status)
						}

						switch  string(redis_status) {

						case "DELETED":
							logger.Infof("Redis %v status is deleted, should remove it from configMap.")
							lock.Lock()
							defer lock.Unlock()
							delete(ConfigMap, redisName)
						default:
							logger.Infof("redis %s status is %s, will do nothing about it.", redisName, redis_status )
						}
					}
				}
			}

			logging.Infof("Last setp on UpdateNode, CLEAN empty key.")
			for key,_ := range ConfigMap {
				if key == "" {
					lock.Lock()
					defer lock.Unlock()
					delete(ConfigMap,"")
				}
			}



		case "NodeRemoved":
			//fmt.Printf( event_path)
			event_path := event.Data.Path()
			logger.Infof("TreeCache event is: NodeRemoved \n, zk path is %s", event_path)
		case "ConnSuspended":
			//fmt.Printf( event_path)
			logger.Infof("TreeCache event is: ConnSuspended \n")
		case "ConnReconnected":
			//fmt.Printf( event_path)
			logger.Infof("TreeCache event is: ConnReconnected \n")
		case "ConnLost":
			//fmt.Printf( event_path)
			logger.Infof("TreeCache event is: ConnLost \n")
		case "Initialized":
			//fmt.Printf( event_path)
			logger.Infof("TreeCache event is: Initialized \n")
		default:
			logger.Infof("TreeCache event is: unknown. \n")
		}
		return nil
	})

	treeCache.Listenable().AddListener(cacheListener)

	logger.Infof("Adding listener for treeCache.")

	wait_ch := make(chan bool)
	<-wait_ch
}

/*
func init() {

	filePath := LogFilePath

	fileMode := os.O_APPEND

	bufferSize := 0

	bufferFlushTime := 30 * time.Second

	inputChanSize := 1

	backupCount := uint32(LogFileMaxBackups)
	// set the maximum size of every file to 100 M bytes
	fileMaxBytes := uint64(LogFileMaxSize)

	// create a handler(which represents a log message destination)
	handler := logging.MustNewRotatingFileHandler(
		filePath, fileMode, bufferSize, bufferFlushTime, inputChanSize,
		fileMaxBytes, backupCount)


	// the format for the whole log message
	format := "%(asctime)s %(levelname)s (%(filename)s:%(lineno)d) " +
		"%(name)s %(message)s"

	// the format for the time part
	dateFormat := "%Y-%m-%d %H:%M:%S.%3n"

	// create a formatter(which controls how log messages are formatted)
	formatter := logging.NewStandardFormatter(format, dateFormat)

	// set formatter for handler
	handler.SetFormatter(formatter)

	//Define logger name of program as redis_proxy
	logger = logging.GetLogger("redis_proxy")

	logger.SetLevel(logging.LevelInfo)

	logger.AddHandler(handler)

	// ensure all log messages are flushed to disk before program exits.
	defer logging.Shutdown()

	fmt.Println("Finish init")

}
*/
func main() {

	//Initialize the global Config map
	ConfigMap = make(map[string]Entry)

	//Initialize the global LocalPorts map
	LocalPortsMap = make(map[string]string)

	//init()
	filePath := LogFilePath

	fileMode := os.O_APPEND

	bufferSize := 0

	bufferFlushTime := 30 * time.Second

	inputChanSize := 1

	backupCount := uint32(LogFileMaxBackups)
	// set the maximum size of every file to 100 M bytes
	fileMaxBytes := uint64(LogFileMaxSize)

	// create a handler(which represents a log message destination)
	handler := logging.MustNewRotatingFileHandler(
		filePath, fileMode, bufferSize, bufferFlushTime, inputChanSize,
		fileMaxBytes, backupCount)

	// the format for the whole log message
	format := "%(asctime)s %(levelname)s (%(filename)s:%(lineno)d) " +
		"%(name)s %(message)s"

	// the format for the time part
	dateFormat := "%Y-%m-%d %H:%M:%S.%3n"

	// create a formatter(which controls how log messages are formatted)
	formatter := logging.NewStandardFormatter(format, dateFormat)

	// set formatter for handler
	handler.SetFormatter(formatter)

	//Define logger name of program as redis_proxy
	logger = logging.GetLogger("redis_proxy")

	logger.SetLevel(logging.LevelInfo)

	logger.AddHandler(handler)

	// ensure all log messages are flushed to disk before program exits.
	defer logging.Shutdown()

	conn := connect()

	defer conn.Close()

	//Initialize zk node /MrRedis-local-ports

	PrepareLocalPorts(conn, "/MrRedisLocalPorts")

	//Initialize existent proxy instance inside zk and added them into ConfigMap

	go InitializeProxy(conn, RedisPath)

	//Watch each redis status and take action if failover occurs or new redis created
	go watchRedisStatus(conn)

	http.HandleFunc("/Update/", HandleHTTPUpdate)

	http.HandleFunc("/Get/", HandleHTTPGet)

	err := http.ListenAndServe(ProxyAddr, nil)

	if err != nil {

		logger.Errorf("Failed to start http server on port %v!, error is %v", ProxyAddr, err.Error())

	} else {

		logger.Infof("Start http server on port %v successfuly!", ProxyAddr)

	}

	//Wait indefinitely
	waitCh := make(chan bool)

	<-waitCh

}
