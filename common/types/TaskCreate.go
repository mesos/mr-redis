package types

type TaskCreate struct {
	M bool //Is this a master or slave
	I *Instance
	C int //count of number of instance to be created
}

func NewTaskCreate(m bool, i *Instance, c int) TaskCreate {
	return TaskCreate{M: m, I: i, C: c}
}

func CreateMaster(i *Instance) TaskCreate {
	return NewTaskCreate(true, i, 1)
}

func CreateSlaves(i *Instance, c int) TaskCreate {
	return NewTaskCreate(false, i, c)
}
