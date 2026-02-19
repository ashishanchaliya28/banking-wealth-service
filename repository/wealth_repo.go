package repository

import (
	"context"
	"time"

	"github.com/banking-superapp/wealth-service/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type MFSchemeRepo interface {
	FindAll(ctx context.Context, category string) ([]model.MFScheme, error)
	FindByCode(ctx context.Context, code string) (*model.MFScheme, error)
}

type SIPRepo interface {
	Create(ctx context.Context, s *model.SIP) error
	FindByUserID(ctx context.Context, userID bson.ObjectID) ([]model.SIP, error)
}

type PortfolioRepo interface {
	FindByUserID(ctx context.Context, userID bson.ObjectID) (*model.Portfolio, error)
	Upsert(ctx context.Context, p *model.Portfolio) error
}

type RiskProfileRepo interface {
	FindByUserID(ctx context.Context, userID bson.ObjectID) (*model.RiskProfile, error)
	Upsert(ctx context.Context, rp *model.RiskProfile) error
}

type mfSchemeRepo struct{ col *mongo.Collection }
type sipRepo struct{ col *mongo.Collection }
type portfolioRepo struct{ col *mongo.Collection }
type riskProfileRepo struct{ col *mongo.Collection }

func NewMFSchemeRepo(db *mongo.Database) MFSchemeRepo   { return &mfSchemeRepo{col: db.Collection("mf_schemes")} }
func NewSIPRepo(db *mongo.Database) SIPRepo             { return &sipRepo{col: db.Collection("sips")} }
func NewPortfolioRepo(db *mongo.Database) PortfolioRepo  { return &portfolioRepo{col: db.Collection("portfolios")} }
func NewRiskProfileRepo(db *mongo.Database) RiskProfileRepo { return &riskProfileRepo{col: db.Collection("risk_profiles")} }

func (r *mfSchemeRepo) FindAll(ctx context.Context, category string) ([]model.MFScheme, error) {
	filter := bson.M{"is_active": true}
	if category != "" {
		filter["category"] = category
	}
	cursor, err := r.col.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var schemes []model.MFScheme
	cursor.All(ctx, &schemes)
	return schemes, nil
}

func (r *mfSchemeRepo) FindByCode(ctx context.Context, code string) (*model.MFScheme, error) {
	var s model.MFScheme
	err := r.col.FindOne(ctx, bson.M{"scheme_code": code}).Decode(&s)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *sipRepo) Create(ctx context.Context, s *model.SIP) error {
	s.CreatedAt = time.Now()
	s.UpdatedAt = time.Now()
	_, err := r.col.InsertOne(ctx, s)
	return err
}

func (r *sipRepo) FindByUserID(ctx context.Context, userID bson.ObjectID) ([]model.SIP, error) {
	cursor, err := r.col.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var sips []model.SIP
	cursor.All(ctx, &sips)
	return sips, nil
}

func (r *portfolioRepo) FindByUserID(ctx context.Context, userID bson.ObjectID) (*model.Portfolio, error) {
	var p model.Portfolio
	err := r.col.FindOne(ctx, bson.M{"user_id": userID}).Decode(&p)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *portfolioRepo) Upsert(ctx context.Context, p *model.Portfolio) error {
	p.UpdatedAt = time.Now()
	_, err := r.col.UpdateOne(ctx,
		bson.M{"user_id": p.UserID},
		bson.M{"$set": p},
		&mongo.UpdateOptions{Upsert: boolPtr(true)},
	)
	return err
}

func (r *riskProfileRepo) FindByUserID(ctx context.Context, userID bson.ObjectID) (*model.RiskProfile, error) {
	var rp model.RiskProfile
	err := r.col.FindOne(ctx, bson.M{"user_id": userID}).Decode(&rp)
	if err != nil {
		return nil, err
	}
	return &rp, nil
}

func (r *riskProfileRepo) Upsert(ctx context.Context, rp *model.RiskProfile) error {
	rp.AssessedAt = time.Now()
	_, err := r.col.UpdateOne(ctx,
		bson.M{"user_id": rp.UserID},
		bson.M{"$set": rp},
		&mongo.UpdateOptions{Upsert: boolPtr(true)},
	)
	return err
}

func boolPtr(b bool) *bool { return &b }
