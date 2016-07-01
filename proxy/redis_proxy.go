package main

import (
	"encoding/json"
	"flag"
<<<<<<< HEAD
=======
	"fmt"
>>>>>>> c0030928f2894a32452cdbd34c194cc4917bd22c
	"io"
	"io/ioutil"
	"log"
	"net"
<<<<<<< HEAD
	"runtime"
)

type Config struct {
	Entries []Entry
=======
	"net/http"
)

var ConfigMap map[string]Entry

type Config struct {
	HTTPPort string  //HTTP server Port number that we should bind to
	Entries  []Entry //List of proxy entries
>>>>>>> c0030928f2894a32452cdbd34c194cc4917bd22c
}

type Entry struct {
	Name string
	Pair PorxyPair
}

type PorxyPair struct {
	From string //IP:PORT pair
	To   string //IP:PORT pair
}

<<<<<<< HEAD
func HandleConnection(E Entry) error {
	log.Printf("HandleConnection() %v", E)
	src, err := net.Listen("tcp", E.Pair.From)

=======
//This structure is used by the HTTP PUT request to change the IP address of the destination on the fly
type HTTPUpdate struct {
	Name string
	Addr string
}

func HandleConnection(E Entry) error {

	var CurrentE Entry //A Temp variable to get the latest Desination proxy value
	var OK bool

	log.Printf("HandleConnection() %v", E)
	src, err := net.Listen("tcp", E.Pair.From)
>>>>>>> c0030928f2894a32452cdbd34c194cc4917bd22c
	if err != nil {
		log.Printf("Error binding to the IP %v", err)
		return err
	}
	defer src.Close()

<<<<<<< HEAD
=======
	//Add this in the global Map so that it can be updated dynamically by HTTP apis
	ConfigMap[E.Name] = E

>>>>>>> c0030928f2894a32452cdbd34c194cc4917bd22c
	for {
		conn, err := src.Accept()
		if err != nil {
			log.Printf("Error accepting a new connection %v", err)
			continue
		}

<<<<<<< HEAD
		//Start a Lamda for performing the proxy
		go func(E Entry, F net.Conn) {

			var buf []byte

=======
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

>>>>>>> c0030928f2894a32452cdbd34c194cc4917bd22c
			T, err := net.Dial("tcp", E.Pair.To)
			if err != nil {
				log.Printf("Unable to connect to the Destination %s %v", E.Pair.To, err)
				return
			}
			defer T.Close()
			defer F.Close()

			go io.Copy(F, T)
			io.Copy(T, F)

<<<<<<< HEAD
		}(E, conn)
=======
		}(CurrentE, conn)
>>>>>>> c0030928f2894a32452cdbd34c194cc4917bd22c
	}
	return nil
}

<<<<<<< HEAD
func main() {

	var Cfg Config
	runtime.GOMAXPROCS(runtime.NumCPU())
=======
func HandleHTTPUpdate(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, Going to Update %s! Method=%s\n", r.URL.Path[1:], r.Method)
	if r.Method == "PUT" {
		//This can be used for updating an existing variable
		content, err := ioutil.ReadAll(r.Body)
		r.Body.Close()
		if err != nil {
			fmt.Fprint(w, "Error understanding the Body %v", err)
			log.Printf("Error understanding the Body %v", err)
			return
		}

		var val HTTPUpdate
		var CurrentE Entry
		var OK bool
		err = json.Unmarshal(content, &val)
		if err != nil {
			fmt.Fprintf(w, "Wrong json format %v", err)
			log.Printf("Wrong json format %v", err)
			return
		}
		if CurrentE, OK = ConfigMap[val.Name]; !OK {
			log.Printf("Error Proxy entry is incorrect / empty for %s", val.Name)
			fmt.Fprintf(w, "Error Proxy entry is incorrect / empty for %s", val.Name)
			return
		}
		log.Printf("Updating From porxy for %s From %s TO %s", val.Name, CurrentE.Pair.To, val.Addr)
		CurrentE.Pair.To = val.Addr
		ConfigMap[val.Name] = CurrentE
		return
	}
	return
}

func HandleHTTPGet(w http.ResponseWriter, r *http.Request) {
	ret_bytes, err := json.MarshalIndent(ConfigMap, " ", "  ")
	if err != nil {
		log.Printf("Error Marshalling HandleHTTPGet() %v", err)
		fmt.Fprintf(w, "Error Marshalling HandleHTTPGet() %v", err)
		return

	}
	fmt.Fprintf(w, "Current Config: %s", string(ret_bytes))
	return
}

func main() {

	var Cfg Config

	//Initialize the global Config map
	ConfigMap = make(map[string]Entry)
>>>>>>> c0030928f2894a32452cdbd34c194cc4917bd22c

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

<<<<<<< HEAD
=======
	http.HandleFunc("/Update/", HandleHTTPUpdate)
	http.HandleFunc("/Get/", HandleHTTPGet)
	log.Fatal(http.ListenAndServe(":"+Cfg.HTTPPort, nil))

>>>>>>> c0030928f2894a32452cdbd34c194cc4917bd22c
	//Wait indefinately
	waitCh := make(chan bool)
	<-waitCh

}
