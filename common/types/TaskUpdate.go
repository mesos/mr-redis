package types

//TaskUpdate type used to community with Maintainer goroutine
type TaskUpdate struct {
	Name  string
	State string
	Data  []byte
}
