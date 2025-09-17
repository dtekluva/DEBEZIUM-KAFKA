package utils

import (
	"context"
	"go_consumer_service/types"
	"log"

	"github.com/Ayobami6/webutils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Utils struct {
	db *mongo.Database
}

func NewUtils(db *mongo.Database) *Utils {
	return &Utils{
		db: db,
	}
}

// utils functions
func (u *Utils) GetLottoSubscribersPhoneList(ctx context.Context, phoneNumber string, gameType *string) ([]string, error) {
	var allPhones []string
	baseFilter := bson.M{"paid": true, "phone": phoneNumber}
	if gameType != nil {
		lottoFilter := bson.M{"paid": true, "phone": phoneNumber, "lottery_type": *gameType}
		phones, err := FetchPhones(ctx, u.db.Collection("lotto_event"), lottoFilter)
		if err != nil {
			return nil, err
		}
		allPhones = append(allPhones, phones...)

	} else {
		// lotto filter without gameType
		lottoFilter := bson.M{"paid": true, "phone": phoneNumber, "channel": "USSD"}
		phones, err := FetchPhones(ctx, u.db.Collection("lotto_event"), lottoFilter)
		if err != nil {
			return nil, err
		}
		allPhones = append(allPhones, phones...)
	}

	filter := baseFilter
	if gameType != nil {
		filter = bson.M{"paid": true, "phone": phoneNumber, "lottery_type": *gameType}
	}
	phones, err := FetchPhones(ctx, u.db.Collection("lottery_model"), filter)
	if err != nil {
		return nil, err
	}
	allPhones = append(allPhones, phones...)
	// Soccer prediction
	phones, err = FetchPhones(ctx, u.db.Collection("soccer_prediction"), filter)
	if err != nil {
		return nil, err
	}
	allPhones = append(allPhones, phones...)
	// AwoofGameTable
	phones, err = FetchPhones(ctx, u.db.Collection("awoof_game_table"), filter)
	if err != nil {
		return nil, err
	}
	allPhones = append(allPhones, phones...)
	return allPhones, nil
}

func (u *Utils) PostSkipDecisioning() bool {
	var status bool
	// get the last record of the constant table
	filter := bson.M{}
	// get the last record of the constant table
	opts := options.FindOne().SetSort(bson.D{{"created_at", -1}})
	// get the last record of the constant table
	result := u.db.Collection("constant").FindOne(context.Background(), filter, opts)
	// map to constanttable type
	var constantTable types.ConstantTable
	// map to constanttable type
	err := result.Decode(&constantTable)
	// map to constanttable type
	if err != nil {
		return false
	}
	postbackCount := constantTable.PostbackCount
	postbackSkipCount := constantTable.PostbackSkipCount
	postbackToSend := constantTable.PostbackToSend
	postbackToSkip := constantTable.PostbackToSkip
	if postbackCount <= postbackToSend {
		constantTable.PostbackCount = postbackCount + 1
		status = false
	} else {
		if postbackSkipCount <= postbackToSkip {
			constantTable.PostbackSkipCount = postbackSkipCount + 1
			status = true
		} else {
			constantTable.PostbackCount = 1
			constantTable.PostbackSkipCount = 1
			status = false
		}
	}
	// resave the newly updated the constant table
	_, err = u.db.Collection("constant").UpdateOne(context.Background(), bson.M{}, bson.M{"$set": constantTable})
	return status

}

// other standalone utils
func FetchPhones(ctx context.Context, col *mongo.Collection, filter interface{}) ([]string, error) {
	cursor, err := col.Find(ctx, filter, options.Find().SetProjection(bson.M{"phone": 1}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var results []struct {
		Phone string `bson:"phone"`
	}
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	phones := make([]string, 0, len(results))
	for _, result := range results {
		phones = append(phones, result.Phone)
	}
	return phones, nil
}

// Subscription holds details of a service plan
type Subscription struct {
	Name     string
	Key      string
	Price    int
	Validity int // corrected from "validty"
}

// SubscriptionServiceCodes returns all subscription codes mapped to their details
func SubscriptionServiceCodes() map[string]Subscription {
	return map[string]Subscription{
		"2076": {Name: "EVERWAGE DAILY AUTO", Key: "AWD", Price: 100, Validity: 1},
		"2077": {Name: "EVERWAGE DAILY ONETIME", Key: "AWD OT", Price: 100, Validity: 1},
		"2078": {Name: "EVERWAGE WEEKLY AUTO", Key: "AWW", Price: 200, Validity: 7},
		"2079": {Name: "EVERWAGE WEEKLY ONETIME", Key: "AWW OT", Price: 200, Validity: 7},
		"2080": {Name: "EVERWAGE MONTHLY AUTO", Key: "AWM", Price: 500, Validity: 30},
		"2081": {Name: "EVERWAGE MONTHLY ONETIME", Key: "AWM OT", Price: 500, Validity: 30},

		"2145": {Name: "AI FORTUNE DAILY AUTO", Key: "AFD", Price: 100, Validity: 1},
		"2146": {Name: "AI FORTUNE DAILY ONETIME", Key: "AFD OT", Price: 100, Validity: 1},
		"2147": {Name: "AI FORTUNE WEEKLY AUTO", Key: "AFW", Price: 200, Validity: 7},
		"2148": {Name: "AI FORTUNE WEEKLY ONETIME", Key: "AFW OT", Price: 100, Validity: 7},
		"2149": {Name: "AI FORTUNE MONTHLY AUTO", Key: "AWM", Price: 500, Validity: 30},
		"2150": {Name: "AI FORTUNE MONTHLY ONETIME", Key: "AWM OT", Price: 500, Validity: 30},

		"2139": {Name: "SMS-TO-AI DAILY AUTO", Key: "ASD", Price: 100, Validity: 1},
		"2140": {Name: "SMS-TO-AI DAILY ONETIME", Key: "ASD OT", Price: 100, Validity: 1},
		"2141": {Name: "SMS-TO-AI WEEKLY AUTO", Key: "ASW", Price: 200, Validity: 7},
		"2142": {Name: "SMS-TO-AI WEEKLY ONETIME", Key: "ASW OT", Price: 100, Validity: 1},
		"2143": {Name: "SMS-TO-AI MONTHLY AUTO", Key: "ASM", Price: 500, Validity: 30},
		"2144": {Name: "SMS-TO-AI MONTHLY ONETIME", Key: "ASM OT", Price: 500, Validity: 30},
	}
}

func GetEquivalentProductCode() map[string]string {
	return map[string]string{
		"1000005720": "2145",
		"1000005726": "2146",
		"1000005729": "2149",
		"1000005730": "2150",
		"1000005727": "2147",
		"1000005728": "2148",
		"1000005718": "2076",
		"1000005731": "2077",
		"1000005736": "2080",
		"1000005737": "2081",
		"1000005734": "2078",
		"1000005735": "2079",
		"1000005716": "2139",
		"1000005721": "2140",
		"1000005724": "2143",
		"1000005725": "2144",
		"1000005722": "2141",
		"1000005723": "2142",
	}
}

func SendSlackNotification(message string) {
	ctx := context.Background()
	// slack webhook url
	webhookUrl := webutils.GetEnv("SLACK_WEBHOOK_URL", "https://hooks.slack.com/services/T02G5QZJQ/B02G5QZJQ/XXXXXXXXXXXXXXXXXXXXXXXX")
	headers := map[string]interface{}{
		"Content-Type": "application/json",
	}
	body := map[string]interface{}{
		"text": message,
	}
	newReq := webutils.NewRequest("POST", webhookUrl, body, headers, ctx)
	//
	response, err := newReq.Send()
	if err != nil {
		log.Println("Error sending slack notification: ", err.Error())
		return
	}
	if response.StatusCode != 200 {
		log.Println("Error sending slack notification: ", response.Status)
		return
	}
	log.Println("Slack notification sent successfully")
}
