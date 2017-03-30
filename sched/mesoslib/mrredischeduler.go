package mesoslib

import (
	"fmt"
	"log"
	"time"
	"syscall"

	"github.com/gogo/protobuf/proto"
	mesos "github.com/mesos/mesos-go/mesosproto"
	util "github.com/mesos/mesos-go/mesosutil"
	sched "github.com/mesos/mesos-go/scheduler"

	typ "github.com/mesos/mr-redis/common/types"
)

//MrRedisScheduler scheudler struct
type MrRedisScheduler struct {
	executor *mesos.ExecutorInfo
}

//NewMrRedisScheduler Constructor
func NewMrRedisScheduler(exec *mesos.ExecutorInfo) *MrRedisScheduler {

	return &MrRedisScheduler{executor: exec}
}

//Registered Scheduler register call back initializes the timestamp and framework id
func (S *MrRedisScheduler) Registered(driver sched.SchedulerDriver, frameworkID *mesos.FrameworkID, masterInfo *mesos.MasterInfo) {
	log.Printf("MrRedis Registered %v", frameworkID)
	typ.IsRegistered = true
	FwIDKey := typ.ETC_CONF_DIR + "/FrameworkID"
	typ.Gdb.Set(FwIDKey, frameworkID.GetValue())
	FwTstamp := typ.ETC_CONF_DIR + "/RegisteredAt"
	typ.Gdb.Set(FwTstamp, time.Now().String())
}

//Reregistered re-register call back simply updates the timestamp
func (S *MrRedisScheduler) Reregistered(driver sched.SchedulerDriver, masterInfo *mesos.MasterInfo) {
	log.Printf("MrRedis Re-registered")
	FwTstamp := typ.ETC_CONF_DIR + "/RegisteredAt"
	typ.Gdb.Set(FwTstamp, time.Now().String())
}

//Disconnected Not implemented call back
func (S *MrRedisScheduler) Disconnected(sched.SchedulerDriver) {
	log.Printf("MrRedis Disconnected")
}

//ResourceOffers The moment we recive some offers we loop through the OfferList (container/list)
//see if we have any task that will fit this offers being sent
func (S *MrRedisScheduler) ResourceOffers(driver sched.SchedulerDriver, offers []*mesos.Offer) {

	//No work to do so reject all the offers we just received
	offerCount := typ.OfferList.Len()
	if offerCount <= 0 || typ.IsRegistered == false {
		//Reject the offers nothing to do now or even the Framework registration isnt completed
		ids := make([]*mesos.OfferID, len(offers))
		for i, offer := range offers {
			ids[i] = offer.Id
		}
		driver.LaunchTasks(ids, []*mesos.TaskInfo{}, &mesos.Filters{})
		//log.Printf("No task to peform reject all the offer")
		if typ.IsRegistered == false {
			log.Printf("Rejecting Offers Framework Not registered yet")
		}
		return
	}

	//We have some task and should consume the offers sent by the masters
	//Pick one task and check if any of the offer is suitable

	//Loop thought he offers
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

		log.Printf("Received Offer with CPU=%v MEM=%v OfferID=%v", cpus, mems, offer.Id.GetValue())
		var tasks []*mesos.TaskInfo

		//Loop through the tasks
		for tskEle := typ.OfferList.Front(); tskEle != nil; {

			tsk := tskEle.Value.(typ.Offer)
			tskCPUFloat := float64(tsk.Cpu)
			tskMemFloat := float64(tsk.Mem)

			var tmpData []byte

			if tsk.IsMaster {
				tmpData = []byte(fmt.Sprintf("%d Master", tsk.Mem))
			} else {
				tmpData = []byte(fmt.Sprintf("%d SlaveOf %s", tsk.Mem, tsk.MasterIpPort))
			}

			if cpus >= tskCPUFloat && mems >= tskMemFloat && typ.Agents.Canfit(offer.SlaveId.GetValue(), tsk.Name, tsk.DValue) {
				tskID := &mesos.TaskID{Value: proto.String(tsk.Taskname)}
				mesosTsk := &mesos.TaskInfo{
					Name:     proto.String(tsk.Taskname),
					TaskId:   tskID,
					SlaveId:  offer.SlaveId,
					Executor: S.executor,
					Resources: []*mesos.Resource{
						util.NewScalarResource("cpus", tskCPUFloat),
						util.NewScalarResource("mem", tskMemFloat),
					},
					Data: tmpData,
				}
				mems -= tskMemFloat
				cpus -= tskCPUFloat

				currentTask := tskEle
				tskEle = tskEle.Next()
				typ.OfferList.Remove(currentTask)
				tasks = append(tasks, mesosTsk)
				typ.Agents.Add(offer.SlaveId.GetValue(), tsk.Name, 1)

			} else {
				tskEle = tskEle.Next()
			}
			//Check if this task is suitable for this offer
		}
		driver.LaunchTasks([]*mesos.OfferID{offer.Id}, tasks, &mesos.Filters{})
		log.Printf("Launched %d tasks from this offer", len(tasks))
	}
	log.Printf("MrRedis Receives offer")
}

//StatusUpdate Simply recives the update and passes it to the Maintainer goroutine
func (S *MrRedisScheduler) StatusUpdate(driver sched.SchedulerDriver, status *mesos.TaskStatus) {

	var ts typ.TaskUpdate
	ts.Name = status.GetTaskId().GetValue()
	ts.State = status.GetState().String()
	ts.SlaveId = status.SlaveId.GetValue()
	ts.Data = status.GetData()
	log.Printf("MrRedis Task Update received")
	log.Printf("Status=%v", ts)

	//Send it across to the channel to maintainer
	typ.Mchan <- &ts

}

//OfferRescinded Not implemented
func (S *MrRedisScheduler) OfferRescinded(_ sched.SchedulerDriver, oid *mesos.OfferID) {
	log.Printf("offer rescinded: %v", oid)
}

//FrameworkMessage not implemented
func (S *MrRedisScheduler) FrameworkMessage(_ sched.SchedulerDriver, eid *mesos.ExecutorID, sid *mesos.SlaveID, msg string) {
	log.Printf("framework message from executor %q slave %q: %q", eid, sid, msg)
}

//SlaveLost Not implemented
func (S *MrRedisScheduler) SlaveLost(_ sched.SchedulerDriver, sid *mesos.SlaveID) {
	log.Printf("slave lost: %v", sid)
}

//ExecutorLost Not implemented
func (S *MrRedisScheduler) ExecutorLost(_ sched.SchedulerDriver, eid *mesos.ExecutorID, sid *mesos.SlaveID, code int) {
	log.Printf("executor %q lost on slave %q code %d", eid, sid, code)
}

//Error Not implemeted
func (S *MrRedisScheduler) Error(_ sched.SchedulerDriver, err string) {
	log.Printf("Scheduler received error:%v", err)
        if err == "Framework has been removed"{
		FwIDKey := typ.ETC_CONF_DIR + "/FrameworkID"
        	typ.Gdb.Set(FwIDKey, "") // This means that when we start up *next time* we will have no framework ID and will register for a new one.
                syscall.Kill(syscall.Getpid(), syscall.SIGINT)
        }
}
