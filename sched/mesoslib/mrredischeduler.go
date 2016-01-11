package mesoslib

import (
	"log"

	mesos "github.com/mesos/mesos-go/mesosproto"
	util "github.com/mesos/mesos-go/mesosutil"
	sched "github.com/mesos/mesos-go/scheduler"
	"golang.org/x/net/context"

	"../../common/types"
)

type MrRedisScheduler struct {
	executor *mesos.ExecutorInfo
}

func NewMrRedisScheduler(exec *mesos.ExecutorInfo) *MrRedisScheduler {

	return &MrRedisScheduler{executor: exec}
}

func (S *MrRedisScheduler) Registered(driver sched.SchedulerDriver, frameworkId *mesos.FrameworkID, masterInfo *mesos.MasterInfo) {
	log.Printf("MrRedis Registered")
}

func (S *MrRedisScheduler) Reregistered(driver sched.SchedulerDriver, masterInfo *mesos.MasterInfo) {
	log.Printf("MrRedis Re-registered")
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
		driver.LaunchTasks(ids, []*mesos.TaskInfo{}, &mesos.Filters{RefuseSeconds: proto.Float64(1)})
		log.Printf("No task to peform reject all the offer")
		return
	}

	//We have some task and should consume the offers sent by the masters
	//Pick one task and check if any of the offer is suitable

	//Loop throught he offers
	for i, offer := range offers {

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

		log.Printf("Recived Offer with CPU=%d MEM=%d OfferID=%v", cpus, mems, offer.Id.GetValue())
		var tasks []*mesos.TaskInfo

		//Loop through the tasks
		for tsk := typ.OfferList.Front(); tsk != nil; tsk = tsk.Next() {

			if cpus >= tsk.Cpu && mems >= tsk.Mem {
				tsk_id := &mesos.TaskID{Value: proto.String(tsk.Taskname)}
				mesos_tsk := &mesos.TaskInfo{
					Name:     proto.String(tsk.Taskname),
					TaskId:   tsk_id,
					SlaveId:  offer.SlaveId,
					Executor: S.executor,
					Resources: []*mesos.Resource{
						util.NewScalarResource("cpus", tsk.Cpu),
						util.NewScalarResource("mem", tsk.Mem),
					},
				}
				mems -= tsk.Mem
				cpus -= tsk.Cpu

				typ.OfferList.Remove(tsk)
				tasks = append(tasks, mesos_tsk)

			}
			//Check if this task is suitable for this offer
		}
		driver.LaunchTasks([]*mesos.OfferID{offer.Id}, tasks, nil)
	}
	log.Printf("MrRedis Recives offer")
}

func (S *MrRedisScheduler) StatusUpdate(driver sched.SchedulerDriver, status *mesos.TaskStatus) {

	log.Printf("MrRedis Recives offer")
}

func (S *MrRedisScheduler) OfferRescinded(_ sched.SchedulerDriver, oid *mesos.OfferID) {
	log.Errorf("offer rescinded: %v", oid)
}

func (S *MrRedisScheduler) FrameworkMessage(_ sched.SchedulerDriver, eid *mesos.ExecutorID, sid *mesos.SlaveID, msg string) {
	log.Errorf("framework message from executor %q slave %q: %q", eid, sid, msg)
}

func (S *MrRedisScheduler) SlaveLost(_ sched.SchedulerDriver, sid *mesos.SlaveID) {
	log.Errorf("slave lost: %v", sid)
}

func (S *MrRedisScheduler) ExecutorLost(_ sched.SchedulerDriver, eid *mesos.ExecutorID, sid *mesos.SlaveID, code int) {
	log.Errorf("executor %q lost on slave %q code %d", eid, sid, code)
}

func (S *MrRedisScheduler) Error(_ sched.SchedulerDriver, err string) {
	log.Errorf("Scheduler received error:", err)
}
