package handlers

import (
	log "github.com/sirupsen/logrus"
	"net/http"

)

// todo: consider moving this into the app package


type Repository struct {
	//db
}

func NewRepository() *Repository {
	return &Repository{

	}
}

func (rp *Repository) HandleHome(w http.ResponseWriter, r *http.Request) {
	//log.Info("HandleHome request made")

	_, err := w.Write([]byte("Good night slothopie <3"))
	if err != nil {
		log.Error("Error writing msg to writer: ", err.Error())
	}
}

func (rp *Repository) HandleCreateMsg(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("create msg"))
	if err != nil {
		log.Error("Error writing msg to writer: ", err.Error())
	}
}

func (rp *Repository) HandleRetrieveMsg(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("retrieved msg"))
	if err != nil {
		log.Error("Error writing msg to writer: ", err.Error())
	}
}

func (rp *Repository) HandleUpdateMsg(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("Update Msg"))
	if err != nil {
		log.Error("Error writing msg to writer: ", err.Error())
	}
}

func (rp *Repository) HandleDeleteMsg(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("Delete Msg"))
	if err != nil {
		log.Error("Error writing msg to writer: ", err.Error())
	}
}

func (rp *Repository) HandleListMsgs(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("List Msgs"))
	if err != nil {
		log.Error("Error writing msg to writer: ", err.Error())
	}
}
