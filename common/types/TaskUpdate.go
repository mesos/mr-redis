package types

//TaskUpdate type used to communicate with Maintainer goroutine
type TaskUpdate struct {
	Name    string
	State   string
	SlaveId string
	Data    []byte
}
