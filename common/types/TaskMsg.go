package types

type TaskMsg struct {
	MSG int
	P   *Proc
}

const (
	TASK_MSG_DESTROY    = iota //0
	TASK_MSG_MAKEMASTER        //1
	TASK_MSG_SLAVEOF
)
