package types

import (
	"log"

	"../store/etcd"
)

func Initialize(dbtype string, config string) (bool, error) {

	switch dbtype {
	case "etcd":
		Gdb = etcd.New()
		err := Gdb.Setup(config)
		if err != nil {
			log.Fatalf("Failed to setup etcd database error:%v", err)
		}
		return Gdb.IsSetup(), nil
		break
	}

	return false, nil
}
