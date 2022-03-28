package handlers

import (
	"github.com/stretchr/testify/assert"
	db2 "github.com/uritrejo/palermo/internal/db"
	"testing"
)

func TestNewRepository(t *testing.T) {
	db := db2.NewBasicMsgDB()
	rp := NewRepository(db)
	assert.NotNil(t, rp)
	assert.Equal(t, db, rp.msgDb)
}

func TestRepository_HandleCreateMsg(t *testing.T) {

}

func TestRepository_HandleCreateMsg_BadRequest(t *testing.T) {

}

func TestRepository_HandleCreateMsg_Conflict(t *testing.T) {

}

func TestRepository_HandleCreateMsg_UnsupportedMediaType(t *testing.T) {

}

func TestRepository_HandleRetrieveMsg(t *testing.T) {

}

func TestRepository_HandleRetrieveMsg_NotFound(t *testing.T) {

}

func TestRepository_HandleRetrieveAllMsgs(t *testing.T) {

}

func TestRepository_HandleUpdateMsg(t *testing.T) {

}

func TestRepository_HandleUpdateMsg_BadRequest(t *testing.T) {

}

func TestRepository_HandleUpdateMsg_NotFound(t *testing.T) {

}

func TestRepository_HandleUpdateMsg_UnsupportedMediaType(t *testing.T) {

}

func TestRepository_HandleDeleteMsg(t *testing.T) {

}

func TestRepository_HandleDeleteMsg_NotFound(t *testing.T) {

}
