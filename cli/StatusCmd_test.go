package main

import (
	"flag"
	"fmt"
	"github.com/codegangsta/cli"
	"net/http"
	"net/http/httptest"
	"testing"
)

//statusof with valid input
func TestStatusOf(T *testing.T) {
	status := `{"Name":"TestInstance","Type":"MS","Status":"RUNNING","Capacity":200,"Master":{"IP":"10.11.12.123","Port":"6382","MemoryCapacity":200,"MemoryUsed":1904432,"Uptime":1623,"ClientsConnected":1,"LastSyncedToMaster":0},"Slaves":[{"IP":"10.11.12.121","Port":"6381","MemoryCapacity":200,"MemoryUsed":834904,"Uptime":1619,"ClientsConnected":2,"LastSyncedToMaster":9},{"IP":"10.11.12.121","Port":"6382","MemoryCapacity":200,"MemoryUsed":834904,"Uptime":1619,"ClientsConnected":2,"LastSyncedToMaster":9}]}`

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, status)
	}))

	defer ts.Close()

	MrRedisFW = ts.URL

	statusOf("TestInstance")

}

//statusof with invalid url
func TestStatusOfWithInvalidURL(T *testing.T) {

	url := DC_INVALID_ENDPOINT

	MrRedisFW = url

	statusOf("Test")

}

//statusof with invalid json
func TestStatusOfWithInvalidJson(T *testing.T) {

	status := `{TestInstance MS CREATING 100 { } [,]}`

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, status)
	}))

	defer ts.Close()

	MrRedisFW = ts.URL

	statusOf("TestInstance")

}

//status with valid input
func TestStatusCmd(T *testing.T) {

	status := `{"Name":"TestInstance","Type":"MS","Status":"RUNNING","Capacity":200,"Master":{"IP":"10.11.12.123","Port":"6382","MemoryCapacity":200,"MemoryUsed":1904432,"Uptime":1623,"ClientsConnected":1,"LastSyncedToMaster":0},"Slaves":[{"IP":"10.11.12.121","Port":"6381","MemoryCapacity":200,"MemoryUsed":834904,"Uptime":1619,"ClientsConnected":2,"LastSyncedToMaster":9},{"IP":"10.11.12.121","Port":"6382","MemoryCapacity":200,"MemoryUsed":834904,"Uptime":1619,"ClientsConnected":2,"LastSyncedToMaster":9}]}`

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, status)
	}))

	defer ts.Close()

	MrRedisFW = ts.URL

	set := flag.NewFlagSet("test", 0)
	set.String("name", "TestInstance", "doc")
	c := cli.NewContext(nil, set, nil)

	StatusCmd(c)
}

//status with empty name
func TestStatusCmdWithEmptyName(T *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "A-ok")
	}))

	defer ts.Close()

	MrRedisFW = ts.URL

	set := flag.NewFlagSet("test", 0)
	set.String("name", "", "doc")
	c := cli.NewContext(nil, set, nil)

	StatusCmd(c)
}

//status with not found
func TestStatusCmdWithNotFound(T *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		fmt.Fprintln(w, "A-Ok")
	}))

	defer ts.Close()

	MrRedisFW = ts.URL

	set := flag.NewFlagSet("test", 0)
	set.String("name", "Test", "doc")
	c := cli.NewContext(nil, set, nil)

	StatusCmd(c)
}
