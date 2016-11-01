//Package agentstate helps us keep a mental map of all the workload that is running across different Agents.  It helps us remember what Instance is running in which node so that we can perform the load distribution or workload affinitiy more easily
package agentstate

import (
	//	"fmt"
	"sync"
)

type Inst map[string]int

//State struct is responsible for maintaining Map of Agents, this should also have methods to let the Package Implementer know that if a particular instace can be provisioned here or not
type State struct {
	lck    sync.Mutex
	Agents map[string]Inst
	Count  int
	IsSet  bool
}

//NewState basic constructor for the struct
func NewState() *State {

	var AS State
	AS.Agents = make(map[string]Inst)
	AS.Count = 0
	AS.IsSet = true
	return &AS
}

//Add adds an Instance to a particular node, if the entry is unavailable it tries to create it first
func (S *State) Add(Node string, Name string, Count int) bool {

	if !S.IsSet {
		return false
	}

	S.lck.Lock()
	defer S.lck.Unlock()

	var I Inst
	var exist bool

	if I, exist = S.Agents[Node]; exist != true {
		I = make(Inst)
		I[Name] = 0
	}
	I[Name] = I[Name] + Count
	S.Agents[Node] = I

	return true
}

//Del This removes an entry from the Map it returns false which means that such an entry itself is not available in the map
func (S *State) Del(Node string, Name string) bool {

	if !S.IsSet {
		return false
	}
	S.lck.Lock()
	defer S.lck.Unlock()

	var I Inst
	var exist bool
	var Count int

	if I, exist = S.Agents[Node]; exist != true {
		return false
	}
	if Count, exist = I[Name]; exist != true {
		return false
	}
	if Count <= 1 {
		delete(I, Name)
	} else {
		I[Name] = Count - 1
	}
	S.Agents[Node] = I
	return true
}

//InstancesRunning This will return how many workload for this instance running on this particular slave, if it returns -1 it means there is no such slave
func (S *State) InstancesRunning(Node string, Name string) int {
	if !S.IsSet {
		return -1
	}

	S.lck.Lock()
	defer S.lck.Unlock()

	var I Inst
	var exist bool
	var Count int

	if I, exist = S.Agents[Node]; exist != true {
		return 0
	}
	if Count, exist = I[Name]; exist != true {
		return 0
	}
	return Count
}

//Canfit This will tell us if a Particualr Instnace with supplied distribution value can fit in that node or not
func (S *State) Canfit(Node string, Name string, DistributionValue int) bool {

	if !S.IsSet {
		return false
	}

	return S.InstancesRunning(Node, Name) < DistributionValue
}

type NElement struct {
	NodeName string
	Count    int
}

//ListDistribution will Retrun an array of Nodes with structure
func (S *State) ListDistribution(Name string) []NElement {

	if !S.IsSet {
		return nil
	}

	var Narray []NElement

	S.lck.Lock()
	defer S.lck.Unlock()

	for k, v := range S.Agents {

		if value, Exist := v[Name]; Exist == true {
			Narray = append(Narray, NElement{k, value})
		}
	}

	return Narray
}
