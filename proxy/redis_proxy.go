package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"github.com/curator-go/curator"
	"github.com/natefinch/lumberjack"
	"github.com/samuel/go-zookeeper/zk"
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

	RedisPortMaxNum = 6300

	ProxyPort = 7979

	CleanUpInterval = 20

	CleanUpZKMaxReties = 3

	CleanUpZKCheckIntervalSecs = 15

	SyncZKIntervalSecs = 2

	RedisPath = "/MrRedis/Instances"

	RedisLocalPortsPath = "/MrRedisLocalPorts"

)

//Config json config structure for the proxy
type Config struct {
	HTTPPort string  //HTTPPort server Port number that we should bind to
	Entries  []Entry //Entries List of proxy entries
}

//Entry Representation of each entry in the proxy config
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
		log.Println("panic")
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

func RandInt64(min, max int) int {
	if min >= max || min == 0 || max == 0 {
		return max
	}
	return rand.Intn(max-min) + min
}

func PrepareLocalPorts(conn *zk.Conn, path string) {
	log.Printf("Begin to prepare redis_local_ports")
	found, _, err := conn.Exists(path)
	must(err)
	if found {
		log.Println(path + " already exist.")
	} else {
		log.Println(path + " doesn't exist, need to create it.")
		flags := int32(0)
		acl := zk.WorldACL(zk.PermAll)

		_, err := conn.Create(path, []byte("Mesos_local_ports_parent"), flags, acl)
		if err != nil {
			log.Println("Failed to create parent node " + path)
		}
	}

	redis_local_ports, _, err := conn.Children(path)

	must(err)

	for _, name := range redis_local_ports {

		local_port, _, _ := conn.Get(path + "/" + name)

		_, ok := LocalPortsMap[name]

		if ok {
			log.Printf("%s local port %s already exist in LocalPortsMap.\n", name, local_port)
		} else {
			LocalPortsMap[name] = string(local_port)
		}

	}

	log.Println("LocalPortsMap is")
	log.Println(LocalPortsMap)

}

func DeleteZKPathRecursive(path string) {
	zksStr := os.Getenv("ZOOKEEPER_SERVERS")
	//zks := strings.Split(zksStr, ",")

	if zksStr != "" {

		retryPolicy := curator.NewExponentialBackoffRetry(time.Second, CleanUpZKMaxReties, CleanUpZKCheckIntervalSecs*time.Second)
		client := curator.NewClient(zksStr, retryPolicy)
		client.Start()
		client.Delete().DeletingChildrenIfNeeded().ForPath(path)
		log.Printf("deleteZKPathRecursive: remove zk znode %s recursively.", path)

		defer client.Close()

	} else {

		log.Printf("deleteZKPathRecursive: failed to get env variable ZOOKEEPER_SERVERS.")

	}
}

func InitializeProxy(conn *zk.Conn, path string) {


	redis_instance, _, err := conn.Children(path)

	must(err)

	//time.Sleep(time.Second * 15)

	log.Println("InitializeProxy: Begin to initialize ans sync proxy servers from zk.")

	for _, name := range redis_instance {

		redis_status, _, _ := conn.Get(RedisPath + "/" + name + "/Status")


		if redis_status != nil && strings.EqualFold(string(redis_status), "RUNNING") {

			// var CurrentE Entry

			//var CurrentE Entry

			redis_id, _, err := conn.Get(RedisPath + "/" + name + "/Mname")

			must(err)

			redis_ip, _, err := conn.Get(RedisPath + "/" + name + "/Procs/" + string(redis_id) + "/IP")

			must(err)

			redis_port, _, err := conn.Get(RedisPath + "/" + name + "/Procs/" + string(redis_id) + "/Port")

			must(err)

			if CurrentE, ok := ConfigMap[name]; ok {

				log.Printf("InitializeProxy: Redis name %s is already in the configMap, only change its backend redis addr. \n", name)

				CurrentE.Pair.To = string(redis_ip) + ":" + string(redis_port)

				ConfigMap[name] = CurrentE

			} else {

				log.Printf("InitializeProxy: Redis name %s not found in the configMap \n", name)

				var redis_tcp_local_port string

				if found, _, _ := conn.Exists(RedisLocalPortsPath + "/" + name); found {

					redis_port_byte, _, _ := conn.Get(RedisLocalPortsPath + "/" + name)

					redis_tcp_local_port = string(redis_port_byte[:])

					log.Printf("redis instance %s exist, redis tcp local port is %s \n", name, redis_tcp_local_port)

				} else {

					//redis_listen_port := RedisPortBaseNum + len(redis_instance)

					var redis_port_found bool = false

					for {
						random_port := RandInt64(RedisPortMinNum, RedisPortMaxNum)

						redis_tcp_local_port = strconv.Itoa(random_port)

						log.Printf("redis %s generate random local_ ort number is %s \n", name, redis_tcp_local_port)

						local_port_num := len(LocalPortsMap)

						log.Printf("local port num is %d \n", local_port_num)

						if local_port_num > 0 {
							for _, value := range LocalPortsMap {
								if strings.EqualFold(redis_tcp_local_port, value) {
									redis_port_found = true
									log.Printf("Redis %sredis port %s is already assigned.\n", name, value)
									break
								}
							}

							if redis_port_found {
								log.Printf("Local tcp port %s is duplicated, will generate a new one.\n", redis_tcp_local_port)
								continue
							} else {
								log.Printf("random_tcp_port not assigned in local, so it can be used, will skip this loop.")
								break
							}
						} else {
							log.Println("LocalPortsMap length is zero, so a random port can be choosen")
							break
						}

						log.Printf("loop redis %s to check local port over\n", name)

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

				log.Printf("InitializeProxy: Redis %s local_tcp_addr is %s, to_tcp_addr is %s \n", name, local_tcp_addr.String(), to_tcp_addr.String())

				currentProxyPair := PorxyPair{From: local_tcp_addr.String(), To: to_tcp_addr.String()}

				CurrentEntry := Entry{Name: name, Pair: currentProxyPair}

				ConfigMap[name] = CurrentEntry

				go HandleConnection(CurrentEntry)

				log.Println("InitializeProxy: End of InitializeProxy")
			}
		}
	}


}

//HandleConnection Actuall proxy implementation per client. Untimatly this performs a implments a duplex io.Copy
func HandleConnection(E Entry) error {

	var CurrentE Entry //A Temp variable to get the latest Desination proxy value
	var OK bool

	log.Printf("HandleConnection() %v", E)
	//src, err := net.Listen("tcp", E.Pair.From)
	listener, err := newTCPListener(E.Pair.From)

	if err != nil {
		log.Printf("Error binding to the IP %v", err)
		return err
	}

	defer listener.Close()

	//Add this in the global Map so that it can be updated dynamically by HTTP apis
	ConfigMap[E.Name] = E

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting a new connection %v", err)
			continue
		}

		//Get the latest Entry from the MAP because it migh thave been updated on the fly.
		if CurrentE, OK = ConfigMap[E.Name]; !OK {
			log.Printf("Error Proxy entry is incorrect / empty for %s", E.Name)
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
				log.Printf("Unable to connect to the Destination %s %v", E.Pair.To, err)
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

			log.Printf("cleanProxy: Sleep %d seconds", CleanUpInterval)

			redis_instances, _, err := conn.Children(RedisPath)

			log.Printf("cleanProxy: redis_instaces nodes are %v", redis_instances)

			for _,name := range redis_instances {

				redis_status,_, err := conn.Get(RedisPath + "/" + name + "/" + "Status")

				log.Printf("cleanProxy: redis %s status is %s", name, redis_status)

				if err != nil {
					log.Printf("cleanProxy: err occured when getting redis %s path %v", name, err)
					return
				}


				log.Printf("cleanProxy:redis %s status is %s.\n", name, redis_status)

				if strings.EqualFold(string(redis_status), "DELETED") || redis_status == nil {


					//delete znode in zk

					DeleteZKPathRecursive(RedisPath + "/" + name)

					CurrentE, ok := ConfigMap[name]

					if ok {
						log.Printf("cleanProxy: redis %s is in the ConfigMap", name)

						from_addr := CurrentE.Pair.From

						log.Printf("cleanProxy: redis %s from_addr is %s", name, from_addr)

						delete(ConfigMap, name)

					} else {

						log.Printf("cleanProxy: redis %s is not in the ConfigMap", name)
					}

				}
			}
			if err != nil {
				return
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
			log.Printf("Error understanding the Body %v", err)
			return
		}

		var val HTTPUpdate
		var CurrentE Entry
		var OK bool
		err = json.Unmarshal(content, &val)
		if err != nil {
			//log.Printf(w, "Wrong json format %v", err)
			log.Printf("Wrong json format %v", err)
			return
		}
		if CurrentE, OK = ConfigMap[val.Name]; !OK {
			log.Printf("Error Proxy entry is incorrect / empty for %s", val.Name)
			//log.Printf(w, "Error Proxy entry is incorrect / empty for %s", val.Name)
			return
		}
		log.Printf("Updating From porxy for %s From %s TO %s", val.Name, CurrentE.Pair.To, val.Addr)
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
		log.Printf("Error Marshalling HandleHTTPGet() %v", err)
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

        //set log rotating policy

       	log.SetOutput(&lumberjack.Logger{
           Filename:   "/data/apps/log/MrRedis-local-proxy.log",
   	   MaxSize:    50, // megabytes
           MaxBackups: 10,
           MaxAge:     3, //days
})

	//Read a config file that has json update the config files
	cfgFileName := flag.String("config", "./config.json", "Supply the location of MrRedis configuration file")
	flag.Parse()

	log.Printf("The config file name is %s ", *cfgFileName)
	cfgFile, err := ioutil.ReadFile(*cfgFileName)

	if err != nil {
		log.Printf("Error Reading the configration file. Resorting to default values")
	}
	err = json.Unmarshal(cfgFile, &Cfg)
	if err != nil {
		log.Fatalf("Error parsing the config file %v", err)
		return
	}
	log.Printf("Configuration file is = %v", Cfg)

	conn := connect()

	defer conn.Close()

	//Initialize zk node /MrRedis-local-ports

	PrepareLocalPorts(conn, "/MrRedisLocalPorts")

	//Initialize existent proxy instance inside zk and added them into ConfigMap

	InitializeProxy(conn, RedisPath)

	//Clean up unused tcp ports, eg: when redis status is DELETED, the local proxy server of that redis will be shutdown.

	conn1 := connect()

	defer conn1.Close()

	cleanProxy(conn1)

	go func() {

		for {

			time.Sleep(time.Second * SyncZKIntervalSecs)

			log.Printf("Routine: Sync redis infomration from zk...")

			InitializeProxy(conn, RedisPath)
		}

	}()

	http.HandleFunc("/Update/", HandleHTTPUpdate)

	http.HandleFunc("/Get/", HandleHTTPGet)

	log.Fatal(http.ListenAndServe(":"+Cfg.HTTPPort, nil))

	//Wait indefinitely
	waitCh := make(chan bool)

	<-waitCh

}
