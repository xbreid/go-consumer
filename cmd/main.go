package main

import (
	"context"
	"database/sql"
	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"go-consumer/ent"
	"log"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Config struct {
	DB                *ent.Client
	SQS               *sqs.Client
	queueUrl          string
	visibilityTimeout int32
	waitingTimeout    int32
}

var retryCount int64

func main() {
	log.Printf("Service started...")

	queueUrl := os.Getenv("SQS_URL")

	// Connect to DB
	client := connectToDB()
	if client == nil {
		log.Panic("Cannot connect to Postgres!")
	}

	// Run the auto migration tool.
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

	// AWS Config Init
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("configuration error, " + err.Error())
	}

	// SQS Client Init
	sqsSvc := sqs.NewFromConfig(cfg)

	app := Config{
		DB:                client,
		SQS:               sqsSvc,
		queueUrl:          queueUrl,
		visibilityTimeout: 60 * 10,
		waitingTimeout:    20,
	}

	app.runConsumer()

	log.Println("service is safely stopped")
}

func OpenDB(dsn string) (*ent.Client, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	// Create an ent.Driver from `db`.
	drv := entsql.OpenDB(dialect.Postgres, db)
	return ent.NewClient(ent.Driver(drv)), nil
}

func connectToDB() *ent.Client {
	dsn := os.Getenv("DSN")

	for {
		connection, err := OpenDB(dsn)
		if err != nil {
			log.Println("Postgres not yet ready...")
			retryCount++
		} else {
			log.Println("Connected to Postgres!")
			return connection
		}

		if retryCount > 10 {
			log.Println(err)
			return nil
		}

		log.Println("Backing off for two seconds...")
		time.Sleep(2 * time.Second)
		continue
	}
}
