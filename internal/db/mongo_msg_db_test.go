package db

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	// set runMongoDBTests to true to run this test suite, must have a mongo db set up at testMongoDBAddr
	runMongoDBTests    = true
	testMongoDBAddr    = "mongodb://localhost:27017"
	testDBName         = "testMsgDB"
	testCollectionName = "testMsgCollection"
)

func TestNewMongoMsgDB(t *testing.T) {
	if !runMongoDBTests {
		t.Skip("MongoDB tests are disabled")
	}
	db, err := NewMongoMsgDB(testMongoDBAddr, testDBName, testCollectionName)
	assert.Nil(t, err)
	defer db.client.Database(testDBName).Drop(context.TODO())

	assert.NotNil(t, db.client)
	assert.NotNil(t, db.msgCollection)
}

func TestMongoMsgDB_CreateGetMsg(t *testing.T) {
	if !runMongoDBTests {
		t.Skip("MongoDB tests are disabled")
	}
	db, err := NewMongoMsgDB(testMongoDBAddr, testDBName, testCollectionName)
	assert.Nil(t, err)
	defer db.client.Database(testDBName).Drop(context.TODO())

	msg1 := NewMsg("unicorn", "kayak")
	err = db.CreateMsg(msg1)
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

func TestMongoMsgDB_CreateMsg_ErrIdUnavailable(t *testing.T) {
	if !runMongoDBTests {
		t.Skip("MongoDB tests are disabled")
	}
	db, err := NewMongoMsgDB(testMongoDBAddr, testDBName, testCollectionName)
	assert.Nil(t, err)
	defer db.client.Database(testDBName).Drop(context.TODO())

	msg1 := NewMsg("fly", "this is the message")
	err = db.CreateMsg(msg1)
	assert.Nil(t, err)

	msg2 := NewMsg("fly", "other message")
	err = db.CreateMsg(msg2)
	assert.NotNil(t, err)
	assert.IsType(t, ErrIdUnavailable{}, err)
}

func TestMongoMsgDB_GetMsg_ErrMsgNotFound(t *testing.T) {
	if !runMongoDBTests {
		t.Skip("MongoDB tests are disabled")
	}
	db, err := NewMongoMsgDB(testMongoDBAddr, testDBName, testCollectionName)
	assert.Nil(t, err)
	defer db.client.Database(testDBName).Drop(context.TODO())

	msg, err := db.GetMsg("potato")
	assert.Nil(t, msg)
	assert.NotNil(t, err)
	assert.IsType(t, ErrMsgNotFound{}, err)

	msg, err = db.GetMsg("88888")
	assert.Nil(t, msg)
	assert.NotNil(t, err)
	assert.IsType(t, ErrMsgNotFound{}, err)
}

func TestMongoMsgDB_GetAllMsgs(t *testing.T) {
	if !runMongoDBTests {
		t.Skip("MongoDB tests are disabled")
	}
	db, err := NewMongoMsgDB(testMongoDBAddr, testDBName, testCollectionName)
	assert.Nil(t, err)
	defer db.client.Database(testDBName).Drop(context.TODO())

	msgs, err := db.GetAllMsgs()
	assert.Equal(t, 0, len(msgs))

	msg1 := NewMsg("horse", "caballo")
	err = db.CreateMsg(msg1)
	assert.Nil(t, err)

	msgs, err = db.GetAllMsgs()
	assert.Equal(t, 1, len(msgs))

	msg2 := NewMsg("quertyuiop", "zxcvbnm,")
	err = db.CreateMsg(msg2)
	assert.Nil(t, err)

	msgs, err = db.GetAllMsgs()
	assert.Equal(t, 2, len(msgs))

	for _, id := range []string{"horse", "quertyuiop"} {
		idFound := false
		for _, m := range msgs {
			if id == m.Id {
				idFound = true
				break
			}
		}
		assert.True(t, idFound)
	}
}

func TestMongoMsgDB_UpdateMsg(t *testing.T) {
	if !runMongoDBTests {
		t.Skip("MongoDB tests are disabled")
	}
	db, err := NewMongoMsgDB(testMongoDBAddr, testDBName, testCollectionName)
	assert.Nil(t, err)
	defer db.client.Database(testDBName).Drop(context.TODO())

	msg := NewMsg("pegasus", "123456")
	err = db.CreateMsg(msg)
	assert.Nil(t, err)

	newMsg := NewMsg(msg.Id, "iAmGroot")
	err = db.UpdateMsg(newMsg)
	assert.Nil(t, err)

	retMsg, err := db.GetMsg(msg.Id)
	assert.Nil(t, err)

	assert.Equal(t, newMsg.Content, retMsg.Content)
	assert.Equal(t, newMsg.IsPalindrome, retMsg.IsPalindrome)
}

func TestMongoMsgDB_UpdateMsg_ErrMsgNotFound(t *testing.T) {
	if !runMongoDBTests {
		t.Skip("MongoDB tests are disabled")
	}
	db, err := NewMongoMsgDB(testMongoDBAddr, testDBName, testCollectionName)
	assert.Nil(t, err)
	defer db.client.Database(testDBName).Drop(context.TODO())

	newMsg := NewMsg("nonexistent", "iAmGroot")
	err = db.UpdateMsg(newMsg)
	assert.NotNil(t, err)
	assert.IsType(t, ErrMsgNotFound{}, err)
}

func TestMongoMsgDB_DeleteMsg(t *testing.T) {
	if !runMongoDBTests {
		t.Skip("MongoDB tests are disabled")
	}

	// delete, then delete again and make sure it returns ErrMsgNotFound

	db, err := NewMongoMsgDB(testMongoDBAddr, testDBName, testCollectionName)
	assert.Nil(t, err)
	defer db.client.Database(testDBName).Drop(context.TODO())

	msg := NewMsg("helloMsg", "goodbye my friend")
	err = db.CreateMsg(msg)
	assert.Nil(t, err)

	err = db.DeleteMsg(msg.Id)
	assert.Nil(t, err)

	// try to delete again now that it's deleted
	err = db.DeleteMsg(msg.Id)
	assert.NotNil(t, err)
	assert.IsType(t, ErrMsgNotFound{}, err)
}
