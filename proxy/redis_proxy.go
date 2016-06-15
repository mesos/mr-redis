package main

import (
	"encoding/json"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"net"
	"runtime"
)

type Config struct {
	Entries []Entry
}

type Entry struct {
	Name string
	Pair PorxyPair
}

type PorxyPair struct {
	From string //IP:PORT pair
	To   string //IP:PORT pair
}

func HandleConnection(E Entry) error {
	log.Printf("HandleConnection() %v", E)
	src, err := net.Listen("tcp", E.Pair.From)

	if err != nil {
		log.Printf("Error binding to the IP %v", err)
		return err
	}
	defer src.Close()

	for {
		conn, err := src.Accept()
		if err != nil {
			log.Printf("Error accepting a new connection %v", err)
			continue
		}

		//Start a Lamda for performing the proxy
		go func(E Entry, F net.Conn) {

			var buf []byte

			T, err := net.Dial("tcp", E.Pair.To)
			if err != nil {
				log.Printf("Unable to connect to the Destination %s %v", E.Pair.To, err)
				return
			}
			defer T.Close()
			defer F.Close()

			go io.Copy(F, T)
			io.Copy(T, F)

		}(E, conn)
	}
	return nil
}

func main() {

	var Cfg Config
	runtime.GOMAXPROCS(runtime.NumCPU())

	//Read a config file that has json update the config files
	cfg_file_name := flag.String("config", "./config.json", "Supply the location of MrRedis configuration file")

	log.Printf("The config file name is %s ", *cfg_file_name)
	cfg_file, err := ioutil.ReadFile(*cfg_file_name)

	if err != nil {
		log.Printf("Error Reading the configration file. Resorting to default values")
	}
	err = json.Unmarshal(cfg_file, &Cfg)
	if err != nil {
		log.Fatalf("Error parsing the config file %v", err)
		return
	}
	log.Printf("Configuration file is = %v", Cfg)

	//Hanlde each connection

	for _, E := range Cfg.Entries {
		go HandleConnection(E)
	}

	//Wait indefinately
	waitCh := make(chan bool)
	<-waitCh

}
