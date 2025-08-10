package conector

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoConnector struct {
	client *mongo.Client
	DBName string
}

func (connector *MongoConnector) Connect(ctx context.Context, dbDsn string) error {
	clientOptions := options.Client().ApplyURI(dbDsn)

	var err error = nil
	connector.client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		return err
	}
	err = connector.client.Ping(ctx, nil)
	if err != nil {
		return err
	}

	return nil
}

func (connector *MongoConnector) GetData(ctx context.Context, collectionName string) ([]map[string]any, error) {
	collection := connector.client.Database(connector.DBName).Collection(collectionName)
	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []map[string]any
	for cursor.Next(ctx) {
		var result map[string]any
		err := cursor.Decode(&result)
		if err != nil {
			return nil, err
		}
		delete(result, "_id")
		results = append(results, result)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func (connector *MongoConnector) WriteData(ctx context.Context, collectionName string, data []map[string]any) error {
	collection := connector.client.Database(connector.DBName).Collection(collectionName)

	var docs []any
	for _, doc := range data {
		docs = append(docs, doc)
	}

	_, err := collection.InsertMany(ctx, docs)
	return err
}

func (connector *MongoConnector) Close(ctx context.Context) error {
	return connector.client.Disconnect(ctx)
}
