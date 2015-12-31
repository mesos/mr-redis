package main

import (
	"./store"
	"./store/etcd"
	"fmt"
)

func main() {

	/* Create a variable for the DB interface */
	var etcd DB.DB
	/* Create a variable on the implemented datastore , in this case etcd */
	etcd = ETCD.New()

	/* Setup etcd with the etcd endpoint*/
	err := etcd.Setup("http://10.11.12.24:2379")
	/* Test if this is setup */
	fmt.Printf("IsSetup %v\n", etcd.IsSetup())

	/*Create a dummy Instance */
	err = etcd.Set("HelloWorld", []byte("How are you"))

	if err != nil {
		fmt.Println("Error setting value ", err)
	}

	err, value := etcd.Get("HelloWorld")
	fmt.Printf("GET err:%v, value:%v\n", err, string(value))

}
