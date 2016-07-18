package types

//A common datatype that implements an in memeory structure and methods for remembering the Instances

//InMem a structure for maintaining InMemory cache about all the instances, All the api's interactive with the DB should first check this MAP before actually reading form the DB.  for changes it should also update this map first before updating the DB
type InMem struct {
	I map[string]*Instance //Map of instances
}

//NewInMem is called when you want to initialize the In Memory cache for the first time
func NewInMem() *InMem {
	inMem := &InMem{}
	inMem.I = make(map[string]*Instance)
	return inMem
}

//IsValid A quick look up function to see if a key is available in the inmemory cache
func (inMem *InMem) IsValid(name string) bool {

	_, ok := inMem.I[name]

	return ok
}

//Add Use add to add a new Instance entry in the inmemory, throws error if the element already exist
func (inMem *InMem) Add(name string, instance *Instance) (bool, error) {

	if inMem.IsValid(name) == true {
		return false, nil
	}

	inMem.I[name] = instance

	return true, nil
}

//Update use this to update an existing value, throws error otherwise
func (inMem *InMem) Update(name string, instance *Instance) (bool, error) {
	if inMem.IsValid(name) == false {
		return false, nil
	}

	inMem.I[name] = instance
	return true, nil
}

//Delete use thsi to Delete an element from the cache
func (inMem *InMem) Delete(name string) (bool, error) {

	if inMem.IsValid(name) == false {
		return false, nil
	}

	delete(inMem.I, name)
	return true, nil
}

//Get get an Instance pointer from the cache
func (inMem *InMem) Get(name string) *Instance {

	if i, ok := inMem.I[name]; ok {
		return i
	}

	return nil
}
