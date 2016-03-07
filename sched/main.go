package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"

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
	MasterIP     string //MesosMaster's IP address
	MasterPort   string //Mesos Masters Port number
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
		MasterPort:   "5050",
		MasterIP:     "127.0.0.1",
		ExecutorPath: "./MrRedisExecutor",
		RedisPath:    "./redis-server",
		DBType:       "etcd",
		DBEndPoint:   "127.0.0.1:2379",
		LogFile:      "stderr",
		ArtifactIP:   "127.0.0.1",
		ArtifactPort: "5454",
		HTTPPort:     "8080",
	}
}

func main() {

	cfg_file_name := flag.String("config", "./config.json", "Supply the location of MrRedis configuration file")
	flag.Parse()

	cfg_file, err := ioutil.ReadFile(*cfg_file_name)

	if err != nil {
		log.Printf("Error Reading the configration file. Resorting to default values")
	}
	Cfg := NewMrRedisDefaultConfig()
	err = json.Unmarshal(cfg_file, &Cfg)
	if err != nil {
		log.Fatalf("Error parsing the config file %v", err)
	}
	log.Printf("Configuration file is = %v", Cfg)

	log.Printf("*****************************************************************")
	log.Printf("*********************Starting MrRedis-Scheduler******************")
	log.Printf("*****************************************************************")
	//Command line argument parsing

	//Initalize the common entities like store, store configuration etc.
	isInit, err := types.Initialize(Cfg.DBType, Cfg.DBEndPoint)
	if err != nil || isInit != true {
		log.Fatalf("Failed to intialize Error:%v return %v", err, isInit)
	}

	//Start the Mesos library
	go mesoslib.Run(Cfg.MasterIP, Cfg.MasterPort, Cfg.ArtifactIP, Cfg.ArtifactPort, Cfg.ExecutorPath, Cfg.RedisPath, Cfg.DBType, Cfg.DBEndPoint)

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
