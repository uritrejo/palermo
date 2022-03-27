package db

// ErrMsgNotFound is used when a message with the Id provided is not found
type ErrMsgNotFound struct {}

func (e ErrMsgNotFound) Error() string {
	return "There was no message associated with the ID provided"
}

func IsErrMsgNotFound(err error) bool {
	_, isErrMsgNotFound := err.(ErrMsgNotFound)
	return isErrMsgNotFound
}

// ErrIdUnavailable is used when the Id provided for a new message is already in use
type ErrIdUnavailable struct {}

func (e ErrIdUnavailable) Error() string {
	return "The ID  provided is already in use"
}

func IsErrIdUnavailable(err error) bool {
	_, isErrIdUnavailable := err.(ErrIdUnavailable)
	return isErrIdUnavailable
}
