package etcdStore

import (
	_ "fmt"
	_ "log"
	"strings"
	"time"

	"github.com/coreos/etcd/Godeps/_workspace/src/golang.org/x/net/context"
	cli "github.com/coreos/etcd/client"

	"../../types"
)

type etcdDB struct {
	C       cli.Client      //The client context
	Kapi    cli.KeysAPI     //The api context for Get/Set/Delete/Update/Watcher etc.,
	Ctx     context.Context //Context for the connection mostly set to context.Background
	Cfg     cli.Config      //Configuration details of the connection should be loaded from a configuration file
	isSetup bool            //Has this been setup
}

func New() *etcdDB {
	return &etcdDB{isSetup: false}
}

func (db *etcdDB) Login() error {

	var err error
	db.C, err = cli.New(db.Cfg)
	if err != nil {

		return err
	}
	db.Kapi = cli.NewKeysAPI(db.C)
	db.Ctx = context.Background()

	return nil
}

// Setup will create/establish connection with the etcd store and also setup
// the nessary environment if etcd is running for the first time
// MrRedis will look for the following location in the central store
// /MrRedis/Instances/...... -> Will have the entries of all the instances
// /MrRedis/Config/....		-> Will have the entries of all the config information

func (db *etcdDB) Setup(config string) error {
	var err error
	db.Cfg = cli.Config{
		Endpoints: []string{config},
		Transport: cli.DefaultTransport,
		// set timeout per request to fail fast when the target endpoint is unavailable
		HeaderTimeoutPerRequest: time.Second,
	}

	err = db.Login()
	if err != nil {
		return err
	}

	err = db.CreateSection(types.ETC_BASE_DIR)
	if err != nil && strings.Contains(err.Error(), "Key already exists") != true {
		return err
	}

	err = db.CreateSection(types.ETC_INST_DIR)
	if err != nil && strings.Contains(err.Error(), "Key already exists") != true {
		return err
	}

	err = db.CreateSection(types.ETC_CONF_DIR)
	if err != nil && strings.Contains(err.Error(), "Key already exists") != true {
		return err
	}

	db.isSetup = true
	return nil
}

func (db *etcdDB) IsSetup() bool {
	return db.isSetup
}

func (db *etcdDB) Set(Key string, Value []byte) error {

	_, err := db.Kapi.Set(db.Ctx, Key, string(Value), nil)
	return err

}

func (db *etcdDB) Get(Key string) (error, []byte) {

	resp, err := db.Kapi.Get(db.Ctx, Key, nil)

	if err != nil {
		return err, []byte{}
	} else {
		return nil, []byte(resp.Node.Value)
	}

}

func (db *etcdDB) IsDir(Key string) (error, bool) {
	resp, err := db.Kapi.Get(db.Ctx, Key, nil)

	if err != nil {
		return err, false
	} else {
		return nil, resp.Node.Dir
	}
}

func (db *etcdDB) Update(Key string, Value []byte, Lock bool) error {

	return nil
}

func (db *etcdDB) Del(Key string) error {

	_, err := db.Kapi.Delete(db.Ctx, Key, nil)

	if err != nil {
		return nil
	}

	return nil

}

//CreateSection will create a directory in etcd store

func (db *etcdDB) CreateSection(Key string) error {

	_, err := db.Kapi.Set(db.Ctx, Key, "", &cli.SetOptions{Dir: true, PrevExist: cli.PrevNoExist})

	if err != nil {
		return err
	}

	return nil
}

// Delete section will delete a directory optionally delete

func (db *etcdDB) DeleteSection(Key string) error {

	_, err := db.Kapi.Delete(db.Ctx, Key, &cli.DeleteOptions{Dir: true})
	return err
}

func (db *etcdDB) ListSection(Key string, Recursive bool) []types.Rec {

	return nil
}

func (db *etcdDB) CleanSlate() error {

	_, err := db.Kapi.Delete(db.Ctx, types.ETC_BASE_DIR, &cli.DeleteOptions{Dir: true, Recursive: true})

	return err
}
