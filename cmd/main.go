package main

import (
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/uritrejo/palermo/internal/handlers"
	"net/http"
	"time"
)

const (
	addr = ":4420"
	logFile = "./log/palermo-server.log"
)

var (
	repo *handlers.Repository
)

func main() {
	// todo: fix log file creation & closing
	err := initLogger()
	if err != nil {
		log.Warn("Failed to initialize log file: ", err.Error(), "; only console output will be available")
	}

	// init repo with some args once you have some
	repo = handlers.NewRepository()

	server := &http.Server{
		Handler: router(),
		Addr: addr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Info("Palermo server is listening on ", addr)
	err = server.ListenAndServe()
	if err != nil {
		log.Fatal("Failed to listen and serve: ", err.Error())
	}
}

// todo: probably will just remove the support for a file and direct output to a file from cli
func initLogger() error {
	//f, err := os.OpenFile(logFile, os.O_APPEND|os.O_APPEND|os.O_WRONLY, 0600)
	//if err != nil {
	//	log.Error("Failed to open/create log file: " + err.Error())
	//}
	//
	//log.SetOutput(io.MultiWriter(f, os.Stdout))

	log.SetLevel(log.DebugLevel)  // could also make it configurable, not for now

	return nil
}

// will later take in a config obj
func router() http.Handler {
	router := mux.NewRouter()

	// add handlers
	router.HandleFunc("/", repo.HandleHome)
	router.HandleFunc("/createMsg", repo.HandleCreateMsg).Methods("POST")
	router.HandleFunc("/listMsgs", repo.HandleListMsgs)
	// aka path variables
	router.HandleFunc("/retrieveMsg/{id}", repo.HandleRetrieveMsg)
	router.HandleFunc("/updateMsg/{id}", repo.HandleUpdateMsg).Methods("POST")
	router.HandleFunc("/deleteMsg/{id}", repo.HandleDeleteMsg)

	// todo: add a handler for all the rest to return a message
	// an example was in the gorilla example documentation

	router.Use(handlers.LoggingMiddleware)

	return router
}
