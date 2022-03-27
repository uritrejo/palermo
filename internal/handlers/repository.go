package handlers

import (
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
)

// todo: consider moving this into the app package


type Repository struct {
	//db
}

// todo: receive the db as an arg
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

// POST
func (rp *Repository) HandleCreateMsg(w http.ResponseWriter, r *http.Request) {
	// try to unmarshall json
	// then add it to the database


	_, err := w.Write([]byte("create msg"))
	if err != nil {
		log.Error("Error writing msg to writer: ", err.Error())
	}
}

func (rp *Repository) HandleListMsgs(w http.ResponseWriter, r *http.Request) {
	// todo: basically the same as the one below

	_, err := w.Write([]byte("List Msgs"))
	if err != nil {
		log.Error("Error writing msg to writer: ", err.Error())
	}
}

func (rp *Repository) HandleRetrieveMsg(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]

	log.Info("The key was: ", key)

	// error 404

	// todo: encode in json and write it to the w
	// json.NewEncoder(w).Encode(article)
	_, err := w.Write([]byte("retrieved msg"))
	if err != nil {
		log.Error("Error writing msg to writer: ", err.Error())
	}
}

func (rp *Repository) HandleUpdateMsg(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	// use strings.trim to reassign the messageid

	// here you can accept both adding the ID or not in the body
	// but if it's added, then it has to match, else return a 404 (bref you can't accept the id to be changed)
	// so you don't have to duplicate it

	_, err := w.Write([]byte("Update Msg"))
	if err != nil {
		log.Error("Error writing msg to writer: ", err.Error())
	}
}

func (rp *Repository) HandleDeleteMsg(w http.ResponseWriter, r *http.Request) {

	// get the id, then delete it
	// this can remove it while iterating
	// Articles = append(Articles[:index], Articles[index+1:]...)

	_, err := w.Write([]byte("Delete Msg"))
	if err != nil {
		log.Error("Error writing msg to writer: ", err.Error())
	}
}
