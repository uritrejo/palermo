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
	defaultConnectTimeout    = 5 * time.Second
)

// MongoMsgDB will store the messages in the database running on the address provided
type MongoMsgDB struct {
	client        *mongo.Client
	msgCollection *mongo.Collection // we could get it from the client, but this saves a lot of redundant code
}

// NewMongoMsgDB returns a new mongo msg db that will connect to the addr provided
// for now it won't take in a password or user name, only unauthenticated access is available
// addr expected is in the format 'mongodb://<host>:<port>' e.g: "mongodb://localhost:27017"
// default values to use for dbName and collectionName are: DefaultMsgDbName and DefaultMsgCollectionName
func NewMongoMsgDB(addr, dbName, collectionName string) (*MongoMsgDB, error) {
	m := &MongoMsgDB{}

	clientOptions := options.Client().ApplyURI(addr)
	ctx, cancel := context.WithTimeout(context.Background(), defaultConnectTimeout)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Errorf("Failed to connect to mongo db with addr: %s; err: %s", addr, err.Error())
		return nil, err
	}
	log.Debug("Successfully connected to mongo db at: ", addr)

	// Check the connection
	ctx, cancel2 := context.WithTimeout(context.Background(), defaultConnectTimeout)
	defer cancel2()
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
	ctx, cancel := context.WithTimeout(context.Background(), defaultConnectTimeout)
	defer cancel()
	err := m.client.Disconnect(ctx)
	if err != nil {
		log.Error("Failed to disconnect client from database: ", err.Error())
	}
}

func (m *MongoMsgDB) GetMsg(id string) (*Msg, error) {
	msg := &Msg{}

	filter := bson.D{primitive.E{Key: "id", Value: id}}

	ctx, cancel := context.WithTimeout(context.Background(), defaultConnectTimeout)
	defer cancel()
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
	ctx, cancel := context.WithTimeout(context.Background(), defaultConnectTimeout)
	defer cancel()
	cursor, err := m.msgCollection.Find(ctx, filter)
	if err != nil {
		log.Error("Failed to find documents: ", err.Error())
		return msgs, err
	}

	ctx, cancel2 := context.WithTimeout(context.Background(), defaultConnectTimeout)
	defer cancel2()
	for cursor.Next(ctx) {
		msg := &Msg{}
		err = cursor.Decode(msg)
		if err != nil {
			return msgs, err
		}
		msgs = append(msgs, msg)
	}

	ctx, cancel3 := context.WithTimeout(context.Background(), defaultConnectTimeout)
	defer cancel3()
	_ = cursor.Close(ctx)

	return msgs, nil
}

func (m *MongoMsgDB) CreateMsg(msg *Msg) error {
	_, err := m.GetMsg(msg.Id)
	if err != nil {
		if IsErrMsgNotFound(err) {
			// we'll only add it if it wasn't found
			ctx, cancel := context.WithTimeout(context.Background(), defaultConnectTimeout)
			defer cancel()
			_, err = m.msgCollection.InsertOne(ctx, msg)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	} else {
		// the message with such Id was already found
		return ErrIdUnavailable{}
	}
	return nil
}

func (m *MongoMsgDB) UpdateMsg(msg *Msg) error {
	filter := bson.D{primitive.E{Key: "id", Value: msg.Id}}

	updater := bson.D{primitive.E{Key: "$set", Value: bson.D{
		primitive.E{Key: "content", Value: msg.Content},
		primitive.E{Key: "isPalindrome", Value: msg.IsPalindrome},
		primitive.E{Key: "modTime", Value: msg.ModTime},
	}}}

	ctx, cancel := context.WithTimeout(context.Background(), defaultConnectTimeout)
	defer cancel()
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

	ctx, cancel := context.WithTimeout(context.Background(), defaultConnectTimeout)
	defer cancel()
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
