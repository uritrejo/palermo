package db

// MsgDB exposes the functionality to create, delete, update and retrieve a message
type MsgDB interface {
	GetMsg(id string) (*Msg, error)
	GetAllMsgs() ([]*Msg, error)
	CreateMsg(msg *Msg) error
	UpdateMsg(id string, msg *Msg) error
	DeleteMsg(id string) error
}

// ErrMsgNotFound is used when a message with the Id provided is not found
type ErrMsgNotFound struct {}

func (e *ErrMsgNotFound) Error() string {
	return "There was no message associated with the ID provided"
}

// ErrIdUnavailable is used when the Id provided for a new message is already in use
type ErrIdUnavailable struct {}

func (e *ErrIdUnavailable) Error() string {
	return "The ID  provided is already in use"
}
