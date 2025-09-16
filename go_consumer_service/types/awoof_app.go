package types

import "time"

type AwoofGameTableEvent struct {
	Payload AwoofGameTablePayload `json:"payload"`
}

type AwoofGameTable struct {
	ID             int64  `bson:"id,omitempty" json:"id,omitempty"`
	ItemID         *int64 `bson:"item_id,omitempty" json:"item_id,omitempty"`
	UserProfileID  *int64 `bson:"user_profile_id,omitempty" json:"user_profile_id,omitempty"`
	AgentProfileID *int64 `bson:"agent_profile_id,omitempty" json:"agent_profile_id,omitempty"`

	Percentage  float64 `bson:"percentage" json:"percentage"`
	Phone       *string `bson:"phone,omitempty" json:"phone,omitempty"`
	AwoofAmount float64 `bson:"awoof_amount" json:"awoof_amount"`

	RTO                                float64 `bson:"rto" json:"rto"`
	RTP                                float64 `bson:"rtp" json:"rtp"`
	RTPPer                             float64 `bson:"rtp_per" json:"rtp_per"`
	UssdTelcoCommission                float64 `bson:"ussd_telco_commission" json:"ussd_telco_commission"`
	UssdTelcoCommissionValue           float64 `bson:"ussd_telco_commission_value" json:"ussd_telco_commission_value"`
	UssdTelcoAggregatorCommission      float64 `bson:"ussd_telco_aggregator_commission" json:"ussd_telco_aggregator_commission"`
	UssdTelcoAggregatorCommissionValue float64 `bson:"ussd_telco_aggregator_commission_value" json:"ussd_telco_aggregator_commission_value"`

	AwoofStakeAmount     float64 `bson:"awoof_stake_amount" json:"awoof_stake_amount"`
	Band                 float64 `bson:"band" json:"band"`
	TicketPrice          float64 `bson:"ticket_price" json:"ticket_price"`
	TicketPercent        float64 `bson:"ticket_percent" json:"ticket_percent"`
	AmountPaid           float64 `bson:"amount_paid" json:"amount_paid"`
	InstantCashoutAmount float64 `bson:"instant_cashout_amount" json:"instant_cashout_amount"`
	IllusionAmount       float64 `bson:"illusion_amount" json:"illusion_amount"`

	Paid                  bool       `bson:"paid" json:"paid"`
	AccountNo             *string    `bson:"account_no,omitempty" json:"account_no,omitempty"`
	BankName              *string    `bson:"bank_name,omitempty" json:"bank_name,omitempty"`
	LotteryType           string     `bson:"lottery_type" json:"lottery_type"` // default "AWOOF"
	PaidDate              *time.Time `bson:"paid_date,omitempty" json:"paid_date,omitempty"`
	GamePlayID            *string    `bson:"game_play_id,omitempty" json:"game_play_id,omitempty"`
	Channel               string     `bson:"channel" json:"channel"`
	TicketID              *string    `bson:"ticket_id,omitempty" json:"ticket_id,omitempty"`
	UniqueTicketID        *string    `bson:"unique_ticket_id,omitempty" json:"unique_ticket_id,omitempty"`
	GamePin               *string    `bson:"game_pin,omitempty" json:"game_pin,omitempty"`
	Payload               *string    `bson:"payload,omitempty" json:"payload,omitempty"`
	Consent               bool       `bson:"consent" json:"consent"`
	CountPercentage       int        `bson:"count_percentage" json:"count_percentage"`
	IsAgent               bool       `bson:"is_agent" json:"is_agent"`
	IsDuplicate           bool       `bson:"is_duplicate" json:"is_duplicate"`
	DateCreated           time.Time  `bson:"date_created" json:"date_created"`
	DateUpdated           time.Time  `bson:"date_updated" json:"date_updated"`
	ICashIsPaid           bool       `bson:"i_cash_is_paid" json:"i_cash_is_paid"`
	PlayedViaTelcoChannel bool       `bson:"played_via_telco_channel" json:"played_via_telco_channel"`
	TelcoChannel          string     `bson:"telco_channel" json:"telco_channel"`
	ServiceType           string     `bson:"service_type" json:"service_type"`
	AIRequest             *string    `bson:"ai_request,omitempty" json:"ai_request,omitempty"`
	AIResponse            *string    `bson:"ai_response,omitempty" json:"ai_response,omitempty"`
	IsUtilized            bool       `bson:"is_utilized" json:"is_utilized"`
	TelcoNetwork          *string    `bson:"telco_network,omitempty" json:"telco_network,omitempty"`
}

type AwoofGameTablePayload struct {
	Before      *AwoofGameTable `json:"before"`
	After       *AwoofGameTable `json:"after"`
	Source      EventSource     `json:"source"`
	Op          string          `json:"op"`
	TsMs        int64           `json:"ts_ms"`
	Transaction any
}
