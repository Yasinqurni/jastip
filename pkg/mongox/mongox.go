package mongox

import (
	"context"
	"errors"
	"sync"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Struct for mongo query builder
type MongoX struct {
	sync.Mutex
	collection *mongo.Collection
}

// Do aggregation function
func (m *MongoX) Aggregate(ctx context.Context, value any, filters ...any) error {
	m.Lock()
	defer m.Unlock()

	if m.collection == nil {
		return errors.New("Collection is undefined")
	}

	cursor, err := m.collection.Aggregate(ctx, filters)
	if err != nil {
		return err
	}

	if err := cursor.All(ctx, value); err != nil {
		return err
	}

	if err := cursor.Err(); err != nil {
		return err
	}
	defer cursor.Close(ctx)

	return nil
}

// Count collection use aggregation
func (m *MongoX) Count(ctx context.Context, filters ...bson.M) (int64, error) {
	m.Lock()
	defer m.Unlock()

	if m.collection == nil {
		return 0, errors.New("Collection is undefined")
	}

	filters = append(filters, bson.M{
		"$count": "count",
	})
	cursor, err := m.collection.Aggregate(ctx, filters)
	if err != nil {
		return 0, err
	}

	var result []map[string]int64
	if err := cursor.All(ctx, &result); err != nil {
		return 0, err
	}

	if err := cursor.Err(); err != nil {
		return 0, err
	}
	defer cursor.Close(ctx)

	if len(result) == 0 {
		return 0, nil
	}

	return result[0]["count"], nil
}

func (m *MongoX) Collection(collection *mongo.Collection) *MongoX {
	m.collection = collection
	return m
}
