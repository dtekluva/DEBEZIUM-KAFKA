package consumers

import (
	"context"
	"encoding/json"
	"go_consumer_service/types"
	"log"

	"github.com/segmentio/kafka-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type KafkaConsumer struct {
	brokerUrl *string
	database  *mongo.Database
}

// NewKafkaConsumer creates a new instance of KafkaConsumer
func NewKafkaConsumer(brokerUrl *string, database *mongo.Database) *KafkaConsumer {
	return &KafkaConsumer{
		brokerUrl: brokerUrl,
		database:  database,
	}
}

// Consume DebeziumMobidTracker Task
func (kc *KafkaConsumer) ConsumeDebeziumMobidTrackerTask() {
	// Create a new consumer group
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{*kc.brokerUrl},
		GroupID: "lotto-mobidtracker-cdc",
		Topic:   "postgres.public.ads_tracker_mobidtracker",
	})
	collection := kc.database.Collection("mobid_tracker")
	log.Println("Kafka Consumer Started for MobidTracker CDC......")

	for {
		msg, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Fatal(err)
		}
		var event types.MobidTrackerEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.Printf("failed to unmarshall message: %v\n", err)
			continue
		}
		log.Println("Received message: ", event.Payload.After)
		ops := event.Payload.Op
		log.Println("Operation: ", ops)
		switch ops {
		case "c", "r":
			log.Println("Inserting MobidTracker: ", event.Payload.After.ID)
			if event.Payload.After != nil {
				_, err := collection.InsertOne(context.Background(), event.Payload.After)
				if err != nil {
					log.Printf("failed to insert mobidtracker: %v\n", err)
				} else {
					log.Println("MobidTracker inserted successfully")
				}
			}
		case "u":
			log.Println("Updating MobidTracker: ", event.Payload.After.ID)
			if event.Payload.After != nil {
				filter := bson.M{"id": event.Payload.After.ID}
				update := bson.M{"$set": event.Payload.After}
				_, err := collection.UpdateOne(context.Background(), filter, update)
				if err != nil {
					log.Printf("failed to update mobidtracker: %v\n", err)
				} else {
					log.Println("MobidTracker updated successfully")
				}
			}
		case "d":
			log.Println("Deleting MobidTracker: ", event.Payload.Before.ID)
			if event.Payload.Before != nil {
				_, err := collection.DeleteOne(context.Background(), event.Payload.Before.ID)
				if err != nil {
					log.Printf("failed to delete mobidtracker: %v\n", err)
				} else {
					log.Println("MobidTracker deleted successfully")
				}
			}
		default:
			log.Println("Unknown operation: ", ops)
		}
	}
}
