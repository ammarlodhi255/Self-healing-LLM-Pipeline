package projectdb

import (
	"context"
	"testing"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestDB_WriteDataToDatabase(t *testing.T) {
	clientOptions := options.Client().ApplyURI(UriDb)
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		t.Fatal(err)
	}

	db := GetDB(clientOptions, TestDB)

	data := ProgramResultData{
		InitialInputQuery: "test",
		SolutionCode:      "test",
		CompilerOutput:    "test",
		TestCases:         "test",
	}

	db.WriteDataToDatabase(data)

	// Check if the data was inserted
	collection := client.Database(string(db.DBType)).Collection("data")
	result := collection.FindOne(context.TODO(), data)
	if result.Err() != nil {
		t.Fatal(result.Err())
	}

	// Drop the collection
	err = collection.Drop(context.TODO())
	if err != nil {
		t.Fatal(err)
	}
}

