package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"go_consumer_service/types"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

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

// helper function for optional int fields
func intPtr(i int) *int {
	return &i
}

// helper function for optional bool fields
func boolPtr(b bool) *bool {
	return &b
}

// helper function for optional string fields
func strPtr(s string) *string {
	return &s
}

type Product struct {
	ServiceName         string  `json:"service_name"`
	ServiceID           string  `json:"service_id"`
	ProductName         string  `json:"product_name"`
	Amount              float64 `json:"amount"`
	ProductID           string  `json:"product_id"`
	LineNumber          *int    `json:"line_number,omitempty"`
	Band                *int    `json:"band,omitempty"`
	IsDailySubscription *bool   `json:"is_daily_subscription,omitempty"`
	Network             *string `json:"network,omitempty"`
}

// SecureDAndUpstreamServiceAndProductDetails looks up product by ID
func SecureDAndUpstreamServiceAndProductDetails(productID string) *Product {
	data := map[string]Product{
		"23410220000024641": {ServiceName: "WYSE_CASH_150", ServiceID: "234102200006769", ProductName: "WYSE_CASH", Amount: 150, ProductID: "23410220000024641"},
		"23410220000024642": {ServiceName: "WYSE_CASH_200", ServiceID: "234102200006769", ProductName: "WYSE_CASH", Amount: 200, ProductID: "23410220000024642"},
		"23410220000024643": {ServiceName: "WYSE_CASH_300", ServiceID: "234102200006769", ProductName: "WYSE_CASH", Amount: 300, ProductID: "23410220000024643", Band: intPtr(10000)},
		"23410220000024644": {ServiceName: "WYSE_CASH_400", ServiceID: "234102200006769", ProductName: "WYSE_CASH", Amount: 400, ProductID: "23410220000024644"},
		"23410220000024645": {ServiceName: "WYSE_CASH_500", ServiceID: "234102200006769", ProductName: "WYSE_CASH", Amount: 500, ProductID: "23410220000024645"},
		"23410220000024646": {ServiceName: "WYSE_CASH_700", ServiceID: "234102200006769", ProductName: "WYSE_CASH", Amount: 700, ProductID: "23410220000024646"},
		"23410220000024647": {ServiceName: "WYSE_CASH_1000", ServiceID: "234102200006769", ProductName: "WYSE_CASH", Amount: 1000, ProductID: "23410220000024647"},
		"23410220000024635": {ServiceName: "INSTANT_CASH_150", ServiceID: "234102200006767", ProductName: "INSTANT_CASH", Amount: 150, ProductID: "23410220000024635", LineNumber: intPtr(1)},
		"23410220000024636": {ServiceName: "INSTANT_CASH_300", ServiceID: "234102200006767", ProductName: "INSTANT_CASH", Amount: 300, ProductID: "23410220000024636", LineNumber: intPtr(1)},
		"23410220000024637": {ServiceName: "INSTANT_CASH_500", ServiceID: "234102200006767", ProductName: "INSTANT_CASH", Amount: 500, ProductID: "23410220000024637", LineNumber: intPtr(2)},
		"23410220000024638": {ServiceName: "INSTANT_CASH_750", ServiceID: "234102200006767", ProductName: "INSTANT_CASH", Amount: 750, ProductID: "23410220000024638", LineNumber: intPtr(3)},
		"23410220000024639": {ServiceName: "INSTANT_CASH_900", ServiceID: "234102200006767", ProductName: "INSTANT_CASH", Amount: 900, ProductID: "23410220000024639", LineNumber: intPtr(4)},
		"23410220000024640": {ServiceName: "INSTANT_CASH_1000", ServiceID: "234102200006767", ProductName: "INSTANT_CASH", Amount: 1000, ProductID: "23410220000024640", LineNumber: intPtr(5)},
		"23410220000024656": {ServiceName: "SOCCER_CASH_100", ServiceID: "234102200006772", ProductName: "SOCCER_CASH", Amount: 100, ProductID: "23410220000024656", LineNumber: intPtr(1)},
		"23410220000024657": {ServiceName: "SOCCER_CASH_500", ServiceID: "234102200006772", ProductName: "SOCCER_CASH", Amount: 500, ProductID: "23410220000024657", LineNumber: intPtr(2)},
		"23410220000024658": {ServiceName: "SOCCER_CASH_800", ServiceID: "234102200006772", ProductName: "SOCCER_CASH", Amount: 800, ProductID: "23410220000024658", LineNumber: intPtr(3)},
		"23410220000024659": {ServiceName: "SOCCER_CASH_1200", ServiceID: "234102200006772", ProductName: "SOCCER_CASH", Amount: 1200, ProductID: "23410220000024659", LineNumber: intPtr(4)},
		"23410220000024648": {ServiceName: "SALARY_FOR_LIFE_200", ServiceID: "234102200006770", ProductName: "SALARY_FOR_LIFE", Amount: 200, ProductID: "23410220000024648", LineNumber: intPtr(1)},
		"23410220000024649": {ServiceName: "SALARY_FOR_LIFE_500", ServiceID: "234102200006770", ProductName: "SALARY_FOR_LIFE", Amount: 500, ProductID: "23410220000024649", LineNumber: intPtr(2)},
		"23410220000024650": {ServiceName: "SALARY_FOR_LIFE_800", ServiceID: "234102200006770", ProductName: "SALARY_FOR_LIFE", Amount: 800, ProductID: "23410220000024650", LineNumber: intPtr(3)},
		"23410220000024651": {ServiceName: "SALARY_FOR_LIFE_1200", ServiceID: "234102200006770", ProductName: "SALARY_FOR_LIFE", Amount: 1200, ProductID: "23410220000024651", LineNumber: intPtr(7)},
		"23410220000027462": {ServiceName: "INSTANT_CASHOUT_50", ServiceID: "234102200006961", ProductName: "INSTANT_CASHOUT", Amount: 50, ProductID: "23410220000027462", LineNumber: intPtr(1), IsDailySubscription: boolPtr(true)},
		"23410220000027463": {ServiceName: "INSTANT_CASHOUT_75", ServiceID: "234102200006961", ProductName: "INSTANT_CASHOUT", Amount: 75, ProductID: "23410220000027463", LineNumber: intPtr(1), IsDailySubscription: boolPtr(true)},
		"23410220000027470": {ServiceName: "WYSE_CASH_50", ServiceID: "234102200006965", ProductName: "WYSE_CASH", Amount: 50, ProductID: "23410220000027470", LineNumber: intPtr(1), IsDailySubscription: boolPtr(true)},
		"23410220000027471": {ServiceName: "WYSE_CASH_75", ServiceID: "234102200006965", ProductName: "WYSE_CASH", Amount: 75, ProductID: "23410220000027471", LineNumber: intPtr(1), IsDailySubscription: boolPtr(true)},
		"23410220000027464": {ServiceName: "SALARY_FOR_LIFE_50", ServiceID: "234102200006962", ProductName: "SALARY_FOR_LIFE", Amount: 50, ProductID: "23410220000027464", LineNumber: intPtr(1), IsDailySubscription: boolPtr(true)},
		"23410220000027465": {ServiceName: "SALARY_FOR_LIFE_75", ServiceID: "234102200006962", ProductName: "SALARY_FOR_LIFE", Amount: 75, ProductID: "23410220000027465", LineNumber: intPtr(1), IsDailySubscription: boolPtr(true)},
		"23410220000027466": {ServiceName: "FAST_FINGER_50", ServiceID: "234102200006963", ProductName: "FAST_FINGER", Amount: 50, ProductID: "23410220000027466", LineNumber: intPtr(1), IsDailySubscription: boolPtr(true)},
		"23410220000027467": {ServiceName: "FAST_FINGER_75", ServiceID: "234102200006963", ProductName: "FAST_FINGER", Amount: 75, ProductID: "23410220000027467", LineNumber: intPtr(1), IsDailySubscription: boolPtr(true)},
		"23410220000027468": {ServiceName: "SOCCER_CASH_50", ServiceID: "234102200006964", ProductName: "SOCCER_CASH", Amount: 75, ProductID: "23410220000027468", LineNumber: intPtr(1), IsDailySubscription: boolPtr(true)},
		"23410220000027469": {ServiceName: "SOCCER_CASH_75", ServiceID: "234102200006964", ProductName: "SOCCER_CASH", Amount: 100, ProductID: "23410220000027469", LineNumber: intPtr(1), IsDailySubscription: boolPtr(true)},
		"0017182000001707":  {ServiceName: "INSTANT_CASH_50", ServiceID: "234102200006964", ProductName: "INSTANT_CASH", Amount: 50, ProductID: "0017182000001707", LineNumber: intPtr(1), IsDailySubscription: boolPtr(true), Network: strPtr("GLO")},
		"0017182000003867":  {ServiceName: "FAST_FINGER_50", ServiceID: "234102200006964", ProductName: "FAST_FINGER", Amount: 50, ProductID: "0017182000003867", LineNumber: intPtr(1), IsDailySubscription: boolPtr(true), Network: strPtr("GLO")},
		"0017182000003868":  {ServiceName: "SALARY_FOR_LIFE_50", ServiceID: "234102200006964", ProductName: "SALARY_FOR_LIFE", Amount: 50, ProductID: "0017182000003868", LineNumber: intPtr(1), IsDailySubscription: boolPtr(true), Network: strPtr("GLO")},
		"0017182000003869":  {ServiceName: "WYSE_CASH_50", ServiceID: "234102200006964", ProductName: "WYSE_CASH", Amount: 50, ProductID: "0017182000003869", LineNumber: intPtr(1), IsDailySubscription: boolPtr(true), Network: strPtr("GLO")},
		"0017182000003870":  {ServiceName: "SOCCER_CASH_50", ServiceID: "234102200006964", ProductName: "SOCCER_CASH", Amount: 50, ProductID: "0017182000003870", LineNumber: intPtr(1), IsDailySubscription: boolPtr(true), Network: strPtr("GLO")},
	}

	if val, ok := data[productID]; ok {
		return &val
	}
	return nil
}

type SecureDDataDumpData struct {
	Msisdn      string `json:"msisdn"`
	Activation  string `json:"activation"`
	ProductID   string `json:"productID"`
	Description string `json:"description"`
	Timestamp   string `json:"timestamp"` // keep as string, parse later
	TrxId       string `json:"trxId"`
}

func (u *Utils) SendMarketingPartnersPostback(instanceId int) {
	// sending postcallback slack message
	// message := fmt.Sprintf("Sending marketing partners postback for instance id: %d", instanceId)
	// go SendSlackNotification(message)
	// get the secure datadump model
	filter := bson.M{"id": instanceId}
	secureDDumpCol := u.db.Collection("secure_data_dump")
	// find one
	var secureDDump types.SecureDDataDump
	err := secureDDumpCol.FindOne(context.Background(), filter).Decode(&secureDDump)
	if err != nil {
		log.Println("Error finding secure data dump: ", err.Error())
		return
	}
	// get the data
	var data SecureDDataDumpData
	raw := secureDDump.Data

	// If the string starts with "b'" and ends with "'", strip them
	if strings.HasPrefix(raw, "b'") && strings.HasSuffix(raw, "'") {
		raw = raw[2 : len(raw)-1]
	}
	fmt.Printf(raw)

	// unquoted, err := strconv.Unquote(`"` + raw + `"`)
	// if err != nil {
	// 	log.Println("Error unquoting secure data dump data:", err)
	// 	return
	// }

	// if err := json.Unmarshal([]byte(unquoted), &data); err != nil {
	// 	log.Println("Error unmarshalling secure data dump data:", err)
	// 	return
	// }
	clean := strings.ReplaceAll(raw, `\n`, "\n")
	clean = strings.ReplaceAll(clean, `\t`, "\t")

	// Step 3: Now unmarshal
	if err := json.Unmarshal([]byte(clean), &data); err != nil {
		log.Println("Error unmarshalling secure data dump data:", err)
		return
	}
	// let's see the data
	log.Println("Data: ", data)
	transRef := data.TrxId
	phoneNumer := data.Msisdn
	activation := data.Activation
	description := data.Description
	productId := data.ProductID
	var subscriptionAmount any
	var gameType string
	source := secureDDump.Source
	if source == "ST_GLO" {
		nitroswitchServiceCodes := SubscriptionServiceCodes()
		nitroswitchEquivalentServiceID := GetEquivalentProductCode()
		productCode, ok := nitroswitchEquivalentServiceID[productId]
		if !ok {
			productCode = ""
		}
		subscription := nitroswitchServiceCodes[productCode]
		subscriptionAmount = subscription.Price
		gameType = subscription.Name
	} else if source == "MTN" {
		game_details := SecureDAndUpstreamServiceAndProductDetails(productId)
		subscriptionAmount = game_details.Amount
		gameType = game_details.ProductName
	} else {
		subscriptionAmount = nil
		gameType = ""

	}
	dataDump := map[string]interface{}{
		"transRef":           transRef,
		"phoneNumer":         phoneNumer,
		"activation":         activation,
		"description":        description,
		"subscriptionAmount": subscriptionAmount,
		"gameType":           gameType,
		"source":             source,
		"networkProvider":    source,
	}
	u.runAndSendTrafficPostback(transRef, dataDump)
	log.Println("Sent marketing partners postback for instance id: ", instanceId)

}

func (u *Utils) runAndSendTrafficPostback(reference string, dataDumps map[string]interface{}) {
	ctx := context.Background()
	var phoneNumber, gameType, networkProvider string
	var subscriptionAmount float64

	if dataDumps != nil {
		phoneNumber, _ = dataDumps["phone_number"].(string)
		reference, _ = dataDumps["reference"].(string)
		if val, ok := dataDumps["subscription_amount"].(float64); ok {
			subscriptionAmount = val
		}
		gameType, _ = dataDumps["game_type"].(string)
		networkProvider, _ = dataDumps["network_provider"].(string)
	} else {
		var instance types.SecureDTransaction
		err := u.db.Collection("secure_d_transaction").
			FindOne(ctx, bson.M{"reference": reference}).
			Decode(&instance)
		if err != nil {
			fmt.Errorf("transaction not found: %w", err)
			return
		}
		phoneNumber = instance.PhoneNumber
		reference = instance.Reference
		subscriptionAmount = instance.SubscriptionAmount
		gameType = *instance.GameType
	}

	isWinwise := strings.HasPrefix(reference, "WINWISE")
	sent := false
	postbackResponse := ""
	postbackSkip := u.PostSkipDecisioning()

	var source string

	// ======================= LOGIC TREE ==========================

	if strings.HasPrefix(reference, "a1") {
		source = "ANGEL_MEDIA"
		if !postbackSkip {
			resp := sendAdsTrackerPostbackURL(source, reference, subscriptionAmount)
			sent, postbackSkip, postbackResponse = true, false, resp
		} else {
			sent, postbackResponse = false, "pass"
		}
	} else if strings.HasPrefix(reference, "aazz") && len(reference) == 55 {
		source = "MORBIDTEK_MEDIA"
		if !postbackSkip {
			resp := sendAdsTrackerPostbackURL(source, reference, subscriptionAmount)
			postbackResponse = resp
			if resp != "Failed" {
				var parsed map[string]interface{}
				if json.Unmarshal([]byte(resp), &parsed) == nil {
					if parsed["error"] == float64(0) && parsed["info"] == "Conversion Received." {
						sent, postbackSkip = true, false
					}
				}
			}
		} else {
			sent, postbackResponse = false, "pass"
		}
	} else if strings.Contains(reference, "mobitech") {
		reference = reference[8:]
		source = "MORBIDTEK_MEDIA"
		if !postbackSkip {
			resp := sendAdsTrackerPostbackURL(source, reference, subscriptionAmount)
			postbackResponse = resp
			if resp != "Failed" {
				var parsed map[string]interface{}
				if json.Unmarshal([]byte(resp), &parsed) == nil {
					if parsed["error"] == float64(0) && parsed["info"] == "Conversion Received." {
						sent, postbackSkip = true, false
					}
				}
			}
		} else {
			sent, postbackResponse = false, "pass"
		}
	} else if len(reference) == 32 {
		source = "MOBPLUS"
		if !postbackSkip {
			resp := sendAdsTrackerPostbackURL(source, reference, subscriptionAmount)
			postbackResponse = resp
			if resp == "success" {
				sent, postbackSkip = true, false
			} else {
				sent = false
			}
		} else {
			sent, postbackResponse = false, "pass"
		}
	} else if (strings.Contains(reference, ",") || strings.Contains(reference, "%")) && len(reference) > 29 {
		source = "TRAFFIC_COMPANY"
		if postbackSkip {
			sent, postbackResponse = false, "pass"
		} else {
			resp := sendAdsTrackerPostbackURL(source, reference, subscriptionAmount)
			postbackResponse, sent, postbackSkip = resp, true, false
		}
	} else if strings.Contains(reference, "phoenix") || strings.HasPrefix(reference, "sci_") || strings.Contains(reference, "_832") {
		source = "Phoenix_MTN"
		if !postbackSkip {
			resp := sendAdsTrackerPostbackURL(source, reference, subscriptionAmount)
			postbackResponse = resp
			var parsed map[string]interface{}
			if json.Unmarshal([]byte(resp), &parsed) == nil {
				if parsed["code"] == float64(0) && parsed["msg"] == "report success" {
					sent, postbackSkip = true, false
				}
			}
		} else {
			sent, postbackResponse = false, "pass"
		}
	} else {
		if strings.HasPrefix(reference, "5") {
			source = "GOLDEN_GOOSE"
			if !postbackSkip {
				resp := sendAdsTrackerPostbackURL(source, reference, subscriptionAmount)
				postbackResponse = resp
				if resp == "Received" {
					sent, postbackSkip = true, false
				}
			} else {
				sent, postbackResponse = false, "pass"
			}
		} else if strings.Contains(reference, "vic") {
			source = "VICTORY_ADS"
			if !postbackSkip {
				reference = reference[3:]
				resp := sendAdsTrackerPostbackURL(source, reference, subscriptionAmount)
				postbackResponse = resp
				var parsed map[string]interface{}
				if json.Unmarshal([]byte(resp), &parsed) == nil {
					if parsed["message"] == "conversion recorded successfully" {
						sent = true
					}
				}
			} else {
				sent, postbackResponse = false, "pass"
			}
		} else if strings.Contains(reference, "upstream_paidMTN") {
			source = "UPSTREAM"
		} else if strings.HasPrefix(reference, "golden") {
			source = "GOLDEN_GOOSE"
			if !postbackSkip {
				reference = reference[6:]
				resp := sendAdsTrackerPostbackURL(source, reference, subscriptionAmount)
				postbackResponse, sent, postbackSkip = resp, true, false
			} else {
				sent, postbackResponse = false, "pass"
			}
		} else if isWinwise {
			source = "WINWISE_UPSELLING"
		} else if len(reference) == 47 {
			source = "CLICK_STREAM"
			if !postbackSkip {
				resp := sendAdsTrackerPostbackURL(source, reference, subscriptionAmount)
				postbackResponse, sent, postbackSkip = resp, true, false
			} else {
				sent, postbackResponse = false, "pass"
			}
		} else if strings.HasPrefix(reference, "9399_") || strings.HasPrefix(reference, "mobipium") {
			source = "MOBIPIUM"
			if strings.HasPrefix(reference, "mobipium") {
				reference = reference[8:]
			}
			if !postbackSkip {
				resp := sendAdsTrackerPostbackURL(source, reference, subscriptionAmount)
				postbackResponse, sent, postbackSkip = resp, true, false
			} else {
				sent, postbackResponse = false, "pass"
			}
		} else if strings.Contains(reference, "CBT") {
			source = "CLICKBYTE_MEDIA"
			reference = reference[3:]
			if !postbackSkip {
				resp := sendAdsTrackerPostbackURL(source, reference, subscriptionAmount)
				postbackResponse, sent, postbackSkip = resp, true, false
			} else {
				sent, postbackResponse = false, "pass"
			}
		} else if strings.Contains(reference, "sigma") {
			source = "SIGMA"
			reference = reference[5:]
			if !postbackSkip {
				resp := sendAdsTrackerPostbackURL(source, reference, subscriptionAmount)
				postbackResponse, sent, postbackSkip = resp, true, false
			} else {
				sent, postbackResponse = false, "pass"
			}
		} else if strings.Contains(reference, "shine") {
			source = "SHINE_DIGITAL"
			reference = reference[5:]
			if !postbackSkip {
				resp := sendAdsTrackerPostbackURL(source, reference, subscriptionAmount)
				postbackResponse, sent, postbackSkip = resp, true, false
			}
		} else if strings.Contains(reference, "greencorp") {
			source = "GREEN_CORP"
			reference = reference[9:]
			if !postbackSkip {
				resp := sendAdsTrackerPostbackURL(source, reference, subscriptionAmount)
				postbackResponse, sent, postbackSkip = resp, true, false
			} else {
				sent, postbackResponse = false, "pass"
			}
		} else {
			source = "ORGANIC"
		}
	}
	// send update post requets
	requestData := map[string]interface{}{
		"postback_sent":     sent,
		"postback_response": postbackResponse,
		"click_id":          reference,
		"controlled":        postbackSkip,
		"converted":         false,
		"amount_played":     subscriptionAmount,
		"phone_number":      phoneNumber,
		"game_type":         gameType,
		"telco_network":     networkProvider,
		"source":            source,
	}
	recordUrl := "https://libertydraw.com/api/v1/ads-tracker/morbid-data-callback"

	headers := map[string]interface{}{
		"Content-Type": "application/json",
	}
	req := webutils.NewRequest("POST", recordUrl, requestData, headers, ctx)
	response, err := req.Send()
	if err != nil {
		log.Println("Error sending postback: ", err.Error())
		return
	}
	if response.StatusCode != 200 {
		log.Println("Error sending postback: ", response.StatusCode)
		return
	}
	log.Println("Postback sent successfully")
}

func sendAdsTrackerPostbackURL(source, clickID string, amount float64) string {
	var url string

	switch source {
	case "MORBIDTEK_MEDIA":
		url = fmt.Sprintf("http://mobtekmedia.hopb0.com/notify/110122/?click_id=%s", clickID)
	case "MOBPLUS":
		url = fmt.Sprintf("http://m.mobplus.net/c/p/74a3eac380d944aab100c0f6112b268f?txid=%s", clickID)
	case "AD_MAVEN":
		url = fmt.Sprintf("https://xml.realtime-bid.com/conversion?c=%s&count=1&value=%v", clickID, amount)
	case "ANGEL_MEDIA":
		url = fmt.Sprintf("http://postback.rustclick.com/pb/377?click_id=%s&payout=", clickID)
	case "TRAFFIC_COMPANY":
		url = fmt.Sprintf("https://postback.level23.nl/?currency=USD&handler=11342&hash=14f1e1b8e3a711e0c49d672857610ea6&tracker=%s", clickID)
	case "GOLDEN_GOOSE":
		url = fmt.Sprintf("http://n.gg.agency/ntf1/?token=af54c5e8739f11b00bca46da47962675&click_id=%s", clickID)
	case "Phoenix_MTN":
		url = fmt.Sprintf("https://adwh.mywkd.com/thirdparty?event_name=Subscribe&phx_click_id=%s&product_id=832", clickID)
	case "VICTORY_ADS":
		url = fmt.Sprintf("http://mnoi.online/trackwsp.php?subid=%s", clickID)
	case "GOLDEN_GOLDEN":
		url = fmt.Sprintf("http://n.gg.agency/ntf1/?token=af54c5e8739f11b00bca46da47962675&click_id=%s", clickID)
	case "CLICK_STREAM":
		url = fmt.Sprintf("https://track.goforclicks.swaarm-clients.com/postback?click_id=%s&security_token=4fe0f47d-81b5-4077-9b0a-b2a402969743", clickID)
	case "MOBIPIUM":
		url = fmt.Sprintf("https://smobipiumlink.com/conversion/index.php?jp=%s&source=WINWISE", clickID)
	case "CLICKBYTE_MEDIA":
		url = fmt.Sprintf("https://click.clickbyte-media.com/postback?cid=%s&payout=%v", clickID, amount)
	case "SIGMA":
		url = fmt.Sprintf("https://cd.sigmamobi.com/callbacks/223?request_id=%s&event=lead", clickID)
	case "SHINE_DIGITAL":
		url = fmt.Sprintf("http://shinedigitalworld.offerstrack.net/advBack.php?click_id=%s", clickID)
	case "GREEN_CORP":
		url = fmt.Sprintf("https://pfwsxng7jj55bf7q4n6v7egmlu0sizta.lambda-url.eu-west-2.on.aws/?aid=3702176&tid=134377&visitor_id=%s&payout=0.50", clickID)
	default:
		return "Failed"
	}

	// HTTP client with 2s timeout
	client := http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return "Failed"
	}
	defer resp.Body.Close()

	fmt.Println("Postback response: ", resp.StatusCode, " for source: ", source)

	// Special cases
	if source == "VICTORY_ADS" {
		body, _ := ioutil.ReadAll(resp.Body)
		return string(body)
	} else if source == "CLICK_STREAM" || source == "MOBIPIUM" || source == "CLICKBYTE_MEDIA" || source == "SHINE_DIGITAL" {
		return fmt.Sprintf("%d", resp.StatusCode)
	}

	// Default: return response body text
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "Failed"
	}
	return string(body)
}
