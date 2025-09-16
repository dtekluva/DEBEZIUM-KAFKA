package types

import (
	"time"
)

type LottoTicketEvent struct {
	Payload LottoTicketEventPayload `json:"payload"`
}

type LottoTicket struct {
	ID                               int64     `bson:"id,omitempty" json:"id,omitempty"`
	UserProfileID                    int64     `bson:"user_profile_id" json:"user_profile_id"`
	AgentProfileID                   *int64    `bson:"agent_profile_id,omitempty" json:"agent_profile_id,omitempty"`
	BatchID                          *int64    `bson:"batch_id,omitempty" json:"batch_id,omitempty"`
	Phone                            string    `bson:"phone" json:"phone"`
	StakeAmount                      float64   `bson:"stake_amount" json:"stake_amount"`
	PotentialWinning                 float64   `bson:"potential_winning" json:"potential_winning"`
	ExpectedAmount                   float64   `bson:"expected_amount" json:"expected_amount"`
	AmountPaid                       float64   `bson:"amount_paid" json:"amount_paid"`
	Illusion                         float64   `bson:"illusion" json:"illusion"`
	Rto                              float64   `bson:"rto" json:"rto"`
	Rtp                              float64   `bson:"rtp" json:"rtp"`
	RtpPer                           float64   `bson:"rtp_per" json:"rtp_per"`
	EffectiveRtp                     float64   `bson:"effective_rtp" json:"effective_rtp"`
	CommissionPer                    float64   `bson:"commission_per" json:"commission_per"`
	CommissionValue                  float64   `bson:"commission_value" json:"commission_value"`
	SalaryForLifeJackpotPer          float64   `bson:"salary_for_life_jackpot_per" json:"salary_for_life_jackpot_per"`
	SalaryForLifeJackpotAmount       float64   `bson:"salary_for_life_jackpot_amount" json:"salary_for_life_jackpot_amount"`
	WinCommissionPer                 float64   `bson:"win_commission_per" json:"win_commission_per"`
	WinCommissionValue               float64   `bson:"win_commission_value" json:"win_commission_value"`
	UssdTelcoCommission              float64   `bson:"ussd_telco_commission" json:"ussd_telco_commission"`
	UssdTelcoCommissionValue         float64   `bson:"ussd_telco_commission_value" json:"ussd_telco_commission_value"`
	UssdTelcoAggregatorCommission    float64   `bson:"ussd_telco_aggregator_commission" json:"ussd_telco_aggregator_commission"`
	UssdTelcoAggregatorCommissionVal float64   `bson:"ussd_telco_aggregator_commission_value" json:"ussd_telco_aggregator_commission_value"`
	Paid                             bool      `bson:"paid" json:"paid"`
	Date                             time.Time `bson:"date" json:"date"`
	UpdatedAt                        time.Time `bson:"updated_at" json:"updated_at"`
	NumberOfTicket                   int       `bson:"number_of_ticket" json:"number_of_ticket"`
	Channel                          string    `bson:"channel" json:"channel"`
	GamePlayID                       *string   `bson:"game_play_id,omitempty" json:"game_play_id,omitempty"`
	UniqueGamePlayID                 *string   `bson:"unique_game_play_id,omitempty" json:"unique_game_play_id,omitempty"`
	AwoofGamePlayID                  *string   `bson:"awoof_game_play_id,omitempty" json:"awoof_game_play_id,omitempty"`
	LotteryType                      string    `bson:"lottery_type" json:"lottery_type"`
	LotterySource                    string    `bson:"lottery_source" json:"lottery_source"`
	GameType                         string    `bson:"game_type" json:"game_type"`
	ServiceType                      string    `bson:"service_type" json:"service_type"`
	HasInterest                      bool      `bson:"has_interest" json:"has_interest"`
	Ticket                           string    `bson:"ticket" json:"ticket"`
	SystemGeneratedNum               *string   `bson:"system_generated_num,omitempty" json:"system_generated_num,omitempty"`
	WinCombo                         *string   `bson:"win_combo,omitempty" json:"win_combo,omitempty"`
	IsAgent                          bool      `bson:"is_agent" json:"is_agent"`
	S4lDrawn                         bool      `bson:"s4l_drawn" json:"s4l_drawn"`
	InstantCashoutDrawn              bool      `bson:"instant_cashout_drawn" json:"instant_cashout_drawn"`
	PosInstantCashoutDrawn           bool      `bson:"pos_instant_cashout_drawn" json:"pos_instant_cashout_drawn"`
	IcashCounted                     bool      `bson:"icash_counted" json:"icash_counted"`
	Icash2Counted                    bool      `bson:"icash_2_counted" json:"icash_2_counted"`
	IcashLocalCounted                bool      `bson:"icash_local_counted" json:"icash_local_counted"`
	IsDuplicate                      bool      `bson:"is_duplicate" json:"is_duplicate"`
	GameIdTreated                    bool      `bson:"game_id_treated" json:"game_id_treated"`
	Pin                              *string   `bson:"pin,omitempty" json:"pin,omitempty"`
	IdentityID                       *string   `bson:"identity_id,omitempty" json:"identity_id,omitempty"`
	PlayedViaTelcoChannel            bool      `bson:"played_via_telco_channel" json:"played_via_telco_channel"`
	TelcoChannel                     string    `bson:"telco_channel" json:"telco_channel"`
	IsNewQuikaGame                   bool      `bson:"is_new_quika_game" json:"is_new_quika_game"`
	TelcoNetwork                     *string   `bson:"telco_network,omitempty" json:"telco_network,omitempty"`
	DrawnFor                         string    `bson:"drawn_for" json:"drawn_for"`
	SeederStatus                     string    `bson:"seeder_status" json:"seeder_status"`
	ContentDeliverySmsSent           bool      `bson:"content_delivery_sms_sent" json:"content_delivery_sms_sent"`
	ProductID                        *string   `bson:"product_id,omitempty" json:"product_id,omitempty"`
}

type LottoEventSource struct {
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

type LottoTicketEventPayload struct {
	Before      *LottoTicket     `json:"before"`
	After       *LottoTicket     `json:"after"`
	Source      LottoEventSource `json:"source"`
	Op          string           `json:"op"`
	TsMs        int64            `json:"ts_ms"`
	Transaction any              `json:"transaction"`
}
