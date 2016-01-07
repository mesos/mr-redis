package types

//A common datatype that implements an in memeory structure and methods for remembering the Instances

type InMem struct {
	I map[string]*Instance //Map of instances
}

func NewInMem() *InMem {
	inMem := &InMem{}
	inMem.I = make(map[string]*Instance)
	return inMem
}

func (inMem *InMem) IsValid(name string) bool {

	_, ok := inMem.I[name]

	return ok
}

//Use add to add a new Instance entry in the inmemory, throws error if the element already exist
func (inMem *InMem) Add(name string, instance *Instance) (bool, error) {

	if inMem.IsValid(name) == true {
		return false, nil
	}

	inMem.I[name] = instance

	return true, nil
}

//use this to update an existing value, throws error otherwise
func (inMem *InMem) Update(name string, instance *Instance) (bool, error) {
	if inMem.IsValid(name) == false {
		return false, nil
	}

	inMem.I[name] = instance
	return true, nil
}

func (inMem *InMem) Delete(name string) (bool, error) {

	if inMem.IsValid(name) == false {
		return false, nil
	}

	delete(inMem.I, name)
	return true, nil
}
