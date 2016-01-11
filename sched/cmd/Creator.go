package cmd

import (
	//"container/list"
	"log"
	//	"time"

	"../../common/id"
	typ "../../common/types"
)

func Creator() {

	for {
		select {
		/*
			case <-time.After(1 * time.Second):
				log.Printf("Creator Heart Beat")
				break
		*/

		case inst := <-typ.Cchan:
			log.Printf("Recived offer %v", inst)

			//Push back the offer in the offer list
			cpu := 1
			mem := inst.Capacity

			moffers := inst.ExpMasters - inst.Masters //Number of offers for creating masters
			soffers := inst.ExpSlaves - inst.Slaves   //Number of offers to be selected for the lsaves

			for i := 0; i < moffers; i++ {

				typ.OfferList.PushBack(typ.NewOffer(inst.Name+"::"+id.NewUIIDstr(), cpu, mem, true))
			}

			for i := 0; i < soffers; i++ {

				typ.OfferList.PushBack(typ.NewOffer(inst.Name+"::"+id.NewUIIDstr(), cpu, mem, false))
			}

			log.Printf("Created %d master offers and %d slave offers capcacity=%d", moffers, soffers, mem)

			break
		}
	}
}
