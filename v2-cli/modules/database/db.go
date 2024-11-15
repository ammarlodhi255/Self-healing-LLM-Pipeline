package projectdb

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type IDB interface {
	WriteDataToDatabase(data ProgramResultData)
}

type DB struct {
	client *mongo.Client
	DBType DBType // ProductionDB or TestDB
}

type ProgramResultData struct {
	InitialInputQuery string
	SolutionCode      string
	CompilerOutput    string
	TestCases         string
	//Timestamp         string
	//TimeTest          string
}

const UriDb = "mongodb://localhost:27017"

type DBType string

const (
	ProductionDB DBType = "production"
	TestDB       DBType = "test"
)

// The connection string for the DB database

func (db *DB) WriteDataToDatabase(data ProgramResultData) {

	// Choice of database and collection
	collection := db.client.Database(string(db.DBType)).Collection("data")

	// Inserting a document
	_, err := collection.InsertOne(context.TODO(), data)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Document inserted")

	// close the connection, and test if it is closed
	err = db.client.Disconnect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Document closed")

}

func GetDB(clientOptions *options.ClientOptions, dbType DBType) *DB {
	// Connect to the DB server
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// test the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(" DB connection failed:", err)
	}
	log.Println("DB connection successful")

	db := &DB{
		client: client,
		DBType: dbType,
	}

	return db
}
