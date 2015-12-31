package main

import (
	"./store"
	"./store/etcd"
	"fmt"
)

func main() {

	/* Create a variable for the DB interface */
	var etcdConn store.DB
	/* Create a variable on the implemented datastore , in this case etcd */
	etcdConn = etcd.New()

	/* Setup etcd with the etcd endpoint*/
	err := etcdConn.Setup("http://10.11.12.24:2379")
	/* Test if this is setup */
	fmt.Printf("IsSetup %v\n", etcdConn.IsSetup())

	/*Create a dummy Instance */
	err = etcdConn.Set("HelloWorld", []byte("How are you"))

	if err != nil {
		fmt.Println("Error setting value ", err)
	}

	err, value := etcdConn.Get("HelloWorld")
	fmt.Printf("GET err:%v, value:%v\n", err, string(value))

}
