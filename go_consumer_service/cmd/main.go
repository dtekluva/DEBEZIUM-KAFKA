package main

import (
	"context"
	"go_consumer_service/cmd/api"
	"go_consumer_service/config"
	"go_consumer_service/consumers"
	"go_consumer_service/utils"
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
	newUtils := utils.NewUtils(dbClient.Database("lotto2"))
	kafkaBrokerUrl := webutils.GetEnv("KAFKA_BROKER_URL", "localhost:9092")
	kafkaConsumer := consumers.NewKafkaConsumer(&kafkaBrokerUrl, dbClient.Database("lotto2"), newUtils)
	// go kafkaConsumer.ConsumeDebeziumMobidTrackerTask()
	go kafkaConsumer.ConsumeDebeziumSecureDataDumpTask()
	// go kafkaConsumer.ConsumeLottoDebeziumEvent()
	// go kafkaConsumer.ConsumeLotteryModelDebeziumEvent()
	// go kafkaConsumer.ConsumeAwoofGameTableDebeziumEvent()
	// go kafkaConsumer.ConsumeSecureDTransactionDebeziumEvent()
	// go kafkaConsumer.ConsumeConstantTableEvent()

	apiServer := api.NewAPIServer(":6000", dbClient)

	if err := apiServer.Start(); err != nil {
		panic(err)
	}

}
