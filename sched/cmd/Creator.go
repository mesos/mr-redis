package cmd

import (
	//"container/list"
	"log"
	//	"time"

	"github.com/mesos/mr-redis/common/id"
	typ "github.com/mesos/mr-redis/common/types"
)

//Creator Goroutine responsible for creating a redis instance/Proc waits on channel Cchan for any work to be done
func Creator() {

	for {
		select {
		case tc := <-typ.Cchan:
			log.Printf("Received offer %v", tc)
			//Push back the offer in the offer list
			inst := tc.I
			cpu := 1
			mem := inst.Capacity

			if tc.M {
		
				//If this is a Master instance then the count (tc.C) should always be 1
				if tc.C != 1{
					inst.ExpMasters = 1
				}
				typ.OfferList.PushBack(typ.NewOffer(inst.Name+"::"+id.NewUIIDstr(), cpu, mem, true, ""))
				log.Printf("Created %d master offers for Instance %v", tc.C, inst.Name)

			} else {

				//Create slaves only if the master is created
				if inst.Masters == inst.ExpMasters {

					p := inst.Procs[inst.Mname]

					for i := 0; i < tc.C; i++ {

						typ.OfferList.PushBack(typ.NewOffer(inst.Name+"::"+id.NewUIIDstr(), cpu, mem, false, p.IP+":"+p.Port))
					}

				}
				log.Printf("Created %d slave offers for Instance %v", tc.C, inst.Name)
			}

			break
		}
	}
}
