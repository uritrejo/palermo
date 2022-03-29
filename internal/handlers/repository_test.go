package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/uritrejo/palermo/internal/db"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewRepository(t *testing.T) {
	basicDb := db.NewBasicMsgDB()
	rp := NewRepository(basicDb)
	assert.NotNil(t, rp)
	assert.Equal(t, basicDb, rp.msgDb)
}

func TestHandleReqErr(t *testing.T) {
	rr := httptest.NewRecorder()
	handleReqErr(rr, "base error", http.StatusNotFound, "secret msg")
	assert.Contains(t, rr.Body.String(), "base error")
	assert.NotContains(t, rr.Body.String(), "secret msg")
}

func TestRepository_HandleCreateMsg(t *testing.T) {
	rp := NewRepository(db.NewBasicMsgDB())

	msg := `{"id": "unicorn", "content": "kayak"}`
	req := httptest.NewRequest("POST", "/v1/createMsg", bytes.NewReader([]byte(msg)))
	req.Header.Set("content-type", "application/json")
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(rp.HandleCreateMsg)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestRepository_HandleCreateMsg_BadRequest(t *testing.T) {
	rp := NewRepository(db.NewBasicMsgDB())

	// json msg missing closing }
	msg := `{"id": "unicorn", "content": "kayak"`
	req := httptest.NewRequest("POST", "/v1/createMsg", bytes.NewReader([]byte(msg)))
	req.Header.Set("content-type", "application/json")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(rp.HandleCreateMsg)

	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	// msg id is empty
	msg = `{"id": "", "content": "kayak"}`
	req = httptest.NewRequest("POST", "/v1/createMsg", bytes.NewReader([]byte(msg)))
	req.Header.Set("content-type", "application/json")

	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestRepository_HandleCreateMsg_Conflict(t *testing.T) {
	rp := NewRepository(db.NewBasicMsgDB())

	msg := `{"id": "unicorn", "content": "kayak"}`
	req := httptest.NewRequest("POST", "/v1/createMsg", bytes.NewReader([]byte(msg)))
	req.Header.Set("content-type", "application/json")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(rp.HandleCreateMsg)

	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	msg = `{"id": "unicorn", "content": "other message"}`
	req = httptest.NewRequest("POST", "/v1/createMsg", bytes.NewReader([]byte(msg)))
	req.Header.Set("content-type", "application/json")

	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusConflict, rr.Code)
}

func TestRepository_HandleCreateMsg_UnsupportedMediaType(t *testing.T) {
	rp := NewRepository(db.NewBasicMsgDB())

	msg := `{"id": "unicorn", "content": "kayak"}`
	req := httptest.NewRequest("POST", "/v1/createMsg", bytes.NewReader([]byte(msg)))
	req.Header.Set("content-type", "text/xml")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(rp.HandleCreateMsg)

	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusUnsupportedMediaType, rr.Code)
}

func TestRepository_HandleRetrieveMsg(t *testing.T) {
	basicDb := db.NewBasicMsgDB()
	rp := NewRepository(basicDb)

	err := basicDb.CreateMsg(db.NewMsg("potato", "11/11/11"))
	assert.Nil(t, err)

	req := httptest.NewRequest("GET", "/v1/retrieveMsg/potato", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "potato"})

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(rp.HandleRetrieveMsg)

	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	var msg db.Msg
	err = json.NewDecoder(rr.Body).Decode(&msg)
	assert.Nil(t, err)

	assert.Equal(t, "potato", msg.Id)
	assert.Equal(t, "11/11/11", msg.Content)
	assert.True(t, msg.IsPalindrome)
}

func TestRepository_HandleRetrieveMsg_NotFound(t *testing.T) {
	basicDb := db.NewBasicMsgDB()
	rp := NewRepository(basicDb)

	req := httptest.NewRequest("GET", "/v1/retrieveMsg/potato", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "888000"})

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(rp.HandleRetrieveMsg)

	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestRepository_HandleRetrieveAllMsgs(t *testing.T) {
	basicDb := db.NewBasicMsgDB()
	rp := NewRepository(basicDb)

	err := basicDb.CreateMsg(db.NewMsg("potato", "le message"))
	assert.Nil(t, err)
	err = basicDb.CreateMsg(db.NewMsg("banana", "anana"))
	assert.Nil(t, err)
	err = basicDb.CreateMsg(db.NewMsg("lemon", "i am a fruit"))
	assert.Nil(t, err)

	req := httptest.NewRequest("GET", "/v1/retrieveAllMsgs", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(rp.HandleRetrieveAllMsgs)

	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	resp := make(map[string][]db.Msg)
	err = json.NewDecoder(rr.Body).Decode(&resp)
	assert.Nil(t, err)
	allMsgs, exists := resp["messages"]
	assert.True(t, exists)

	for _, id := range []string{} {
		found := false
		for _, msg := range allMsgs {
			if id == msg.Id {
				found = true
				break
			}
		}
		assert.True(t, found)
	}
}

func TestRepository_HandleRetrieveAllMsgs_Empty(t *testing.T) {
	basicDb := db.NewBasicMsgDB()
	rp := NewRepository(basicDb)

	req := httptest.NewRequest("GET", "/v1/retrieveAllMsgs", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(rp.HandleRetrieveAllMsgs)

	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	resp := make(map[string][]db.Msg)
	err := json.NewDecoder(rr.Body).Decode(&resp)
	assert.Nil(t, err)
	allMsgs, exists := resp["messages"]
	assert.True(t, exists)
	assert.EqualValues(t, 0, len(allMsgs))
}

func TestRepository_HandleUpdateMsg(t *testing.T) {
	basicDb := db.NewBasicMsgDB()
	rp := NewRepository(basicDb)

	err := basicDb.CreateMsg(db.NewMsg("pony", "dskahfbgkalfjsd[a"))
	assert.Nil(t, err)

	msg := `{"id": "pony", "content": "chocolate123"}`
	req := httptest.NewRequest("POST", "/v1/updateMsg/pony", bytes.NewReader([]byte(msg)))
	req.Header.Set("content-type", "application/json")
	req = mux.SetURLVars(req, map[string]string{"id": "pony"})
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(rp.HandleUpdateMsg)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestRepository_HandleUpdateMsg_BadRequest(t *testing.T) {
	basicDb := db.NewBasicMsgDB()
	rp := NewRepository(basicDb)

	err := basicDb.CreateMsg(db.NewMsg("pony", "dskahfbgkalfjsd[a"))
	assert.Nil(t, err)

	msg := `{"id": "pony", "content": "chocolate123"}`
	req := httptest.NewRequest("POST", "/v1/updateMsg/pony", bytes.NewReader([]byte(msg)))
	req.Header.Set("content-type", "application/json")
	req = mux.SetURLVars(req, map[string]string{"id": "notPony"}) // id doesn't match the one in body msg
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(rp.HandleUpdateMsg)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestRepository_HandleUpdateMsg_NotFound(t *testing.T) {
	basicDb := db.NewBasicMsgDB()
	rp := NewRepository(basicDb)

	// we never create the entry for 'pony'

	msg := `{"id": "pony", "content": "chocolate123"}`
	req := httptest.NewRequest("POST", "/v1/updateMsg/pony", bytes.NewReader([]byte(msg)))
	req.Header.Set("content-type", "application/json")
	req = mux.SetURLVars(req, map[string]string{"id": "pony"})
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(rp.HandleUpdateMsg)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestRepository_HandleUpdateMsg_UnsupportedMediaType(t *testing.T) {
	basicDb := db.NewBasicMsgDB()
	rp := NewRepository(basicDb)

	err := basicDb.CreateMsg(db.NewMsg("pony", "dskahfbgkalfjsd[a"))
	assert.Nil(t, err)

	msg := `{"id": "pony", "content": "chocolate123"}`
	req := httptest.NewRequest("POST", "/v1/updateMsg/pony", bytes.NewReader([]byte(msg)))
	req.Header.Set("content-type", "text/xml") // xml
	req = mux.SetURLVars(req, map[string]string{"id": "pony"})
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(rp.HandleUpdateMsg)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnsupportedMediaType, rr.Code)
}

func TestRepository_HandleDeleteMsg(t *testing.T) {
	basicDb := db.NewBasicMsgDB()
	rp := NewRepository(basicDb)

	err := basicDb.CreateMsg(db.NewMsg("elephant", "they don't live in the forest"))
	assert.Nil(t, err)

	req := httptest.NewRequest("GET", "/v1/deleteMsg/elephant", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "elephant"})

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(rp.HandleDeleteMsg)

	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestRepository_HandleDeleteMsg_NotFound(t *testing.T) {
	basicDb := db.NewBasicMsgDB()
	rp := NewRepository(basicDb)

	// no entry with id 'elephant'

	req := httptest.NewRequest("GET", "/v1/deleteMsg/elephant", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "elephant"})

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(rp.HandleDeleteMsg)

	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusNotFound, rr.Code)
}
