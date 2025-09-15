package types

import "time"

type SecureDDataDump struct {
	ID                int       `bson:"id,omitempty" json:"id"`
	Data              string    `bson:"data" json:"data"`
	IPAddress         string    `bson:"ip_address,omitempty" json:"ip_address"`
	Source            string    `bson:"source" json:"source"`
	PostbackProcessed *bool     `bson:"postback_processed,omitempty" json:"postback_processed"`
	CreatedAt         time.Time `bson:"created_at,omitempty" json:"created_at"`
	UpdatedAt         time.Time `bson:"updated_at,omitempty" json:"updated_at"`
}

type SecureDTransaction struct {
	ID                 int       `bson:"id,omitempty" json:"id"`
	PhoneNumber        string    `bson:"phone_number" json:"phone_number"`
	Reference          string    `bson:"reference" json:"reference"`
	SubscriptionAmount float64   `bson:"subscription_amount" json:"subscription_amount"`
	GameType           *string   `bson:"game_type,omitempty" json:"game_type,omitempty"`
	TransactionStatus  string    `bson:"transaction_status" json:"transaction_status"`
	IsSuccessful       bool      `bson:"is_successful" json:"is_successful"`
	Activation         bool      `bson:"activation" json:"activation"`
	Payload            string    `bson:"payload" json:"payload"`
	CreatedAt          time.Time `bson:"created_at,omitempty" json:"created_at"`
	UpdatedAt          time.Time `bson:"updated_at,omitempty" json:"updated_at"`
}

type CampaignSkipTable struct {
	ID            int        `bson:"id,omitempty" json:"id"`
	Source        *string    `bson:"source,omitempty" json:"source,omitempty"`
	GameType      *string    `bson:"game_type,omitempty" json:"game_type,omitempty"`
	TelcoNetwork  *string    `bson:"telco_network,omitempty" json:"telco_network,omitempty"`
	SkipInterval  int        `bson:"skip_interval" json:"skip_interval"`
	LastSkipStart *time.Time `bson:"last_skip_start,omitempty" json:"last_skip_start,omitempty"`
	LastSkipEnd   *time.Time `bson:"last_skip_end,omitempty" json:"last_skip_end,omitempty"`
	DateCreated   time.Time  `bson:"date_created,omitempty" json:"date_created"`
	LastUpdated   time.Time  `bson:"last_updated,omitempty" json:"last_updated"`
}
