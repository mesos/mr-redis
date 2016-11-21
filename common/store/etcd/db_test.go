package etcd

import (
	"fmt"
	cli "github.com/coreos/etcd/client"
	"golang.org/x/net/context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestMain(M *testing.M) {

	//Run the tests
	M.Run()

}

//new with input
func TestNew(T *testing.T) {
	var db etcdDB

	New()

	if db.isSetup != false {
		T.Fail()
	}
}

// Login with endpoint
func TestLoginWithEndPoint(T *testing.T) {
	var db etcdDB

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "")
	}))

	defer ts.Close()

	config := ts.URL

	db.Cfg = cli.Config{
		Endpoints: []string{config},
	}

	err := db.Login()
	if err != nil {
		T.FailNow()
	}
}

// Login without endpoint
func TestLogWithoutEndPoint(T *testing.T) {
	var db etcdDB

	db.Cfg = cli.Config{
		Endpoints: []string{},
	}

	err := db.Login()

	if err == nil {
		//Error cannot be nil
		T.Fail()
	}

}

//setup with correct config values
func TestSetUpWithConfig(T *testing.T) {
	var db etcdDB
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"action":"set","node":{"dir": true,"key":"/test","value":"Hello","modifiedIndex":4,"createdIndex":4}}`)
	}))
	defer ts.Close()

	config := ts.URL

	db.Setup(config)
}

//setup without  config values
func TestSetUpWithoutConfig(T *testing.T) {

	var db etcdDB

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "A-ok")
	}))
	defer ts.Close()

	config := ts.URL

	err := db.Setup(config)

	if err != nil && !strings.Contains(err.Error(), "response is invalid json. The endpoint is probably not valid etcd cluster endpoint") {
		//If its some other error then fail
		T.Fail()
	}

}

//issetup with input
func TestIsSetup(T *testing.T) {

	var db etcdDB

	resp := db.IsSetup()

	if resp != false {
		T.Fail()
	}
}

//createsection inside etcd with key
func TestCreateSectionWithKey(T *testing.T) {

	var db etcdDB

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"action":"set","node":{"key":"/testDir","value":"Hello","modifiedIndex":4,"createdIndex":4}}`)
	}))
	defer ts.Close()

	config := ts.URL

	db.Cfg = cli.Config{
		Endpoints: []string{config},
		Transport: cli.DefaultTransport,
		// set timeout per request to fail fast when the target endpoint is unavailable
		HeaderTimeoutPerRequest: time.Second,
	}

	db.C, _ = cli.New(db.Cfg)

	db.Kapi = cli.NewKeysAPI(db.C)
	db.Ctx = context.Background()

	err := db.CreateSection("/testDir")

	if err != nil {
		T.Fail()
	}
}

//createsection without key
func TestCreateSectionWithoutKey(T *testing.T) {

	var db etcdDB

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "A-ok")
	}))
	defer ts.Close()

	config := ts.URL

	db.Cfg = cli.Config{
		Endpoints: []string{config},
		Transport: cli.DefaultTransport,
		// set timeout per request to fail fast when the target endpoint is unavailable
		HeaderTimeoutPerRequest: time.Second,
	}

	db.C, _ = cli.New(db.Cfg)

	db.Kapi = cli.NewKeysAPI(db.C)
	db.Ctx = context.Background()

	err := db.CreateSection("")

	if err != nil && !strings.Contains(err.Error(), "response is invalid json. The endpoint is probably not valid etcd cluster endpoint") {
		//If its some other error then fail
		T.Fail()
	}
}

//set with key
func TestSetWithKey(T *testing.T) {

	var db etcdDB

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"action":"set","node":{"key":"/testDir","value":"Hello","modifiedIndex":4,"createdIndex":4}}`)
	}))
	defer ts.Close()

	config := ts.URL

	db.Cfg = cli.Config{
		Endpoints: []string{config},
		Transport: cli.DefaultTransport,
		// set timeout per request to fail fast when the target endpoint is unavailable
		HeaderTimeoutPerRequest: time.Second,
	}

	db.C, _ = cli.New(db.Cfg)

	db.Kapi = cli.NewKeysAPI(db.C)
	db.Ctx = context.Background()

	err := db.Set("/testDir", "Hello")

	if err != nil {
		T.Fail()
	}

}

//set without key
func TestSetWithoutKey(T *testing.T) {

	var db etcdDB

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "A-OK")
	}))
	defer ts.Close()

	config := ts.URL

	db.Cfg = cli.Config{
		Endpoints: []string{config},
		Transport: cli.DefaultTransport,
		// set timeout per request to fail fast when the target endpoint is unavailable
		HeaderTimeoutPerRequest: time.Second,
	}

	db.C, _ = cli.New(db.Cfg)

	db.Kapi = cli.NewKeysAPI(db.C)
	db.Ctx = context.Background()

	err := db.Set("", "Hello")

	if err != nil && !strings.Contains(err.Error(), "response is invalid json. The endpoint is probably not valid etcd cluster endpoint") {
		//If its some other error then fail
		T.Fail()
	}
}

//get with valid config
func TestGetValidConfig(T *testing.T) {

	var db etcdDB

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"action":"get","node":{"key":"/testDir","value":"Hello","modifiedIndex":4,"createdIndex":4}}`)
	}))
	defer ts.Close()

	config := ts.URL

	db.Cfg = cli.Config{
		Endpoints: []string{config},
		Transport: cli.DefaultTransport,
		// set timeout per request to fail fast when the target endpoint is unavailable
		HeaderTimeoutPerRequest: time.Second,
	}

	db.C, _ = cli.New(db.Cfg)

	db.Kapi = cli.NewKeysAPI(db.C)
	db.Ctx = context.Background()

	db.Get("/testDir")

}

//get without key
func TestGetWithoutKey(T *testing.T) {

	var db etcdDB

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "A-ok")
	}))
	defer ts.Close()

	config := ts.URL

	db.Cfg = cli.Config{
		Endpoints: []string{config},
		Transport: cli.DefaultTransport,
		// set timeout per request to fail fast when the target endpoint is unavailable
		HeaderTimeoutPerRequest: time.Second,
	}

	db.C, _ = cli.New(db.Cfg)

	db.Kapi = cli.NewKeysAPI(db.C)
	db.Ctx = context.Background()

	resp, err := db.Get("")

	if resp != "" && !strings.Contains(err.Error(), "response is invalid json. The endpoint is probably not valid etcd cluster endpoint") {
		//If its some other error then fail
		T.Fail()
	}
}

//isDir with correct input
func TestIsDir(T *testing.T) {

	var db etcdDB

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"action":"get","node":{"key":"/testDir","value":"Hello","modifiedIndex":4,"createdIndex":4}}`)
	}))
	defer ts.Close()

	config := ts.URL

	db.Cfg = cli.Config{
		Endpoints: []string{config},
		Transport: cli.DefaultTransport,
		// set timeout per request to fail fast when the target endpoint is unavailable
		HeaderTimeoutPerRequest: time.Second,
	}

	db.C, _ = cli.New(db.Cfg)

	db.Kapi = cli.NewKeysAPI(db.C)
	db.Ctx = context.Background()

	db.IsDir("/testDir")

}

//isDir without key
func TestIsDirWithoutKey(T *testing.T) {

	var db etcdDB

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "A-ok")
	}))
	defer ts.Close()

	config := ts.URL

	db.Cfg = cli.Config{
		Endpoints: []string{config},
		Transport: cli.DefaultTransport,
		// set timeout per request to fail fast when the target endpoint is unavailable
		HeaderTimeoutPerRequest: time.Second,
	}

	db.C, _ = cli.New(db.Cfg)

	db.Kapi = cli.NewKeysAPI(db.C)
	db.Ctx = context.Background()

	err, resp := db.IsDir("")

	if resp != false && !strings.Contains(err.Error(), "response is invalid json. The endpoint is probably not valid etcd cluster endpoint") {
		//If its some other error then fail
		T.Fail()
	}

}

//IsKey with correct key
func TestIsKey(T *testing.T) {

	var db etcdDB

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"action":"get","node":{"key":"/testDir","value":"Hello","modifiedIndex":4,"createdIndex":4}}`)
	}))
	defer ts.Close()

	config := ts.URL

	db.Cfg = cli.Config{
		Endpoints: []string{config},
		Transport: cli.DefaultTransport,
		// set timeout per request to fail fast when the target endpoint is unavailable
		HeaderTimeoutPerRequest: time.Second,
	}

	db.C, _ = cli.New(db.Cfg)

	db.Kapi = cli.NewKeysAPI(db.C)
	db.Ctx = context.Background()

	db.IsKey("/testDir")

}

//IsKey without key
func TestIsKeyWithoutKey(T *testing.T) {

	var db etcdDB

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "A-ok")
	}))
	defer ts.Close()

	config := ts.URL

	db.Cfg = cli.Config{
		Endpoints: []string{config},
		Transport: cli.DefaultTransport,
		// set timeout per request to fail fast when the target endpoint is unavailable
		HeaderTimeoutPerRequest: time.Second,
	}

	db.C, _ = cli.New(db.Cfg)

	db.Kapi = cli.NewKeysAPI(db.C)
	db.Ctx = context.Background()

	resp, err := db.IsKey("")

	if resp != false && !strings.Contains(err.Error(), "response is invalid json. The endpoint is probably not valid etcd cluster endpoint") {
		//If its some other error then fail
		T.Fail()
	}

}

//update with correct key
func TestUpdate(T *testing.T) {
	var db etcdDB
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"action":"set","node":{"key":"/testDir","value":"Hello","modifiedIndex":4,"createdIndex":4}}`)
	}))

	defer ts.Close()

	config := ts.URL
	db.Cfg = cli.Config{
		Endpoints: []string{config},
		Transport: cli.DefaultTransport,
		// set timeout per request to fail fast when the target endpoint is unavailable
		HeaderTimeoutPerRequest: time.Second,
	}

	db.C, _ = cli.New(db.Cfg)

	db.Kapi = cli.NewKeysAPI(db.C)
	db.Ctx = context.Background()

	err := db.Update("/testDir", "Hello", true)

	if err != nil {
		T.Fail()
	}

}

//update without key
func TestUpdateWithoutKey(T *testing.T) {
	var db etcdDB
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "A-ok")
	}))

	defer ts.Close()
	config := ts.URL

	db.Cfg = cli.Config{
		Endpoints: []string{config},
		Transport: cli.DefaultTransport,
		// set timeout per request to fail fast when the target endpoint is unavailable
		HeaderTimeoutPerRequest: time.Second,
	}

	db.C, _ = cli.New(db.Cfg)

	db.Kapi = cli.NewKeysAPI(db.C)
	db.Ctx = context.Background()

	err := db.Update("", "Hello", true)

	if err != nil && !strings.Contains(err.Error(), "response is invalid json. The endpoint is probably not valid etcd cluster endpoint") {
		//If its some other error then fail
		T.Fail()
	}
}

//delete section inside etcd without key
func TestDeleteSectionWithoutKey(T *testing.T) {

	var db etcdDB

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "A-ok")
	}))
	defer ts.Close()

	config := ts.URL

	db.Cfg = cli.Config{
		Endpoints: []string{config},
		Transport: cli.DefaultTransport,
		// set timeout per request to fail fast when the target endpoint is unavailable
		HeaderTimeoutPerRequest: time.Second,
	}

	db.C, _ = cli.New(db.Cfg)

	db.Kapi = cli.NewKeysAPI(db.C)
	db.Ctx = context.Background()

	err := db.DeleteSection("")

	if err != nil && !strings.Contains(err.Error(), "response is invalid json. The endpoint is probably not valid etcd cluster endpoint") {
		//If its some other error then fail
		T.Fail()
	}
}

//delete with key
func TestDel(T *testing.T) {

	var db etcdDB

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"action":"delete","node":{"key":"/testDir","value":"Hello","modifiedIndex":4,"createdIndex":4}}`)
	}))
	defer ts.Close()

	config := ts.URL

	db.Cfg = cli.Config{
		Endpoints: []string{config},
		Transport: cli.DefaultTransport,
		// set timeout per request to fail fast when the target endpoint is unavailable
		HeaderTimeoutPerRequest: time.Second,
	}

	db.C, _ = cli.New(db.Cfg)

	db.Kapi = cli.NewKeysAPI(db.C)
	db.Ctx = context.Background()

	db.Del("/testDir")

}

//delete without key
func TestDelWithoutKey(T *testing.T) {

	var db etcdDB

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "A-ok")
	}))
	defer ts.Close()

	config := ts.URL

	db.Cfg = cli.Config{
		Endpoints: []string{config},
		Transport: cli.DefaultTransport,
		// set timeout per request to fail fast when the target endpoint is unavailable
		HeaderTimeoutPerRequest: time.Second,
	}

	db.C, _ = cli.New(db.Cfg)

	db.Kapi = cli.NewKeysAPI(db.C)
	db.Ctx = context.Background()

	err := db.Del("")

	if err != nil && !strings.Contains(err.Error(), "response is invalid json. The endpoint is probably not valid etcd cluster endpoint") {
		//If its some other error then fail
		T.Fail()
	}
}

//list section with key
func TestListSection(T *testing.T) {

	var db etcdDB

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"action":"get","node":{"key":"/testDir","value":"Hello","modifiedIndex":4,"createdIndex":4}}`)
	}))
	defer ts.Close()

	config := ts.URL

	db.Cfg = cli.Config{
		Endpoints: []string{config},
		Transport: cli.DefaultTransport,
		// set timeout per request to fail fast when the target endpoint is unavailable
		HeaderTimeoutPerRequest: time.Second,
	}

	db.C, _ = cli.New(db.Cfg)

	db.Kapi = cli.NewKeysAPI(db.C)
	db.Ctx = context.Background()

	db.ListSection("/testDir", true)

}

//list section without key
func TestListSectionWithoutKey(T *testing.T) {

	var db etcdDB

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "A-ok")
	}))
	defer ts.Close()

	config := ts.URL

	db.Cfg = cli.Config{
		Endpoints: []string{config},
		Transport: cli.DefaultTransport,
		// set timeout per request to fail fast when the target endpoint is unavailable
		HeaderTimeoutPerRequest: time.Second,
	}

	db.C, _ = cli.New(db.Cfg)

	db.Kapi = cli.NewKeysAPI(db.C)
	db.Ctx = context.Background()

	_, err := db.ListSection("", true)

	if err != nil && !strings.Contains(err.Error(), "response is invalid json. The endpoint is probably not valid etcd cluster endpoint") {
		//If its some other error then fail
		T.Fail()
	}
}

//test clean slate 
func TestCleanSlate(T *testing.T) {

	var db etcdDB

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "A-ok")
	}))
	defer ts.Close()

	config := ts.URL

	db.Cfg = cli.Config{
		Endpoints: []string{config},
		Transport: cli.DefaultTransport,
		// set timeout per request to fail fast when the target endpoint is unavailable
		HeaderTimeoutPerRequest: time.Second,
	}

	db.C, _ = cli.New(db.Cfg)

	db.Kapi = cli.NewKeysAPI(db.C)
	db.Ctx = context.Background()


	err := db.CleanSlate()

	if err != nil && !strings.Contains(err.Error(), "response is invalid json. The endpoint is probably not valid etcd cluster endpoint") {
		//If its some other error then fail
		T.Fail()
	}
}
