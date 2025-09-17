package types

import "time"

type SecureDataDumpEvent struct {
	Payload SecureDataDumpPayload `json:"payload"`
}

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

type SecureDataDumpSource struct {
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

type SecureDataDumpPayload struct {
	Before      *SecureDDataDump     `json:"before"`
	After       *SecureDDataDump     `json:"after"`
	Source      SecureDataDumpSource `json:"source"`
	Op          string               `json:"op"`
	TsMs        int64                `json:"ts_ms"`
	Transaction any                  `json:"transaction"`
}

type DataDumpData struct {
	Msisdn      string `json:"msisdn"`
	Activation  string `json:"activation"`
	ProductID   string `json:"productID"`
	Description string `json:"description"`
	Timestamp   string `json:"timestamp"` // keep as string, parse later
	TrxId       string `json:"trxId"`
}

type SecureDTransactionEvent struct {
	Payload SecureDTransactionPayload `json:"payload"`
}

type SecureDTransactionPayload struct {
	Before      *SecureDTransaction  `json:"before"`
	After       *SecureDTransaction  `json:"after"`
	Source      SecureDataDumpSource `json:"source"`
	Op          string               `json:"op"`
	TsMs        int64                `json:"ts_ms"`
	Transaction any                  `json:"transaction"`
}

type SoccerPrediction struct {
	ID                  int64  `bson:"id,omitempty" json:"id,omitempty"`
	UserProfileID       *int64 `bson:"user_profile_id,omitempty" json:"user_profile_id,omitempty"`
	AgentID             *int64 `bson:"agent_id,omitempty" json:"agent_id,omitempty"`
	FootballTableID     *int64 `bson:"football_table_id,omitempty" json:"football_table_id,omitempty"`
	BoughtLotteryTicket *int64 `bson:"bought_lottery_ticket_id,omitempty" json:"bought_lottery_ticket_id,omitempty"`

	Phone            string     `bson:"phone" json:"phone"`
	GameID           *string    `bson:"game_id,omitempty" json:"game_id,omitempty"`
	HomeChoice       *int       `bson:"home_choice,omitempty" json:"home_choice,omitempty"`
	AwayChoice       *int       `bson:"away_choice,omitempty" json:"away_choice,omitempty"`
	BandPlayed       *string    `bson:"band_played,omitempty" json:"band_played,omitempty"`
	StakeAmount      float64    `bson:"stake_amount" json:"stake_amount"`
	PotentialWinning float64    `bson:"potential_winning" json:"potential_winning"`
	AmountPaid       float64    `bson:"amount_paid" json:"amount_paid"`
	Paid             bool       `bson:"paid" json:"paid"`
	GameFixtureID    *string    `bson:"game_fixture_id,omitempty" json:"game_fixture_id,omitempty"`
	AccountNo        *string    `bson:"account_no,omitempty" json:"account_no,omitempty"`
	BankName         *string    `bson:"bank_name,omitempty" json:"bank_name,omitempty"`
	BankCode         *string    `bson:"bank_code,omitempty" json:"bank_code,omitempty"`
	Date             time.Time  `bson:"date" json:"date"`
	PaidDate         *time.Time `bson:"paid_date,omitempty" json:"paid_date,omitempty"`
	Channel          string     `bson:"channel" json:"channel"`
	Freemium         bool       `bson:"freemium" json:"freemium"`
	IsDrawn          bool       `bson:"is_drawn" json:"is_drawn"`
	ScoreChecked     bool       `bson:"score_checked" json:"score_checked"`
	Won              bool       `bson:"won" json:"won"`
	Active           bool       `bson:"active" json:"active"`
	IsScoreValid     bool       `bson:"is_score_valid" json:"is_score_valid"`
	PlayType         string     `bson:"play_type" json:"play_type"`
	LotteryType      string     `bson:"lottery_type" json:"lottery_type"`
}

type SoccerPredictionPayload struct {
	Before      *SoccerPrediction    `json:"before"`
	After       *SoccerPrediction    `json:"after"`
	Source      SecureDataDumpSource `json:"source"`
	Op          string               `json:"op"`
	TsMs        int64                `json:"ts_ms"`
	Transaction any                  `json:"transaction"`
}

type SoccerPredictionEvent struct {
	Payload SoccerPredictionPayload `json:"payload"`
}
