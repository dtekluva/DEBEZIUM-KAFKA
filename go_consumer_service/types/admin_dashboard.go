package types

import "time"

type ConstantTable struct {
	ID                        int64     `bson:"id,omitempty" json:"id,omitempty"`
	SalesTarget               int       `bson:"sales_target" json:"sales_target"`
	CreatedAt                 time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt                 time.Time `bson:"updated_at" json:"updated_at"`
	IsActive                  bool      `bson:"is_active" json:"is_active"`
	WinwiseAgentsSalaryAmount float64   `bson:"winwise_agents_salary_amount" json:"winwise_agents_salary_amount"`
	PostbackCount             int       `bson:"postback_count" json:"postback_count"`
	PostbackSkipCount         int       `bson:"postback_skip_count" json:"postback_skip_count"`
	PostbackToSend            int       `bson:"postback_to_send" json:"postback_to_send"`
	PostbackToSkip            int       `bson:"postback_to_skip" json:"postback_to_skip"`
}

type ConstantTableDebeziumEvent struct {
	Payload ConstantTableEventPayload `json:"payload"`
}

type ConstantTableEventPayload struct {
	Before      *ConstantTable `json:"before"`
	After       *ConstantTable `json:"after"`
	Source      EventSource    `json:"source"`
	Op          string         `json:"op"`
	TsMs        int64          `json:"ts_ms"`
	Transaction any
}
