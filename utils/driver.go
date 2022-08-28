package utils

import (
	"api-desatanggap/config"
	"context"
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"gorm.io/gorm"
)

type DatabaseDriver string

const (
	MongoDB DatabaseDriver = "mongodb"
	Mysql   DatabaseDriver = "mysql"
)

type DatabaseConnection struct {
	Driver DatabaseDriver

	MongoDB     *mongo.Database
	mongoClient *mongo.Client
	Mysql       *gorm.DB
}

func NewConnectionDatabase(config *config.AppConfig) *DatabaseConnection {
	var db DatabaseConnection
	dbName := os.Getenv("MONGO_DBNAME")
	db.mongoClient = newMongodb(config)
	db.MongoDB = db.mongoClient.Database(dbName)

	return &db
}

func newMongodb(config *config.AppConfig) *mongo.Client {
	// dbUser := os.Getenv("MONGO_USERNAME")
	// dbPass := os.Getenv("MONGO_PASSWORD")
	dbUrl := os.Getenv("MONGO_HOST")
	dbport := os.Getenv("MONGO_PORT")
	dbPass := os.Getenv("MONGO_PASS")
	dbUser := os.Getenv("MONGO_USER")
	dbEtc := os.Getenv("MONGO_ETC")

	url := "mongodb://" + dbUser + ":" + dbPass + "@" + dbUrl + ":" + dbport + dbEtc
	fmt.Println(url)

	client, err := mongo.NewClient(options.Client().ApplyURI(url))
	if err != nil {
		panic(err)
	}

	err = client.Connect(context.Background())
	if err != nil {
		panic(err)
	}

	err = client.Ping(context.Background(), readpref.Primary())
	if err != nil {
		panic(err)
	}

	return client
}

func (db *DatabaseConnection) CloseConnection() {
	db.mongoClient.Disconnect(context.Background())
}
