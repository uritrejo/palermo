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
	router.HandleFunc("/createMsg", repo.HandleCreateMsg)
	router.HandleFunc("/retrieveMsg", repo.HandleRetrieveMsg)
	router.HandleFunc("/updateMsg", repo.HandleUpdateMsg)
	router.HandleFunc("/deleteMsg", repo.HandleDeleteMsg)
	router.HandleFunc("/listMsgs", repo.HandleListMsgs)

	router.Use(handlers.LoggingMiddleware)

	return router
}
