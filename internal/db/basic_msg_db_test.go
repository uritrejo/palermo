package db

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewBasicMsgDB(t *testing.T) {
	db := NewBasicMsgDB()
	assert.NotNil(t, db)
	assert.NotNil(t, db.msgs)
}

func TestBasicMsgDB_CreateGetMsg(t *testing.T) {
	db := NewBasicMsgDB()

	msg1 := NewMsg("unicorn", "kayak")
	err := db.CreateMsg(msg1)
	assert.Nil(t, err)

	msg2 := NewMsg("1234", "message")
	err = db.CreateMsg(msg2)
	assert.Nil(t, err)

	retMsg1, err := db.GetMsg("unicorn")
	assert.Nil(t, err)
	assert.Equal(t, msg1.Id, retMsg1.Id)
	assert.Equal(t, msg1.Content, retMsg1.Content)
	assert.Equal(t, msg1.IsPalindrome, retMsg1.IsPalindrome)
	assert.True(t, msg1.ModTime.Equal(retMsg1.ModTime))

	retMsg2, err := db.GetMsg("1234")
	assert.Nil(t, err)
	assert.Equal(t, msg2.Id, retMsg2.Id)
	assert.Equal(t, msg2.Content, retMsg2.Content)
	assert.Equal(t, msg2.IsPalindrome, retMsg2.IsPalindrome)
	assert.True(t, msg2.ModTime.Equal(retMsg2.ModTime))
}

func TestBasicMsgDB_CreateMsg_ErrIdUnavailable(t *testing.T) {

}

func TestBasicMsgDB_GetMsg_ErrMsgNotFound(t *testing.T) {

}

func TestBasicMsgDB_GetAllMsgs(t *testing.T) {

}

func TestBasicMsgDB_UpdateMsg(t *testing.T) {

}

func TestBasicMsgDB_UpdateMsg_ErrMsgNotFound(t *testing.T) {

}

func TestBasicMsgDB_DeleteMsg(t *testing.T) {
	// delete, then delete again and make sure it returns ErrMsgNotFound
}
