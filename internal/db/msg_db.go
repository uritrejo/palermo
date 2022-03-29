package db

// MsgDB exposes the functionality to create, delete, update and retrieve a message
type MsgDB interface {
	// GetMsg returns the msg  if available,
	// returns ErrMsgNotFound if a msg with such id wasn't found
	GetMsg(id string) (*Msg, error)

	// GetAllMsgs will return all the messages in the DB, an empty slice if none
	GetAllMsgs() ([]*Msg, error)

	// CreateMsg will add the msg provided into the database
	// returns ErrIdUnavailable if the msg.Id is already in use
	CreateMsg(msg *Msg) error

	// UpdateMsg will update the msg stored with the provided msg.Id
	// returns ErrMsgNotFound if a msg with such id wasn't found
	UpdateMsg(msg *Msg) error

	// DeleteMsg will delete the message associated with the id provided
	// returns ErrMsgNotFound if a msg with such id wasn't found
	DeleteMsg(id string) error

	Close()
}
