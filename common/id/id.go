package id

import (
	"github.com/nu7hatch/gouuid"
)

//NewUUID This creates a new UUID for the tasks
func NewUUID() (*uuid.UUID, error) {

	rUUID, err := uuid.NewV4()

	if err != nil {
		return nil, err
	}

	return rUUID, nil

}

//NewUIIDstr Simply convert the uid to a string
func NewUIIDstr() string {
	id, err := NewUUID()
	if err != nil {
		return ""
	}
	return id.String()
}
