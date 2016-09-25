package cmd

import (

	"log"
	"bytes"
	"strings"
	"testing"
	typ "github.com/mesos/mr-redis/common/types"
)

var LogBuffer bytes.Buffer

func TestMain(M *testing.M) {
	//Initialize the Json files
	
	log.SetOutput(&LogBuffer)
	
	typ.Initialize("testing", "something")

	//Run the tests
	go Creator()
	
	M.Run()
}

func TestCreaterSingleMaster(t *testing.T) {

	tmpInstance := typ.NewInstance("Test1", "M", 1, 0, 100)
	tmpInstance.Status = typ.INST_STATUS_CREATING

	typ.Cchan <- typ.CreateMaster(tmpInstance)
	

	//
	for tskEle := typ.OfferList.Front(); tskEle != nil; {
		tsk := tskEle.Value.(typ.Offer)
		t.Logf("Task is %v", tsk)
		if strings.Contains(tsk.Taskname,"Test1") && tsk.Cpu == 1 && tsk.Mem == 100 && tsk.IsMaster == true && tsk.MasterIpPort == "" {
			typ.OfferList.Remove(tskEle)
			return
		}
		tskEle = tskEle.Next()
	}
	t.Fail()
}

func TestCreaterMultiMaster(t *testing.T) {

	tmpInstance := typ.NewInstance("Test1", "M", 10, 0, 100)
	tmpInstance.Status = typ.INST_STATUS_CREATING
	InstanceCount := 0
	typ.Cchan <- typ.CreateMaster(tmpInstance)

	//
	for {
		if typ.OfferList.Len() > 0 {
			break
		}
	}
	for tskEle := typ.OfferList.Front(); tskEle != nil; {
		tsk := tskEle.Value.(typ.Offer)
		log.Printf("Task is %v", tsk)
		if strings.Contains(tsk.Taskname,"Test1") && tsk.Cpu == 1 && tsk.Mem == 100 && tsk.IsMaster == true && tsk.MasterIpPort == "" {
			typ.OfferList.Remove(tskEle)
			InstanceCount++
		}
		tskEle = tskEle.Next()
	}
	if InstanceCount != 1 {
		log.Printf("Instance Count = %v", InstanceCount)
		t.Fail()
	}
}

func TestCreaterSlaves (t *testing.T) {

	tmpInstance := typ.NewInstance("Test1", "MS", 1, 10, 100)
	tmpInstance.Status = typ.INST_STATUS_CREATING
	InstanceCount := 0
	
	tmpInstance.Masters = 1
	tmpInstance.Mname = "Master::test"
	tmpProc := typ.NewProc("Master::test",100,"M","")
	tmpProc.IP = "127.0.0.1"
	tmpProc.Port = "8080"
	tmpInstance.Procs[tmpInstance.Mname] = tmpProc
	
	typ.Cchan <- typ.CreateSlaves(tmpInstance, 10)
	//
	for {
		if typ.OfferList.Len() >= 10 {
			t.Logf("OfferList len is %v", typ.OfferList.Len())
			break
		}
	}
	
	for tskEle := typ.OfferList.Front(); tskEle != nil; {
		tsk := tskEle.Value.(typ.Offer)
		t.Logf("Task is %v", tsk)
		if strings.Contains(tsk.Taskname,"Test1") && tsk.Cpu == 1 && tsk.Mem == 100 && tsk.IsMaster == false && tsk.MasterIpPort == "127.0.0.1:8080" {
			//typ.OfferList.Remove(tskEle)
			InstanceCount++
		}
		tskEle = tskEle.Next()
	}
	if InstanceCount != 10 {
		t.Logf("InstanceCount = %v", InstanceCount)
		t.Fail()
	}
	
}