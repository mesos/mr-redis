package mesoslib

import (
	"log"

	mesos "github.com/mesos/mesos-go/mesosproto"
	util "github.com/mesos/mesos-go/mesosutil"
	sched "github.com/mesos/mesos-go/scheduler"
	"golang.org/x/net/context"
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
