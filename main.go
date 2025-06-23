package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/Quanscendence/updateImageURL/internal/config"
	utils "github.com/Quanscendence/updateImageURL/internal/util"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	var getenv func(string) string
	getenv = func(key string) string {
		return os.Getenv(key)
	}

	err = run(context.Background(), getenv, os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

}

func run(ctx context.Context, getenv func(string) string, w io.Writer) error {

	cfg, err := config.LoadEnv(getenv)
	if err != nil {
		return err
	}

	mongoClient, err := ConnectDB(cfg)
	if err != nil {
		return err
	}

	mongoDatabase := mongoClient.Database(cfg.MongoDataBase)

	err = utils.UpdateFieldsWithProxy(ctx, cfg.Domain, mongoDatabase, cfg.CollectionName, cfg.CollectionKeyName)
	if err != nil {
		return err
	}
	log.Println("update all the field with proxy image url")

	return nil

}

func ConnectDB(cfg *config.Config) (*mongo.Client, error) {

	db, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		log.Println("MongoDB restore failed:", err.Error())
		return nil, fmt.Errorf("failed to connect to db: %w", err)
	}

	err = db.Ping(context.TODO(), nil)
	if err != nil {
		log.Println("failed to ping db:", err.Error())
		return nil, fmt.Errorf("failed to ping db: %w", err)
	}

	return db, err
}
