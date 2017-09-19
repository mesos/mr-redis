package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/curator-go/curator"
	"github.com/samuel/go-zookeeper/zk"
	"github.com/hhkbp2/go-logging"
	"io"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)


var (
	//ConfigMap A map of name of the proxy vs its actually backend endpoint
	ConfigMap map[string]Entry

	//LocalPortsMap A map of local ports which fetch the data from zk once proxy daemon restarts
	LocalPortsMap map[string]string

)
const (

	//RedisPortBaseNum Local redis listen port range from 6100
	RedisPortMinNum = 6100

	RedisPortMaxNum = 6400

	ProxyPort = 7979

	CleanUpInterval = 3

	CleanUpZKMaxReties = 3

	CleanUpZKCheckIntervalSecs = 15

	SyncZKIntervalSecs = 3

	RedisPath = "/MrRedis/Instances"

	RedisLocalPortsPath = "/MrRedisLocalPorts"

	LogFilePath = "/data/apps/log/MrRedis-local-proxy.log"

	LogFileMaxSize = 100 * 1024 * 1024   // megabytes

	LogFileMaxBackups = 10

	LogFileMaxAge = 7    //days

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

/*func random(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}
*/


var logger = logging.GetLogger("redis_proxy")



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

	must(err)

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

func DeleteZKPathRecursive(path string) {

	zksStr := os.Getenv("ZOOKEEPER_SERVERS")

	if zksStr != "" {
		retryPolicy := curator.NewExponentialBackoffRetry(time.Second, CleanUpZKMaxReties, CleanUpZKCheckIntervalSecs*time.Second)
		client := curator.NewClient(zksStr, retryPolicy)
		client.Start()
		client.Delete().DeletingChildrenIfNeeded().ForPath(path)
		logger.Infof("deleteZKPathRecursive: remove zk znode %s recursively.", path)

		defer client.Close()

	} else {

		logger.Error("deleteZKPathRecursive: failed to get env variable ZOOKEEPER_SERVERS.")

	}
}

func InitializeProxy(conn *zk.Conn, path string) {


	redis_instance, _, err := conn.Children(path)

	if err != nil {
		logger.Error("Failed to load all redis instances from zk mr-redis path /MrRedis/Instances .")
		panic(err)
	}

	for _, name := range redis_instance {

		redis_status, _, _ := conn.Get(RedisPath + "/" + name + "/Status")


		if redis_status != nil && strings.EqualFold(string(redis_status), "RUNNING") {

			logger.Infof("redis instance %s status is running.", name)

			redis_id, _, err := conn.Get(RedisPath + "/" + name + "/Mname")

			must(err)

			redis_ip, _, err := conn.Get(RedisPath + "/" + name + "/Procs/" + string(redis_id) + "/IP")

			must(err)

			redis_port, _, err := conn.Get(RedisPath + "/" + name + "/Procs/" + string(redis_id) + "/Port")

			must(err)

			var redis_tcp_local_port string

			if CurrentE, ok := ConfigMap[name]; ok {

				logger.Infof("Redis instance %s is in the configMap. \n", name)

				if found, _, _ := conn.Exists(RedisLocalPortsPath + "/" + name); found {

					redis_port_byte, _, _ := conn.Get(RedisLocalPortsPath + "/" + name)

					redis_tcp_local_port = string(redis_port_byte[:])

					logger.Infof("InitializeProxy: Redis %s local port %s is already in the MrRedisLocalPort, sync with zk to keep it consistent . \n", name, redis_tcp_local_port)

					CurrentE.Pair.From = "127.0.0.1" + ":" + redis_tcp_local_port

					logger.Infof("Set redis instance %s Pair.From properties to %s" , name, CurrentE.Pair.From)

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

					var redis_port_found bool = false

					for {
						random_port := RandInt64(RedisPortMinNum, RedisPortMaxNum)

						redis_tcp_local_port = strconv.Itoa(random_port)

						logger.Infof("redis %s generate random local_ ort number is %s \n", name, redis_tcp_local_port)

						local_port_num := len(LocalPortsMap)

						logger.Infof("redis %s local port num is %d \n", name, local_port_num)

						if local_port_num > 0 {
							for _, value := range LocalPortsMap {
								if strings.EqualFold(redis_tcp_local_port, value) {
									redis_port_found = true
									logger.Infof("Redis %sredis port %s is already assigned.\n", name, value)
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
		go func(E Entry, F net.Conn) {

			T, err := net.Dial("tcp", E.Pair.To)
			if err != nil {
				logger.Errorf("Unable to connect to the Destination %s %v", E.Pair.To, err)
				return
			}
			defer T.Close()
			defer F.Close()

			go io.Copy(F, T)
			io.Copy(T, F)

		}(CurrentE, conn)
	}
}

func cleanProxy(conn *zk.Conn) {

	go func(){
		for {

			time.Sleep(time.Second * CleanUpInterval)

			logger.Infof("cleanProxy: Sleep %d seconds", CleanUpInterval)

			redis_instances, _, err := conn.Children(RedisPath)

			if err != nil {
				logger.Errorf("Failed to get redis instances.")
				return
			}

			logger.Infof("cleanProxy: redis_instaces nodes are %v", redis_instances)

			for _,name := range redis_instances {

				redis_status,_, err := conn.Get(RedisPath + "/" + name + "/" + "Status")

				logger.Infof("cleanProxy: redis %s status is %s", name, redis_status)

				if err != nil {
					logger.Errorf("cleanProxy: err occured when getting redis %s path %v", name, err)
					return
				}


				logger.Infof("cleanProxy:redis %s status is %s.\n", name, redis_status)

				if strings.EqualFold(string(redis_status), "DELETED") || redis_status == nil {


					//delete znode in zk

					DeleteZKPathRecursive(RedisPath + "/" + name)

					CurrentE, ok := ConfigMap[name]

					if ok {
						logger.Infof("cleanProxy: redis %s is in the ConfigMap", name)

						from_addr := CurrentE.Pair.From

						logger.Infof("cleanProxy: redis %s from_addr is %s", name, from_addr)

						delete(ConfigMap, name)

					} else {

						logger.Infof("cleanProxy: redis %s is not in the ConfigMap", name)
					}

				}
			}

		}
	}()
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
	fmt.Fprintf(w, string(retBytes) )
	return
}

func main() {

	var Cfg Config

	//Initialize the global Config map
	ConfigMap = make(map[string]Entry)

	//Initialize the global LocalPorts map
	LocalPortsMap = make(map[string]string)

	//Read a config file that has json update the config files
	cfgFileName := flag.String("config", "./config.json", "Supply the location of MrRedis configuration file")
	flag.Parse()

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


	logger.SetLevel(logging.LevelInfo)

	logger.AddHandler(handler)
	

	// ensure all log messages are flushed to disk before program exits.
	defer logging.Shutdown()

	logger.Infof("The config file name is %s ", *cfgFileName)
	cfgFile, err := ioutil.ReadFile(*cfgFileName)

	if err != nil {
		logger.Error("Error Reading the configration file. Resorting to default values")
	}
	err = json.Unmarshal(cfgFile, &Cfg)
	if err != nil {
		logger.Errorf("Error parsing the config file %v", err)
		return
	}
	logger.Infof("Configuration file is = %v", Cfg)

	conn := connect()

	defer conn.Close()

	//Initialize zk node /MrRedis-local-ports

	PrepareLocalPorts(conn, "/MrRedisLocalPorts")

	//Initialize existent proxy instance inside zk and added them into ConfigMap

	InitializeProxy(conn, RedisPath)

	//Clean up unused tcp ports, eg: when redis status is DELETED, the local proxy server of that redis will be shutdown.
	//comment cleanProxy to avoid concurrent write Map error
	//cleanProxy(conn)

	go func() {

		for {

			time.Sleep(time.Second * SyncZKIntervalSecs)

			logger.Info("Routine: Sync redis infomration from zk...")

			InitializeProxy(conn, RedisPath)
		}

	}()

	http.HandleFunc("/Update/", HandleHTTPUpdate)

	http.HandleFunc("/Get/", HandleHTTPGet)

	logger.Fatal(http.ListenAndServe(":"+Cfg.HTTPPort, nil))

	//Wait indefinitely
	waitCh := make(chan bool)

	<-waitCh

}
