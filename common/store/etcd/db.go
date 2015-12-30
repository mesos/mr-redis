package ETCD

import (
	"fmt"
	_ "log"
	"time"

	"github.com/coreos/etcd/Godeps/_workspace/src/golang.org/x/net/context"
	cli "github.com/coreos/etcd/client"
)

type etcdDB struct {
	C    cli.Client  //The client context
	Kapi cli.KeysAPI //The api context for Get/Set/Delete/Update/Watcher etc.,
	Ctx	context.Context //Context for the connection mostly set to context.Background
	Cfg cli.Config	//Configuration details of the connection should be loaded from a configuration file
	isSetup bool	//Has this been setup
}

func New() *etcdDB {
	return &etcdDB {isSetup:false}
}


func (db *etcdDB) Setup (config string) error {
	var err error
	db.Cfg = cli.Config{
        Endpoints:               []string{"http://10.11.12.24:2379"},
        Transport:               cli.DefaultTransport,
        // set timeout per request to fail fast when the target endpoint is unavailable
        HeaderTimeoutPerRequest: time.Second,
    }
	
	db.C, err = cli.New(db.Cfg)
    if err != nil {
        fmt.Printf("Error creatingthe client")
		return err
    }	
	db.Kapi = cli.NewKeysAPI(db.C)
	db.Ctx = context.Background()
	db.isSetup = true
	return nil
}

func (db *etcdDB) IsSetup () bool {
	return db.isSetup
}

func (db *etcdDB) Set (Key string, Value []byte) error {

	_, err := db.Kapi.Set(db.Ctx, Key, string(Value), nil)
	return err
	
}

func (db *etcdDB) Get (Key string) (error, []byte) {
	
	resp, err := db.Kapi.Get(db.Ctx, Key, nil)
	
	if err != nil {
		return err , []byte{}
	} else {
		return nil, []byte(resp.Node.Value)
	}
	
}

func (db *etcdDB) Del (Key string) error {

	return nil

}

func (db *etcdDB) CreateSection (Key string) error {

	return nil
}

func (db *etcdDB) DeleteSection (Key string) error {

	return nil
}

func (db *etcdDB) ListSection (Key string, Recursive bool) error {

	return nil
}

func (db *etcdDB) CleanSlate() error {

	return nil
}