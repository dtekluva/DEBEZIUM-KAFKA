package types

import "time"

type MobidTrackerEvent struct {
	Payload Payload `json:"payload"`
}
type MobidTracker struct {
	ID                     int       `json:"id" bson:"id"`
	ClickID                string    `json:"click_id" bson:"click_id"`
	PhoneNumber            string    `json:"phone_number" bson:"phone_number"`
	Converted              bool      `json:"converted" bson:"converted"`
	AmountPlayed           float64   `json:"amount_played" bson:"amount_played"`
	Source                 string    `json:"source" bson:"source"`
	DateCreated            time.Time `json:"date_created" bson:"date_created"`
	LastUpdated            time.Time `json:"last_updated" bson:"last_updated"`
	GameType               string    `json:"game_type" bson:"game_type"`
	NumberOfRenewals       int       `json:"number_of_renewals" bson:"number_of_renewals"`
	AmountPaid             float64   `json:"amount_paid" bson:"amount_paid"`
	TelcoNetwork           any       `json:"telco_network" bson:"telco_network"`
	PostbackSent           bool      `json:"postback_sent" bson:"postback_sent"`
	PostbackResponse       string    `json:"postback_response" bson:"postback_response"`
	Controled              bool      `json:"controled" bson:"controled"`
	Unsubscribed           bool      `json:"unsubscribed" bson:"unsubscribed"`
	DateUnsubscribed       any       `json:"date_unsubscribed" bson:"date_unsubscribed"`
	NumberOfUnsubsriptions int       `json:"number_of_unsubsriptions" bson:"number_of_unsubsriptions"`
}
type Source struct {
	Version   string `json:"version"`
	Connector string `json:"connector"`
	Name      string `json:"name"`
	TsMs      int64  `json:"ts_ms"`
	Snapshot  string `json:"snapshot"`
	Db        string `json:"db"`
	Sequence  string `json:"sequence"`
	Schema    string `json:"schema"`
	Table     string `json:"table"`
	TxID      int    `json:"txId"`
	Lsn       int64  `json:"lsn"`
	Xmin      any    `json:"xmin"`
}
type Payload struct {
	Before      *MobidTracker `json:"before"`
	After       *MobidTracker `json:"after"`
	Source      Source        `json:"source"`
	Op          string        `json:"op"`
	TsMs        int64         `json:"ts_ms"`
	Transaction any           `json:"transaction"`
}
