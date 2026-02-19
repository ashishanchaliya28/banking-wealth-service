package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type MFScheme struct {
	ID           bson.ObjectID `bson:"_id,omitempty" json:"id"`
	SchemeCode   string        `bson:"scheme_code" json:"scheme_code"`
	SchemeName   string        `bson:"scheme_name" json:"scheme_name"`
	AMC          string        `bson:"amc" json:"amc"`
	Category     string        `bson:"category" json:"category"` // equity | debt | hybrid | liquid
	SubCategory  string        `bson:"sub_category" json:"sub_category"`
	NAV          float64       `bson:"nav" json:"nav"`
	NAVDate      time.Time     `bson:"nav_date" json:"nav_date"`
	Returns1Y    float64       `bson:"returns_1y" json:"returns_1y"`
	Returns3Y    float64       `bson:"returns_3y" json:"returns_3y"`
	Returns5Y    float64       `bson:"returns_5y" json:"returns_5y"`
	Risk         string        `bson:"risk" json:"risk"` // low | moderate | high
	MinSIP       float64       `bson:"min_sip" json:"min_sip"`
	MinLumpsum   float64       `bson:"min_lumpsum" json:"min_lumpsum"`
	IsActive     bool          `bson:"is_active" json:"is_active"`
}

type SIP struct {
	ID          bson.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID      bson.ObjectID `bson:"user_id" json:"user_id"`
	SchemeCode  string        `bson:"scheme_code" json:"scheme_code"`
	SchemeName  string        `bson:"scheme_name" json:"scheme_name"`
	Amount      float64       `bson:"amount" json:"amount"`
	Frequency   string        `bson:"frequency" json:"frequency"` // monthly | weekly
	StartDate   time.Time     `bson:"start_date" json:"start_date"`
	NextSIPDate time.Time     `bson:"next_sip_date" json:"next_sip_date"`
	Status      string        `bson:"status" json:"status"` // active | paused | cancelled
	TotalUnits  float64       `bson:"total_units" json:"total_units"`
	TotalAmount float64       `bson:"total_amount" json:"total_amount"`
	CreatedAt   time.Time     `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time     `bson:"updated_at" json:"updated_at"`
}

type Portfolio struct {
	ID          bson.ObjectID   `bson:"_id,omitempty" json:"id"`
	UserID      bson.ObjectID   `bson:"user_id" json:"user_id"`
	Holdings    []Holding       `bson:"holdings" json:"holdings"`
	TotalValue  float64         `bson:"total_value" json:"total_value"`
	TotalReturn float64         `bson:"total_return" json:"total_return"`
	ReturnPct   float64         `bson:"return_pct" json:"return_pct"`
	UpdatedAt   time.Time       `bson:"updated_at" json:"updated_at"`
}

type Holding struct {
	SchemeCode    string  `bson:"scheme_code" json:"scheme_code"`
	SchemeName    string  `bson:"scheme_name" json:"scheme_name"`
	Units         float64 `bson:"units" json:"units"`
	CurrentNAV    float64 `bson:"current_nav" json:"current_nav"`
	CurrentValue  float64 `bson:"current_value" json:"current_value"`
	InvestedValue float64 `bson:"invested_value" json:"invested_value"`
	GainLoss      float64 `bson:"gain_loss" json:"gain_loss"`
}

type RiskProfile struct {
	ID              bson.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID          bson.ObjectID `bson:"user_id" json:"user_id"`
	Score           int           `bson:"score" json:"score"`
	RiskCategory    string        `bson:"risk_category" json:"risk_category"` // conservative | moderate | aggressive
	RecommendedMix  map[string]int `bson:"recommended_mix" json:"recommended_mix"` // {"equity": 60, "debt": 30, "hybrid": 10}
	AssessedAt      time.Time     `bson:"assessed_at" json:"assessed_at"`
}

// Request types
type CreateSIPRequest struct {
	SchemeCode string    `json:"scheme_code"`
	Amount     float64   `json:"amount"`
	Frequency  string    `json:"frequency"`
	StartDate  time.Time `json:"start_date"`
}

type RiskProfileRequest struct {
	Answers []int `json:"answers"` // questionnaire answers (1-5 scale)
}

type LinkExternalRequest struct {
	FolioNumber string `json:"folio_number"`
	PAN         string `json:"pan"`
}

type PortfolioAnalytics struct {
	TotalInvested  float64            `json:"total_invested"`
	CurrentValue   float64            `json:"current_value"`
	TotalGainLoss  float64            `json:"total_gain_loss"`
	ReturnPct      float64            `json:"return_pct"`
	CategoryBreakdown map[string]float64 `json:"category_breakdown"`
	TopHoldings    []Holding          `json:"top_holdings"`
}
