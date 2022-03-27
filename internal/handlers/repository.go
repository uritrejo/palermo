package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/uritrejo/palermo/internal/db"
	"net/http"
	"strings"
)

// Repository will implement the handlers for our REST API
// it will store all messages in msgDb
type Repository struct {
	msgDb db.MsgDB
}

func NewRepository(msgDb db.MsgDB) *Repository {
	return &Repository{
		msgDb: msgDb,
	}
}

func (rp *Repository) HandleCreateMsg(w http.ResponseWriter, r *http.Request) {
	var msg db.Msg
	err := json.NewDecoder(r.Body).Decode(&msg)
	if err != nil {
		handleReqErr(w, "Failed to decode body into msg object", http.StatusBadRequest, err.Error())
		return
	}

	msg.Id = strings.TrimSpace(msg.Id)
	if msg.Id == "" {
		handleReqErr(w, "Message id must not be empty", http.StatusBadRequest, "")
		return
	}

	err = rp.msgDb.CreateMsg(&msg)
	if err != nil {
		if db.IsErrIdUnavailable(err) {
			handleReqErr(w, "CreateMsg request failed, " + msg.Id + " is already in use", http.StatusConflict, err.Error())
			return
		}
		handleReqErr(w, "Unexpected error during creation of message", http.StatusInternalServerError, err.Error())
		return
	}

	log.Debug("A message was successfully created: ", msg)
}

func (rp *Repository) HandleRetrieveAllMsgs(w http.ResponseWriter, r *http.Request) {
	msgs, err := rp.msgDb.GetAllMsgs()
	if err != nil {
		handleReqErr(w, "Unexpected error during retrieval of all messages", http.StatusInternalServerError, err.Error())
		return
	}

	msgJson, err := json.Marshal(msgs)
	if err != nil {
		handleReqErr(w, "Unexpected error during marshalling of messages", http.StatusInternalServerError, err.Error())
		return
	}

	_, err = w.Write(msgJson)
	if err != nil {
		handleReqErr(w, "Unexpected error during encoding of messages into json", http.StatusInternalServerError, err.Error())
		return
	}
}

func (rp *Repository) HandleRetrieveMsg(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	msg, err := rp.msgDb.GetMsg(id)
	if err != nil {
		if db.IsErrMsgNotFound(err) {
			handleReqErr(w, "Msg with id " + id + " was not found", http.StatusNotFound, err.Error())
			return
		}
		handleReqErr(w, "Unexpected error during retrieval of message", http.StatusInternalServerError, err.Error())
		return
	}

	log.Debug("Successfully retrieved message: ", msg)

	msgJson, err := json.Marshal(msg)
	if err != nil {
		handleReqErr(w, "Unexpected error during marshalling of message", http.StatusInternalServerError, err.Error())
		return
	}

	_, err = w.Write(msgJson)
	if err != nil {
		handleReqErr(w, "Unexpected error during encoding of message into json", http.StatusInternalServerError, err.Error())
		return
	}
}

func (rp *Repository) HandleUpdateMsg(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	var msg db.Msg
	err := json.NewDecoder(r.Body).Decode(&msg)
	if err != nil {
		handleReqErr(w, "Failed to decode body into msg object", http.StatusBadRequest, err.Error())
		return
	}

	if id != msg.Id {
		handleReqErr(w, "The id in the request doesn't match the id in the msg object", http.StatusBadRequest, "")
		return
	}

	err = rp.msgDb.UpdateMsg(&msg)
	if err != nil {
		if db.IsErrMsgNotFound(err) {
			handleReqErr(w, "Msg with id " + id + " was not found", http.StatusNotFound, err.Error())
			return
		}
		handleReqErr(w, "Unexpected error during creation of message", http.StatusInternalServerError, err.Error())
		return
	}

	log.Debug("A message was successfully created: ", msg)
}

func (rp *Repository) HandleDeleteMsg(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	err := rp.msgDb.DeleteMsg(id)
	if err != nil {
		if db.IsErrMsgNotFound(err) {
			handleReqErr(w, "Msg with id " + id + " was not found", http.StatusNotFound, err.Error())
			return
		}
		handleReqErr(w, "Unexpected error during deletion of message" , http.StatusInternalServerError, err.Error())
		return
	}

	log.Debug("Successfully deleted message with id: ", id)
}

// handleReqErr logs the error and replies to the request
// baseErrorMsg is the error message that will be sent back,
// internalErrorMsg will be added to local logs
// this distinction is done to avoid exposing internal details
func handleReqErr(w http.ResponseWriter, baseErrorMsg string, code int, internalErrorMsg string) {
	log.Error(baseErrorMsg + ": " + internalErrorMsg + "; returned code " + http.StatusText(code))
	http.Error(w, baseErrorMsg, code)
}
