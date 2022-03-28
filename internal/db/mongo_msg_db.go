package db

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"strconv"
	"sync"
	"time"
)

type ToDoItem struct {
	//Id          primitive.ObjectID `bson:"_id"`  // seems like it's automatically added
	Code        string    `bson:"code"`
	CreatedAt   time.Time `bson:"created_at"`
	Title       string    `bson:"title"`
	Description string    `bson:"description"`
	Completed   bool      `bson:"completed"`
}

func (t *ToDoItem) ToString() string {
	return "Task title: " + t.Title + ", Description: " + t.Description +
		", Completed: " + strconv.FormatBool(t.Completed) + ", Created at: " + t.CreatedAt.String()
}

// AddToDoItem will add an item into the database
func AddToDoItem(item *ToDoItem) error {
	client, err := GetMongoClient()
	if err != nil {
		log.Println("Error getting the mongo client: ", err.Error())
		return err
	}

	// ??++ me pregunto si esta database y collection ya tienen que haber existido?
	// y si no, no es ineficiente pullearlas siempre?
	//Create a handle to the respective collection in the database.
	collection := client.Database(dbName).Collection(collectionName)

	_, err = collection.InsertOne(context.TODO(), item) // there's also InsertMany
	if err != nil {
		return err
	}

	return nil
}

// ??++ For sure there's gotta be a way to get the item by its id
// GetToDoItemByCode will get the to do item associated with the given code
func GetToDoItemByCode(code string) (*ToDoItem, error) {
	item := &ToDoItem{}

	//Define filter query for fetching specific document from collection
	filter := bson.D{primitive.E{Key: "code", Value: code}}

	//Get MongoDB connection using connectionhelper.
	client, err := GetMongoClient()
	if err != nil {
		log.Println("Failed to get mongo client: ", err.Error())
		return nil, err
	}

	collection := client.Database(dbName).Collection(collectionName)

	err = collection.FindOne(context.TODO(), filter).Decode(item)
	if err != nil {
		log.Println("Failed to find document: ", err.Error())
		return nil, err
	}

	return item, nil
}

// GetAllToDoItems will get all the to-do items from the database
func GetAllToDoItems() ([]*ToDoItem, error) {
	var items []*ToDoItem

	//bson.D{{}} specifies 'all documents'
	filter := bson.D{{}}

	//Get MongoDB connection using connectionhelper.
	client, err := GetMongoClient()
	if err != nil {
		log.Println("Failed to get mongo client: ", err.Error())
		return items, err
	}

	collection := client.Database(dbName).Collection(collectionName)

	// cursor is like an iterator
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		log.Println("Failed to find documents: ", err.Error())
		return items, err
	}

	for cursor.Next(context.TODO()) {
		t := &ToDoItem{}
		err := cursor.Decode(t)
		if err != nil {
			return items, err
		}
		items = append(items, t)
	}
	// once exhausted, close the cursor
	_ = cursor.Close(context.TODO())
	if len(items) == 0 {
		return items, mongo.ErrNoDocuments
	}

	return items, nil
}

// MarkCompleted will mark an item as completed
func MarkCompleted(code string) (*ToDoItem, error) {
	item := &ToDoItem{}

	//Define filter query for fetching specific document from collection
	filter := bson.D{primitive.E{Key: "code", Value: code}}

	//Define updater for to specifiy change to be updated.
	updater := bson.D{primitive.E{Key: "$set", Value: bson.D{
		primitive.E{Key: "completed", Value: true},
	}}}

	//Get MongoDB connection using connectionhelper.
	client, err := GetMongoClient()
	if err != nil {
		log.Println("Failed to get mongo client: ", err.Error())
		return nil, err
	}

	collection := client.Database(dbName).Collection(collectionName)

	_, err = collection.UpdateOne(context.TODO(), filter, updater)
	if err != nil {
		log.Println("Failed to update document: ", err.Error())
		return nil, err
	}

	return item, nil
}

// DeleteToDoItem will delete an item from the collection
func DeleteToDoItem(code string) error {
	//Define filter query for fetching specific document from collection
	filter := bson.D{primitive.E{Key: "code", Value: code}}

	//Get MongoDB connection using connectionhelper.
	client, err := GetMongoClient()
	if err != nil {
		log.Println("Failed to get mongo client: ", err.Error())
		return err
	}

	collection := client.Database(dbName).Collection(collectionName)

	_, err = collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		log.Println("Failed to find document: ", err.Error())
		return err
	}

	return nil
}

/* Used to create a singleton object of MongoDB client.
Initialized and exposed through  GetMongoClient().*/
var clientInstance *mongo.Client

//Used during creation of singleton client object in GetMongoClient().
var clientInstanceError error

//Used to execute client creation procedure only once.
var mongoOnce sync.Once

//I have used below constants just to hold required database config's.
const (
	mongoAddr      = "mongodb://localhost:27017"
	dbName         = "db_todo_list"
	collectionName = "col_todo"
)

//GetMongoClient - Return mongodb connection to work with
func GetMongoClient() (*mongo.Client, error) {
	//Perform connection creation operation only once.
	mongoOnce.Do(func() {
		// Set client options
		clientOptions := options.Client().ApplyURI(mongoAddr)
		// Connect to MongoDB
		client, err := mongo.Connect(context.TODO(), clientOptions)
		if err != nil {
			clientInstanceError = err
		}
		// Check the connection
		// ??++ weird error handling
		err = client.Ping(context.TODO(), nil)
		if err != nil {
			clientInstanceError = err
		}
		clientInstance = client
	})
	return clientInstance, clientInstanceError
}

// todo: delete this once the rest work:
// test ones for the DB

func addItemQuick() {
	err := AddToDoItem(&ToDoItem{
		Code:        "002",
		CreatedAt:   time.Now(),
		Title:       "Espero que ya este",
		Description: "a salir",
		Completed:   false,
	})
	if err != nil {
		log.Fatal("Failed to add todo item: ", err.Error())
	}
}

func getAllItemsQuick() {
	items, err := GetAllToDoItems()
	if err != nil {
		log.Fatal("Failed to add todo item: ", err.Error())
	}

	log.Println("All items are: ")
	for _, item := range items {
		log.Println(item.ToString())
	}
}
