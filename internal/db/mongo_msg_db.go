package db

import (
	"context"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

const (
	DefaultMsgDbName         = "msgDB"
	DefaultMsgCollectionName = "msgCollection"
	defaultConnectTimeout    = 7 * time.Second
)

// MongoMsgDB will store the messages in the database running on the address provided
type MongoMsgDB struct {
	client        *mongo.Client
	msgCollection *mongo.Collection // we could get it from the client, but this saves a lot of redundant code
}

// todo: look into the auth stuff, add a username and password

// NewMongoMsgDB returns a new mongo msg db that will connect to the addr provided
// addr expected is in the format 'mongodb://<host>:<port>' e.g: "mongodb://localhost:27017"
// default values to use for dbName and collectionName are: DefaultMsgDbName and DefaultMsgCollectionName
func NewMongoMsgDB(addr, dbName, collectionName string) (*MongoMsgDB, error) {
	m := &MongoMsgDB{}

	clientOptions := options.Client().ApplyURI(addr)
	ctx, _ := context.WithTimeout(context.Background(), defaultConnectTimeout)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Errorf("Failed to connect to mongo db with addr: %s; err: %s", addr, err.Error())
		return nil, err
	}
	log.Debug("Successfully connected to mongo db at: ", addr)

	// Check the connection
	ctx, _ = context.WithTimeout(context.Background(), defaultConnectTimeout)
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Error("Failed to ping mongo db: ", err.Error())
		return nil, err
	}

	m.client = client
	m.msgCollection = client.Database(dbName).Collection(collectionName)

	return m, nil
}

func (m *MongoMsgDB) Close() {
	ctx, _ := context.WithTimeout(context.Background(), defaultConnectTimeout)
	err := m.client.Disconnect(ctx)
	if err != nil {
		log.Error("Failed to disconnect client from database: ", err.Error())
	}
}

func (m *MongoMsgDB) GetMsg(id string) (*Msg, error) {
	msg := &Msg{}

	filter := bson.D{primitive.E{Key: "id", Value: id}}

	ctx, _ := context.WithTimeout(context.Background(), defaultConnectTimeout)
	err := m.msgCollection.FindOne(ctx, filter).Decode(msg)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrMsgNotFound{}
		}
		log.Error("Failed to find msg: ", err.Error())
		return nil, err
	}

	return msg, nil
}

// todo: verify that if there's nothing it just returns an empty slice
func (m *MongoMsgDB) GetAllMsgs() ([]*Msg, error) {
	var msgs []*Msg

	//bson.D{{}} specifies 'all documents'
	filter := bson.D{{}}

	// cursor is like an iterator
	ctx, _ := context.WithTimeout(context.Background(), defaultConnectTimeout)
	cursor, err := m.msgCollection.Find(ctx, filter)
	if err != nil {
		log.Error("Failed to find documents: ", err.Error())
		return msgs, err
	}

	ctx, _ = context.WithTimeout(context.Background(), defaultConnectTimeout)
	for cursor.Next(ctx) {
		msg := &Msg{}
		err = cursor.Decode(msg)
		if err != nil {
			return msgs, err
		}
		msgs = append(msgs, msg)
	}

	ctx, _ = context.WithTimeout(context.Background(), defaultConnectTimeout)
	_ = cursor.Close(ctx)

	return msgs, nil
}

func (m *MongoMsgDB) CreateMsg(msg *Msg) error {
	_, err := m.GetMsg(msg.Id)
	if err != nil {
		if IsErrMsgNotFound(err) {
			// we'll only add it if it wasn't found
			ctx, _ := context.WithTimeout(context.Background(), defaultConnectTimeout)
			_, err = m.msgCollection.InsertOne(ctx, msg)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	return ErrIdUnavailable{}
}

func (m *MongoMsgDB) UpdateMsg(msg *Msg) error {
	filter := bson.D{primitive.E{Key: "id", Value: msg.Id}}

	updater := bson.D{primitive.E{Key: "$set", Value: bson.D{
		primitive.E{Key: "content", Value: msg.Content},
		primitive.E{Key: "isPalindrome", Value: msg.IsPalindrome},
		primitive.E{Key: "modTime", Value: msg.ModTime},
	}}}

	ctx, _ := context.WithTimeout(context.Background(), defaultConnectTimeout)
	result, err := m.msgCollection.UpdateOne(ctx, filter, updater)
	if err != nil {
		log.Error("Failed to update document: ", err.Error())
		return err
	}
	if result.MatchedCount == 0 {
		return ErrMsgNotFound{}
	}

	return nil
}

func (m *MongoMsgDB) DeleteMsg(id string) error {
	filter := bson.D{primitive.E{Key: "id", Value: id}}

	ctx, _ := context.WithTimeout(context.Background(), defaultConnectTimeout)
	result, err := m.msgCollection.DeleteOne(ctx, filter)
	if err != nil {
		log.Error("Failed to delete document: ", err.Error())
		return err
	}
	if result.DeletedCount == 0 {
		return ErrMsgNotFound{}
	}

	return nil
}
