package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/mesos/mr-redis/common/types"
	"github.com/mesos/mr-redis/sched/cmd"
	"github.com/mesos/mr-redis/sched/httplib"
	"github.com/mesos/mr-redis/sched/mesoslib"
)

//Declare all the Constants to be used in this file

const (
	HTTP_SERVER_PORT = "8080"
)

//Main execution point

type MrRedisConfig struct {
	Master       string //MesosMaster's endpoint zk://mesos.master/2181 or 10.11.12.13:5050
	ExecutorPath string //Executor's Path from where to distribute
	RedisPath    string //Path where redis-server executable is available
	DBType       string //Type of the database etcd/zk/mysql/consul etcd.,
	DBEndPoint   string //Endpoint of the database
	LogFile      string //Name of the logfile
	ArtifactIP   string //The IP to which we should bind to for distributing the executor among the interfaces
	ArtifactPort string //The port to which we shoudl bind to for distributing the executor
	HTTPPort     string //Defaults to 8080 if otherwise specify explicitly
}

func NewMrRedisDefaultConfig() MrRedisConfig {
	return MrRedisConfig{
		Master:       "127.0.0.1:5050",
		ExecutorPath: "./MrRedisExecutor",
		RedisPath:    "./redis-server",
		DBType:       "etcd",
		DBEndPoint:   "127.0.0.1:2379",
		LogFile:      "stderr",
		ArtifactIP:   "127.0.0.1",
		ArtifactPort: "5454",
		HTTPPort:     "5656",
	}
}

func main() {

	cfg_file_name := flag.String("config", "./config.json", "Supply the location of MrRedis configuration file")
	dumpConfig := flag.Bool("DumpEmptyConfig", false, "Dump Empty Config file")
	flag.Parse()

	Cfg := NewMrRedisDefaultConfig()

	if *dumpConfig == true {
		config_bytes, err := json.MarshalIndent(Cfg, " ", "  ")
		if err != nil {
			log.Printf("Error marshalling the dummy config file. Exiting %v", err)
			return
		}
		fmt.Printf("%s\n", string(config_bytes))
		return
	}

	cfg_file, err := ioutil.ReadFile(*cfg_file_name)

	if err != nil {
		log.Printf("Error Reading the configration file. Resorting to default values")
	}
	err = json.Unmarshal(cfg_file, &Cfg)
	if err != nil {
		log.Fatalf("Error parsing the config file %v", err)
	}
	log.Printf("Configuration file is = %v", Cfg)

	log.Printf("*****************************************************************")
	log.Printf("*********************Starting MrRedis-Scheduler******************")
	log.Printf("*****************************************************************")
	//Command line argument parsing

	//Facility to overwrite the etcd endpoint for scheduler if its running in the same docker container and expose a different one for executors

	db_endpoint := os.Getenv("ETCD_LOCAL_ENDPOINT")

	if db_endpoint == "" {
		db_endpoint = Cfg.DBEndPoint
	}

	//Initalize the common entities like store, store configuration etc.
	isInit, err := types.Initialize(Cfg.DBType, db_endpoint)
	if err != nil || isInit != true {
		log.Fatalf("Failed to intialize Error:%v return %v", err, isInit)
	}

	//Start the Mesos library
	go mesoslib.Run(Cfg.Master, Cfg.ArtifactIP, Cfg.ArtifactPort, Cfg.ExecutorPath, Cfg.RedisPath, Cfg.DBType, Cfg.DBEndPoint)

	//Start the creator
	go cmd.Creator()

	//Start the Mainterainer
	go cmd.Maintainer()

	//Start the Destroyer
	go cmd.Destoryer()

	//Start HTTP server and related things to handle restfull calls to the scheduler
	httplib.Run(Cfg.HTTPPort)

	//Wait for termination signal

	log.Printf("*****************************************************************")
	log.Printf("*********************Finished MrRedis-Scheduler******************")
	log.Printf("*****************************************************************")

}
