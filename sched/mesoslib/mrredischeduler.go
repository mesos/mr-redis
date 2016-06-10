package mesoslib

import (
	"fmt"
	"log"
	"time"

	"github.com/gogo/protobuf/proto"
	mesos "github.com/mesos/mesos-go/mesosproto"
	util "github.com/mesos/mesos-go/mesosutil"
	sched "github.com/mesos/mesos-go/scheduler"

	typ "github.com/mesos/mr-redis/common/types"
)

type MrRedisScheduler struct {
	executor *mesos.ExecutorInfo
}

func NewMrRedisScheduler(exec *mesos.ExecutorInfo) *MrRedisScheduler {

	return &MrRedisScheduler{executor: exec}
}

func (S *MrRedisScheduler) Registered(driver sched.SchedulerDriver, frameworkId *mesos.FrameworkID, masterInfo *mesos.MasterInfo) {
	log.Printf("MrRedis Registered %v", frameworkId)

	FwIDKey := typ.ETC_CONF_DIR + "/FrameworkID"
	typ.Gdb.Set(FwIDKey, frameworkId.GetValue())
	FwTstamp := typ.ETC_CONF_DIR + "/RegisteredAt"
	typ.Gdb.Set(FwTstamp, time.Now().String())
}

func (S *MrRedisScheduler) Reregistered(driver sched.SchedulerDriver, masterInfo *mesos.MasterInfo) {
	log.Printf("MrRedis Re-registered")
	FwTstamp := typ.ETC_CONF_DIR + "/RegisteredAt"
	typ.Gdb.Set(FwTstamp, time.Now().String())
}
func (S *MrRedisScheduler) Disconnected(sched.SchedulerDriver) {
	log.Printf("MrRedis Disconnected")
}

func (S *MrRedisScheduler) ResourceOffers(driver sched.SchedulerDriver, offers []*mesos.Offer) {

	//No work to do so reject all the offers we just recived
	offer_count := typ.OfferList.Len()
	if offer_count <= 0 {
		//Reject the offers nothing to do now
		ids := make([]*mesos.OfferID, len(offers))
		for i, offer := range offers {
			ids[i] = offer.Id
		}
		driver.LaunchTasks(ids, []*mesos.TaskInfo{}, &mesos.Filters{})
		//log.Printf("No task to peform reject all the offer")
		return
	}

	//We have some task and should consume the offers sent by the masters
	//Pick one task and check if any of the offer is suitable

	//Loop throught he offers
	for _, offer := range offers {

		cpuResources := util.FilterResources(offer.Resources, func(res *mesos.Resource) bool {
			return res.GetName() == "cpus"
		})
		cpus := 0.0
		for _, res := range cpuResources {
			cpus += res.GetScalar().GetValue()
		}

		memResources := util.FilterResources(offer.Resources, func(res *mesos.Resource) bool {
			return res.GetName() == "mem"
		})
		mems := 0.0
		for _, res := range memResources {
			mems += res.GetScalar().GetValue()
		}

		log.Printf("Recived Offer with CPU=%v MEM=%v OfferID=%v", cpus, mems, offer.Id.GetValue())
		var tasks []*mesos.TaskInfo

		//Loop through the tasks
		for tsk_ele := typ.OfferList.Front(); tsk_ele != nil; {

			tsk := tsk_ele.Value.(typ.Offer)
			tskCpu_float := float64(tsk.Cpu)
			tskMem_float := float64(tsk.Mem)

			var tmp_data []byte

			if tsk.IsMaster {
				tmp_data = []byte(fmt.Sprintf("%d Master", tsk.Mem))
			} else {
				tmp_data = []byte(fmt.Sprintf("%d SlaveOf %s", tsk.Mem, tsk.MasterIpPort))
			}

			if cpus >= tskCpu_float && mems >= tskMem_float {
				tsk_id := &mesos.TaskID{Value: proto.String(tsk.Taskname)}
				mesos_tsk := &mesos.TaskInfo{
					Name:     proto.String(tsk.Taskname),
					TaskId:   tsk_id,
					SlaveId:  offer.SlaveId,
					Executor: S.executor,
					Resources: []*mesos.Resource{
						util.NewScalarResource("cpus", tskCpu_float),
						util.NewScalarResource("mem", tskMem_float),
					},
					Data: tmp_data,
				}
				mems -= tskMem_float
				cpus -= tskCpu_float

				current_task := tsk_ele
				tsk_ele = tsk_ele.Next()
				typ.OfferList.Remove(current_task)
				tasks = append(tasks, mesos_tsk)

			} else {
				tsk_ele = tsk_ele.Next()
			}
			//Check if this task is suitable for this offer
		}
		driver.LaunchTasks([]*mesos.OfferID{offer.Id}, tasks, &mesos.Filters{})
		log.Printf("Launched %d tasks from this offer", len(tasks))
	}
	log.Printf("MrRedis Recives offer")
}

func (S *MrRedisScheduler) StatusUpdate(driver sched.SchedulerDriver, status *mesos.TaskStatus) {

	var ts typ.TaskUpdate
	ts.Name = status.GetTaskId().GetValue()
	ts.State = status.GetState().String()
	ts.Data = status.GetData()
	log.Printf("MrRedis Task Update recived")
	log.Printf("Status=%v", ts)

	//Send it across to the channel to maintainer
	typ.Mchan <- &ts

}

func (S *MrRedisScheduler) OfferRescinded(_ sched.SchedulerDriver, oid *mesos.OfferID) {
	log.Printf("offer rescinded: %v", oid)
}

func (S *MrRedisScheduler) FrameworkMessage(_ sched.SchedulerDriver, eid *mesos.ExecutorID, sid *mesos.SlaveID, msg string) {
	log.Printf("framework message from executor %q slave %q: %q", eid, sid, msg)
}

func (S *MrRedisScheduler) SlaveLost(_ sched.SchedulerDriver, sid *mesos.SlaveID) {
	log.Printf("slave lost: %v", sid)
}

func (S *MrRedisScheduler) ExecutorLost(_ sched.SchedulerDriver, eid *mesos.ExecutorID, sid *mesos.SlaveID, code int) {
	log.Printf("executor %q lost on slave %q code %d", eid, sid, code)
}

func (S *MrRedisScheduler) Error(_ sched.SchedulerDriver, err string) {
	log.Printf("Scheduler received error:", err)
}
