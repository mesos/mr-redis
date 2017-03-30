package zookeeper

import (
	"fmt"
	"strings"
	"time"

	"github.com/samuel/go-zookeeper/zk"
)

const (
	ETC_BASE_DIR = "/MrRedis"
	ETC_INST_DIR = ETC_BASE_DIR + "/Instances"
	ETC_CONF_DIR = ETC_BASE_DIR + "/Config"
)

var (
	DEF_ACL = zk.WorldACL(zk.PermAll)
)

type zkDB struct {
	Con     *zk.Conn
	Cfg     string
	isSetup bool
}

func New() *zkDB {
	return &zkDB{isSetup: false}
}

func (db *zkDB) Login() error {
	var err error
	db.Con, _, err = zk.Connect([]string{db.Cfg}, time.Second*60)
	if err != nil {
		panic(err)
	}
	return nil
}

func (db *zkDB) IsSetup() bool {
	return db.isSetup
}

func (db *zkDB) Set(Key string, Value string) error {
	fmt.Printf("ZOO K=%s V=%s\n", Key, Value)
	if _, err := db.Con.Set(Key, []byte(Value), -1); err != nil {
		_, err := db.Con.Create(Key, []byte(Value), 0, DEF_ACL)
		if err != nil {
			fmt.Printf("Create error %v\n", err)
			return err
		}
	}
	return nil
}

func (db *zkDB) Get(Key string) (string, error) {
	result, _, err := db.Con.Get(Key)
	return string(result), err
}

func (db *zkDB) IsDir(Key string) (error, bool) {
	result, _, err := db.Con.Children(Key)
	if err != nil || len(result) == 0 {
		return err, false
	}
	return nil, true
}

func (db *zkDB) IsKey(Key string) (bool, error) {

	result, _, err := db.Con.Exists(Key)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (db *zkDB) Update(Key string, Value string, Lock bool) error {
	return nil
}

func (db *zkDB) Del(Key string) error {
	return nil
}

//CreateSection will create a directory in zookeeper store
func (db *zkDB) CreateSection(Key string) error {
	fmt.Printf("ZOO CREATE SECTION K=%s \n", Key)
	Key = strings.TrimSuffix(Key, "/")
	if _, err := db.Con.Set(Key, []byte{'.'}, -1); err != nil {
		_, err = db.Con.Create(Key, []byte{'.'}, 0, DEF_ACL)
		if err != nil {
			fmt.Printf("Create Error %v\n", err)
			return err
		}
	}
	return nil
}

func (db *zkDB) Setup(config string) error {
	var err error
	i := strings.Index(config, "//")
	if i > -1 {
		db.Cfg = config[i+2:]
	} else {
		db.Cfg = config
	}

	err = db.Login()
	if err != nil {
		return err
	}

	err = db.CreateSection(ETC_BASE_DIR)
	if err != nil && strings.Contains(err.Error(), "Key already exists") != true {
		return err
	}
	err = db.CreateSection(ETC_INST_DIR)
	if err != nil && strings.Contains(err.Error(), "Key already exists") != true {
		return err
	}
	err = db.CreateSection(ETC_CONF_DIR)
	if err != nil && strings.Contains(err.Error(), "Key already exists") != true {
		return err
	}
	db.isSetup = true
	return nil
}

func (db *zkDB) CleanSlate() error {
	//	_, err := db.Kapi.Delete(db.Ctx, ETC_BASE_DIR, &cli.DeleteOptions{Dir: true, Recursive: true})
	return nil
}

//DeleteSection section will delete a directory optionally delete
func (db *zkDB) DeleteSection(Key string) error {

	//	_, err := db.Kapi.Delete(db.Ctx, Key, &cli.DeleteOptions{Dir: true})
	return nil
}

//ListSection will list a directory
func (db *zkDB) ListSection(Key string, Recursive bool) ([]string, error) {
	Key = strings.TrimSuffix(Key, "/")
	result, _, err := db.Con.Children(Key)
	return result, err
}
