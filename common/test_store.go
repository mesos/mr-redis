package main

import (
	"./id"
	"./store/etcd"
	"./types"
	"log"
)

func main() {

	/* Create a variable on the implemented datastore , in this case etcd */
	types.Gdb = etcd.New()

	/* Setup etcd with the etcd endpoint*/
	err := types.Gdb.Setup("http://10.11.12.24:2379")
	/* Test if this is setup */
	log.Printf("IsSetup %v\n", types.Gdb.IsSetup())

	/*Create a dummy Instance */
	uid, _ := id.NewUUID()
	uid_str := uid.String()
	err = types.Gdb.Set(uid_str, "How are you")

	if err != nil {
		log.Println("Error setting value ", err)
	}

	value, err := types.Gdb.Get(uid_str)
	log.Printf("GET err:%v, value:%v\n", err, value)

	/* Use the instance library */
	I := types.NewInstance("TestInstance", "S", 1, 0)
	I.Status = "Innactive"

	/* Write the entier data to db store */
	if !I.Sync() {
		log.Printf("Failed to write the structure to etcd")
	}

	J := &types.Instance{}

	J.Name = "TestInstance"

	if !J.Load() {
		log.Printf("Failed to read the structure from the store")
	}

	/* Print the contents of the structure read fromt he store */
	log.Printf("Value read from store {Name : %s, Type : %s, Masters : %d, Slaves : %d, Status : %s}", J.Name, J.Type, J.Masters, J.Slaves, J.Status)

}
