package types

//TaskMsg Type that is used to communicate with Destroyer
type TaskMsg struct {
	MSG int
	P   *Proc
}

const (
	TASK_MSG_DESTROY    = iota //0 TASK_MSG_DESTROY enum of Message tpe starts with 0
	TASK_MSG_MAKEMASTER        //1
	TASK_MSG_SLAVEOF
)
