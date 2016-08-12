package main

import (
	"flag"
	"fmt"
	"github.com/codegangsta/cli"
	"net/http"
	"net/http/httptest"
	"testing"
)

var DC_INVALID_ENDPOINT = "http://127.127.127.127:5656"

//create with valid input
func TestCreatCommand(T *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(201)
		fmt.Fprintln(w, "Created OK")
	}))

	defer ts.Close()

	MrRedisFW = ts.URL
	set := flag.NewFlagSet("test", 0)
	set.String("name", "Test", "doc")
	set.Int("mem", 10, "doc")
	set.Int("slaves", 1, "doc")
	c := cli.NewContext(nil, set, nil)

	CreateCmd(c)

}

//create with empty name
func TestCreatCommandWithEmptyName(T *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "A-OK")
	}))

	defer ts.Close()

	MrRedisFW = ts.URL
	set := flag.NewFlagSet("test", 0)
	set.String("name", "", "doc")
	set.Int("mem", 10, "doc")
	set.Int("slaves", 120, "doc")
	set.Bool("wait", false, "doc")
	c := cli.NewContext(nil, set, nil)

	CreateCmd(c)

}

//create with bad server
func TestCreatCommandBadServer(T *testing.T) {

	MrRedisFW = DC_INVALID_ENDPOINT + "/v1/CREATE/TestInstance"

	set := flag.NewFlagSet("test", 0)
	set.String("name", "TestInstance", "doc")
	set.Int("mem", 10, "doc")
	set.Int("slaves", 1, "doc")
	set.Bool("wait", false, "doc")
	c := cli.NewContext(nil, set, nil)

	CreateCmd(c)

}
