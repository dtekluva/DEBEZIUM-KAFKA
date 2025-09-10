package main

import (
	"context"
	"go_consumer_service/cmd/api"
	"go_consumer_service/config"
	"go_consumer_service/consumers"
	"time"

	"github.com/Ayobami6/webutils"
	_ "github.com/joho/godotenv/autoload"
)

func main() {

	dbUrl := webutils.GetEnv("MONGO_URL", "mongodb://localhost:27017")
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	dbClient, err := config.ConnectDb(ctx, dbUrl)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = dbClient.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
	kafkaBrokerUrl := webutils.GetEnv("KAFKA_BROKER_URL", "localhost:9092")
	kafkaConsumer := consumers.NewKafkaConsumer(&kafkaBrokerUrl, dbClient.Database("lotto2"))
	go kafkaConsumer.ConsumeDebeziumMobidTrackerTask()

	apiServer := api.NewAPIServer(":6000", dbClient)

	if err := apiServer.Start(); err != nil {
		panic(err)
	}

}
