package main

import (
	"errors"
	"flag"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/uritrejo/palermo/internal/db"
	"github.com/uritrejo/palermo/internal/handlers"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

const (
	defaultPort = 4422
	defaultDbType = "basic"
	defaultMongoDbPort = 27017
	defaultLogLevel = "debug"
	logFile = "palermo.log"
)

var (
	repo *handlers.Repository
)

func main() {
	// flags
	var dbType, logLevel string
	var port, mongoDbPort int
	flag.IntVar(&port, "port", defaultPort, "-port=<port>: port on which to listen and serve; the default is " + strconv.Itoa(defaultPort))
	flag.StringVar(&dbType, "dbtype", defaultDbType, "-dbtype=<type>: types are 'basic' (local memory) and 'mongodb; the default is " + defaultDbType)
	flag.IntVar(&mongoDbPort, "mongodbport", defaultMongoDbPort, "-mongodbport=<port>: port where mongo db is listening; the default is " + strconv.Itoa(defaultMongoDbPort))
	flag.StringVar(&logLevel, "loglevel", defaultLogLevel, "-loglevel=<level>: levels are info, debug, trace; default is " + defaultLogLevel)
	flag.Parse()

	closer, err := initLogger(logLevel)
	if err != nil {
		log.Warn("Failed to initialize log file: ", err.Error(), "; only console output will be available")
	}
	defer closer.Close()

	msgDb, err := initDb(dbType, mongoDbPort)
	if err != nil {
		log.Fatal("Failed to initialize database: ", err.Error())
	}

	repo = handlers.NewRepository(msgDb)

	addr := "localhost:" + strconv.Itoa(port)
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

// initDb creates the required database instance
// returns a
func initDb(dbType string, mongoDbPort int) (db.MsgDB, error) {
	var err error
	var msgDb db.MsgDB
	if dbType == "basic" {
		msgDb = db.NewBasicMsgDB()
	} else if dbType == "mongodb" {
		log.Info("Mongo DB port set to ", mongoDbPort)
		return nil, errors.New("MongoDB not implemented yet")
	} else {
		return nil, errors.New("unsupported db type: " + dbType)
	}

	log.Info("Database type set to ", dbType)
	return msgDb, err
}

// initLogger sets the log level and attempts to open a log file
// return an error and a closer that should be used to close the log file at the end of its lifetime
func initLogger(logLevel string) (io.Closer, error) {
	f, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		log.Warn("Failed to open/create log file, logging will only be enable to std out: " + err.Error())
	}
	log.SetOutput(io.MultiWriter(f, os.Stdout))

	if logLevel == "info" {
		log.SetLevel(log.InfoLevel)
	} else if logLevel == "debug" {
		log.SetLevel(log.DebugLevel)
	} else if logLevel == "trace" {
		log.SetLevel(log.TraceLevel)
	} else {
		return nil, errors.New("unsupported log level: " + logLevel)
	}

	log.Info("Log level set to ", logLevel)
	return f, nil
}

// will later take in a config obj
func router() http.Handler {
	router := mux.NewRouter()

	// handlers
	router.HandleFunc("/createMsg", repo.HandleCreateMsg).Methods("POST")
	router.HandleFunc("/listMsgs", repo.HandleListMsgs)
	router.HandleFunc("/retrieveMsg/{id}", repo.HandleRetrieveMsg)
	router.HandleFunc("/updateMsg/{id}", repo.HandleUpdateMsg).Methods("POST")
	router.HandleFunc("/deleteMsg/{id}", repo.HandleDeleteMsg)

	// middlewares
	router.Use(handlers.LoggingMiddleware)

	return router
}
