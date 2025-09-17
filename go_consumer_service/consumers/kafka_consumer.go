package consumers

import (
	"context"
	"encoding/json"
	"go_consumer_service/types"
	"go_consumer_service/utils"
	"log"

	"github.com/segmentio/kafka-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type KafkaConsumer struct {
	brokerUrl *string
	database  *mongo.Database
	utils     *utils.Utils
}

// NewKafkaConsumer creates a new instance of KafkaConsumer
func NewKafkaConsumer(brokerUrl *string, database *mongo.Database, utils *utils.Utils) *KafkaConsumer {
	return &KafkaConsumer{
		brokerUrl: brokerUrl,
		database:  database,
		utils:     utils,
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
				opts := options.UpdateOptions{}
				opts.SetUpsert(true)
				_, err := collection.UpdateOne(context.Background(), filter, update, &opts)
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

// ConsumeDebeziumSecureDataDump consumes debezium messages to topic securedatadump
func (kc *KafkaConsumer) ConsumeDebeziumSecureDataDumpTask() {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{*kc.brokerUrl},
		GroupID: "lotto-securedatadump-cdc-test",
		Topic:   "postgres.public.wyse_ussd_secureddatadump",
	})
	collection := kc.database.Collection("secure_data_dump")
	log.Println("Kafka Consumer Started for SecureDataDump CDC......")

	for {
		msg, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Fatal(err)
		}
		var event types.SecureDataDumpEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.Printf("failed to unmarshall message: %v\n", err)
			continue
		}
		log.Println("Received message: ", event.Payload.After)
		ops := event.Payload.Op
		log.Println("Operation: ", ops)
		switch ops {
		case "c", "r":
			log.Println("Inserting SecureDataDump: ", event.Payload.After.ID)
			log.Println("SecureDataDump: ", event.Payload.After.Data)
			if event.Payload.After != nil {
				_, err := collection.InsertOne(context.Background(), event.Payload.After)
				if err != nil {
					log.Printf("failed to insert securedatadump: %v\n", err)
				} else {
					log.Println("SecureDataDump inserted successfully")
					// send marketing message
					kc.utils.SendMarketingPartnersPostback(event.Payload.After.ID)
				}
			}
		case "u":
			log.Println("Updating SecureDataDump: ", event.Payload.After.ID)
			if event.Payload.After != nil {
				filter := bson.M{"id": event.Payload.After.ID}
				update := bson.M{"$set": event.Payload.After}
				opts := options.UpdateOptions{}
				opts.SetUpsert(true)
				_, err := collection.UpdateOne(context.Background(), filter, update, &opts)
				if err != nil {
					log.Printf("failed to update securedatadump: %v\n", err)
				} else {
					log.Println("SecureDataDump updated successfully")
				}
			}
		case "d":
			log.Println("Deleting SecureDataDump: ", event.Payload.Before.ID)
			if event.Payload.Before != nil {
				_, err := collection.DeleteOne(context.Background(), event.Payload.Before.ID)
				if err != nil {
					log.Printf("failed to delete securedatadump: %v\n", err)
				} else {
					log.Println("SecureDataDump deleted successfully")
				}
			}
		default:
			log.Println("Unknown operation: ", ops)
		}

	}
}

func (kc *KafkaConsumer) ConsumeLottoDebeziumEvent() {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{*kc.brokerUrl},
		GroupID: "lotto-event-cdc-test",
		Topic:   "postgres.public.main_lottoticket",
	})
	collection := kc.database.Collection("lotto_event")
	log.Println("Kafka Consumer Started for LottoEvent CDC......")

	for {
		msg, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Fatal(err)
		}
		var event types.LottoTicketEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.Printf("failed to unmarshall message: %v\n", err)
			continue
		}
		log.Println("Received message: ", event.Payload.After)
		ops := event.Payload.Op
		log.Println("Operation: ", ops)
		switch ops {
		case "c", "r":
			log.Println("Inserting LottoEvent: ", event.Payload.After.ID)
			if event.Payload.After != nil {
				_, err := collection.InsertOne(context.Background(), event.Payload.After)
				if err != nil {
					log.Printf("failed to insert lottoevent: %v\n", err)
				} else {
					log.Println("LottoEvent inserted successfully")
				}
			}
		case "u":
			log.Println("Updating LottoEvent: ", event.Payload.After.ID)
			if event.Payload.After != nil {
				filter := bson.M{"id": event.Payload.After.ID}
				update := bson.M{"$set": event.Payload.After}
				opts := options.UpdateOptions{}
				opts.SetUpsert(true)
				_, err := collection.UpdateOne(context.Background(), filter, update, &opts)
				if err != nil {
					log.Printf("failed to update lottoevent: %v\n", err)
				} else {
					log.Println("LottoEvent updated successfully")
				}
			}
		case "d":
			log.Println("Deleting LottoEvent: ", event.Payload.Before.ID)
			if event.Payload.Before != nil {
				_, err := collection.DeleteOne(context.Background(), event.Payload.Before.ID)
				if err != nil {
					log.Printf("failed to delete lottoevent: %v\n", err)
				} else {
					log.Println("LottoEvent deleted successfully")
				}
			}
		default:
			log.Println("Unknown operation: ", ops)
		}
	}
}

func (kc *KafkaConsumer) ConsumeLotteryModelDebeziumEvent() {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{*kc.brokerUrl},
		GroupID: "lottery-model-cdc-test",
		Topic:   "postgres.public.main_lotterymodel",
	})

	collection := kc.database.Collection("lottery_model")
	log.Println("Kafka Consumer Started for LotteryModel CDC......")

	for {
		msg, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Fatal(err)
		}
		var event types.LotteryEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.Printf("failed to unmarshall message: %v\n", err)
			continue
		}
		log.Println("Received message: ", event.Payload.After)
		ops := event.Payload.Op
		log.Println("Operation: ", ops)
		switch ops {
		case "c", "r":
			log.Println("Inserting LotteryModel: ", event.Payload.After.ID)
			if event.Payload.After != nil {
				_, err := collection.InsertOne(context.Background(), event.Payload.After)
				if err != nil {
					log.Printf("failed to insert lotterymodel: %v\n", err)
				} else {
					log.Println("LotteryModel inserted successfully")
				}
			}
		case "u":
			log.Println("Updating LotteryModel: ", event.Payload.After.ID)
			if event.Payload.After != nil {
				filter := bson.M{"id": event.Payload.After.ID}
				update := bson.M{"$set": event.Payload.After}
				opts := options.UpdateOptions{}
				opts.SetUpsert(true)
				_, err := collection.UpdateOne(context.Background(), filter, update, &opts)
				if err != nil {
					log.Printf("failed to update lotterymodel: %v\n", err)
				} else {
					log.Println("LotteryModel updated successfully")
				}
			}
		case "d":
			log.Println("Deleting LotteryModel: ", event.Payload.Before.ID)
			if event.Payload.Before != nil {
				_, err := collection.DeleteOne(context.Background(), event.Payload.Before.ID)
				if err != nil {
					log.Printf("failed to delete lotterymodel: %v\n", err)
				} else {
					log.Println("LotteryModel deleted successfully")
				}
			}
		default:
			log.Println("Unknown operation: ", ops)
		}
	}

}

func (kc *KafkaConsumer) ConsumeAwoofGameTableDebeziumEvent() {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{*kc.brokerUrl},
		GroupID: "awoof-game-table-cdc-test",
		Topic:   "postgres.public.awoof_app_awoofgametable",
	})
	collection := kc.database.Collection("awoof_game_table")
	log.Println("Kafka Consumer Started for AwoofGameTable CDC......")
	for {
		msg, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Fatal(err)
		}
		var event types.AwoofGameTableEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.Printf("failed to unmarshall message: %v\n", err)
			continue
		}
		log.Println("Received message: ", event.Payload.After)
		ops := event.Payload.Op
		log.Println("Operation: ", ops)
		switch ops {
		case "c", "r":
			log.Println("Inserting AwoofGameTable: ", event.Payload.After.ID)
			if event.Payload.After != nil {
				_, err := collection.InsertOne(context.Background(), event.Payload.After)
				if err != nil {
					log.Printf("failed to insert awoofgametable: %v\n", err)
				} else {
					log.Println("AwoofGameTable inserted successfully")
				}
			}
		case "u":
			log.Println("Updating AwoofGameTable: ", event.Payload.After.ID)
			if event.Payload.After != nil {
				filter := bson.M{"id": event.Payload.After.ID}
				update := bson.M{"$set": event.Payload.After}
				opts := options.UpdateOptions{}
				opts.SetUpsert(true)
				_, err := collection.UpdateOne(context.Background(), filter, update, &opts)
				if err != nil {
					log.Printf("failed to update awoofgametable: %v\n", err)
				} else {
					log.Println("AwoofGameTable updated successfully")
				}
			}
		case "d":
			log.Println("Deleting AwoofGameTable: ", event.Payload.Before.ID)
			if event.Payload.Before != nil {
				_, err := collection.DeleteOne(context.Background(), event.Payload.Before.ID)
				if err != nil {
					log.Printf("failed to delete awoofgametable: %v\n", err)
				} else {
					log.Println("AwoofGameTable deleted successfully")
				}
			}
		default:
			log.Println("Unknown operation: ", ops)
		}
	}
}

func (kc *KafkaConsumer) ConsumeSecureDTransactionDebeziumEvent() {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{*kc.brokerUrl},
		GroupID: "secure-d-transaction-cdc-test",
		Topic:   "postgres.public.wyse_ussd_securedtransaction",
	})
	collection := kc.database.Collection("secure_d_transaction")
	log.Println("Kafka Consumer Started for SecureDTransaction CDC......")
	for {
		msg, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Fatal(err)
		}
		var event types.SecureDTransactionEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.Printf("failed to unmarshall message: %v\n", err)
			continue
		}
		log.Println("Received message: ", event.Payload.After)
		ops := event.Payload.Op
		log.Println("Operation: ", ops)
		switch ops {
		case "c", "r":
			log.Println("Inserting SecureDTransaction: ", event.Payload.After.ID)
			if event.Payload.After != nil {
				_, err := collection.InsertOne(context.Background(), event.Payload.After)
				if err != nil {
					log.Printf("failed to insert securedtransaction: %v\n", err)
				} else {
					log.Println("SecureDTransaction inserted successfully")
				}
			}
		case "u":
			log.Println("Updating SecureDTransaction: ", event.Payload.After.ID)
			if event.Payload.After != nil {
				filter := bson.M{"id": event.Payload.After.ID}
				update := bson.M{"$set": event.Payload.After}
				opts := options.UpdateOptions{}
				opts.SetUpsert(true)
				_, err := collection.UpdateOne(context.Background(), filter, update, &opts)
				if err != nil {
					log.Printf("failed to update securedtransaction: %v\n", err)
				} else {
					log.Println("SecureDTransaction updated successfully")
				}
			}
		case "d":
			log.Println("Deleting SecureDTransaction: ", event.Payload.Before.ID)
			if event.Payload.Before != nil {
				_, err := collection.DeleteOne(context.Background(), event.Payload.Before.ID)
				if err != nil {
					log.Printf("failed to delete securedtransaction: %v\n", err)
				} else {
					log.Println("SecureDTransaction deleted successfully")
				}
			}
		default:
			log.Println("Unknown operation: ", ops)
		}
	}
}

func (kc *KafkaConsumer) ConsumeConstantTableEvent() {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{*kc.brokerUrl},
		GroupID: "constant-table-cdc-test",
		Topic:   "postgres.public.admin_dashboard_constanttable",
	})
	collection := kc.database.Collection("constant_table")
	log.Println("Kafka Consumer Started for ConstantTable CDC......")
	for {
		msg, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Fatal(err)
		}
		var event types.ConstantTableDebeziumEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.Printf("failed to unmarshall message: %v\n", err)
			continue
		}
		log.Println("Received message: ", event.Payload.After)
		ops := event.Payload.Op
		log.Println("Operation: ", ops)
		switch ops {
		case "c", "r":
			log.Println("Inserting ConstantTable: ", event.Payload.After.ID)
			if event.Payload.After != nil {
				_, err := collection.InsertOne(context.Background(), event.Payload.After)
				if err != nil {
					log.Printf("failed to insert constanttable: %v\n", err)
				} else {
					log.Println("ConstantTable inserted successfully")
				}
			}
		case "u":
			log.Println("Updating ConstantTable: ", event.Payload.After.ID)
			if event.Payload.After != nil {
				filter := bson.M{"id": event.Payload.After.ID}
				update := bson.M{"$set": event.Payload.After}
				opts := options.UpdateOptions{}
				opts.SetUpsert(true)
				_, err := collection.UpdateOne(context.Background(), filter, update, &opts)
				if err != nil {
					log.Printf("failed to update constanttable: %v\n", err)
				} else {
					log.Println("ConstantTable updated successfully")
				}
			}
		case "d":
			log.Println("Deleting ConstantTable: ", event.Payload.Before.ID)
			if event.Payload.Before != nil {
				_, err := collection.DeleteOne(context.Background(), event.Payload.Before.ID)
				if err != nil {
					log.Printf("failed to delete constanttable: %v\n", err)
				} else {
					log.Println("ConstantTable deleted successfully")
				}
			}
		default:
			log.Println("Unknown operation: ", ops)
		}
	}
}

func (kc *KafkaConsumer) ConsumeSoccerPredictionDebeziumEvent() {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{*kc.brokerUrl},
		GroupID: "soccer-prediction-cdc-test",
		Topic:   "postgres.public.main_soccerprediction",
	})
	collection := kc.database.Collection("soccer_prediction")
	log.Println("Kafka Consumer Started for SoccerPrediction CDC......")
	for {
		msg, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Fatal(err)
		}
		var event types.SoccerPredictionEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.Printf("failed to unmarshall message: %v\n", err)
			continue
		}
		log.Println("Received message: ", event.Payload.After)
		ops := event.Payload.Op
		log.Println("Operation: ", ops)
		switch ops {
		case "c", "r":
			log.Println("Inserting SoccerPrediction: ", event.Payload.After.ID)
			if event.Payload.After != nil {
				_, err := collection.InsertOne(context.Background(), event.Payload.After)
				if err != nil {
					log.Printf("failed to insert soccerprediction: %v\n", err)
				} else {
					log.Println("SoccerPrediction inserted successfully")
				}
			}
		case "u":
			log.Println("Updating SoccerPrediction: ", event.Payload.After.ID)
			if event.Payload.After != nil {
				filter := bson.M{"id": event.Payload.After.ID}
				update := bson.M{"$set": event.Payload.After}
				opts := options.UpdateOptions{}
				opts.SetUpsert(true)
				_, err := collection.UpdateOne(context.Background(), filter, update, &opts)
				if err != nil {
					log.Printf("failed to update soccerprediction: %v\n", err)
				} else {
					log.Println("SoccerPrediction updated successfully")
				}
			}
		case "d":
			log.Println("Deleting SoccerPrediction: ", event.Payload.Before.ID)
			if event.Payload.Before != nil {
				_, err := collection.DeleteOne(context.Background(), event.Payload.Before.ID)
				if err != nil {
					log.Printf("failed to delete soccerprediction: %v\n", err)
				} else {
					log.Println("SoccerPrediction deleted successfully")
				}
			}
		default:
			log.Println("Unknown operation: ", ops)
		}
	}
}
