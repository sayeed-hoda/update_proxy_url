package utils

import (
	"context"
	"errors"
	"fmt"
	urlNet "net/url"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func ProxyUrlBuilder(domain, url string) (string, error) {
	_, err := urlNet.Parse(url)
	if err != nil {
		return "", errors.New("url should be valid")
	}

	parseURL, err := urlNet.Parse(fmt.Sprintf("https://%s/api/v1/file-server/files/images/load-url", domain))
	if err != nil {
		return "", errors.New("domain should be valid")
	}

	encodedURL := urlNet.QueryEscape(url)
	finalClientRequestURL := fmt.Sprintf("%s?url=%s", parseURL.String(), encodedURL)

	return finalClientRequestURL, nil
}

func UpdateFieldsWithProxy(ctx context.Context, domain string, db *mongo.Database, collectionName string, keys []string) error {
	collection := db.Collection(collectionName)

	// Get all documents
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return fmt.Errorf("find error: %w", err)
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var doc bson.M
		if err := cursor.Decode(&doc); err != nil {
			return fmt.Errorf("decode error: %w", err)
		}

		id := doc["_id"]
		updateFields := bson.M{}

		for _, key := range keys {
			val, ok := doc[key]
			if !ok {
				continue
			}

			strVal, ok := val.(string)
			if !ok {
				continue
			}

			newVal, err := ProxyUrlBuilder(domain, strVal)
			if err != nil {
				return err
			}
			if newVal != strVal {
				updateFields[key] = newVal
			}
		}

		// Only update if there are fields to change
		if len(updateFields) > 0 {
			filter := bson.M{"_id": id}
			update := bson.M{"$set": updateFields}

			_, err := collection.UpdateOne(ctx, filter, update)
			if err != nil {
				return fmt.Errorf("update error for _id=%v: %w", id, err)
			}
		}
	}

	if err := cursor.Err(); err != nil {
		return fmt.Errorf("cursor error: %w", err)
	}

	return nil
}
