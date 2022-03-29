package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/uritrejo/palermo/internal/db"
	"testing"
)

func TestInitDb_Basic(t *testing.T) {
	msgDb, err := initDb("basic", "")
	assert.Nil(t, err)
	assert.IsType(t, &db.BasicMsgDB{}, msgDb)
}

func TestInitDb_Mongo(t *testing.T) {
	t.Skip("MongoDB test deactivated for now")

	msgDb, err := initDb("mongodb", defaultMongoDbAddr)
	assert.Nil(t, err)
	assert.IsType(t, &db.MongoMsgDB{}, msgDb)
}

func TestInitLogger(t *testing.T) {
	f, err := initLogger("debug")
	assert.Nil(t, err)
	defer f.Close()

	assert.Equal(t, log.DebugLevel, log.GetLevel())
}

func TestRouter(t *testing.T) {
	assert.NotNil(t, router())
}
