package cmd

import (
	//"container/list"
	"log"
	//	"time"

	"github.com/mesos/mr-redis/common/id"
	typ "github.com/mesos/mr-redis/common/types"
)

func Creator() {

	for {
		select {
		/*
			case <-time.After(1 * time.Second):
				log.Printf("Creator Heart Beat")
				break
		*/

		case tc := <-typ.Cchan:
			log.Printf("Received offer %v", tc)
			//Push back the offer in the offer list
			inst := tc.I
			cpu := 1
			mem := inst.Capacity

			if tc.M {

				for i := 0; i < tc.C; i++ {

					typ.OfferList.PushBack(typ.NewOffer(inst.Name+"::"+id.NewUIIDstr(), cpu, mem, true, ""))
				}
				log.Printf("Created %d master offers for Instnace %v", tc.C, inst.Name)

			} else {

				//Create slaves only if the master is created
				if inst.Masters == inst.ExpMasters {

					p := inst.Procs[inst.Mname]

					for i := 0; i < tc.C; i++ {

						typ.OfferList.PushBack(typ.NewOffer(inst.Name+"::"+id.NewUIIDstr(), cpu, mem, false, p.IP+":"+p.Port))
					}

				}
				log.Printf("Created %d slave offers for Instnace %v", tc.C, inst.Name)
			}

			break
		}
	}
}
