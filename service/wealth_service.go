package service

import (
	"context"
	"errors"
	"time"

	"github.com/banking-superapp/wealth-service/model"
	"github.com/banking-superapp/wealth-service/repository"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var (
	ErrSchemeNotFound = errors.New("scheme not found")
	ErrUnauthorized   = errors.New("unauthorized")
)

type WealthService interface {
	GetCatalogue(ctx context.Context, category string) ([]model.MFScheme, error)
	CreateSIP(ctx context.Context, userID string, req *model.CreateSIPRequest) (*model.SIP, error)
	GetPortfolio(ctx context.Context, userID string) (*model.Portfolio, error)
	GetPortfolioAnalytics(ctx context.Context, userID string) (*model.PortfolioAnalytics, error)
	AssessRiskProfile(ctx context.Context, userID string, req *model.RiskProfileRequest) (*model.RiskProfile, error)
	GetRiskProfile(ctx context.Context, userID string) (*model.RiskProfile, error)
}

type wealthService struct {
	mfRepo      repository.MFSchemeRepo
	sipRepo     repository.SIPRepo
	portRepo    repository.PortfolioRepo
	riskRepo    repository.RiskProfileRepo
}

func NewWealthService(mr repository.MFSchemeRepo, sr repository.SIPRepo, pr repository.PortfolioRepo, rr repository.RiskProfileRepo) WealthService {
	return &wealthService{mr, sr, pr, rr}
}

func (s *wealthService) GetCatalogue(ctx context.Context, category string) ([]model.MFScheme, error) {
	return s.mfRepo.FindAll(ctx, category)
}

func (s *wealthService) CreateSIP(ctx context.Context, userID string, req *model.CreateSIPRequest) (*model.SIP, error) {
	oid, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return nil, ErrUnauthorized
	}

	scheme, err := s.mfRepo.FindByCode(ctx, req.SchemeCode)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrSchemeNotFound
		}
		return nil, err
	}

	startDate := req.StartDate
	if startDate.IsZero() {
		startDate = time.Now().AddDate(0, 1, 0)
	}

	sip := &model.SIP{
		UserID:      oid,
		SchemeCode:  req.SchemeCode,
		SchemeName:  scheme.SchemeName,
		Amount:      req.Amount,
		Frequency:   req.Frequency,
		StartDate:   startDate,
		NextSIPDate: startDate,
		Status:      "active",
	}

	if err := s.sipRepo.Create(ctx, sip); err != nil {
		return nil, err
	}
	return sip, nil
}

func (s *wealthService) GetPortfolio(ctx context.Context, userID string) (*model.Portfolio, error) {
	oid, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return nil, ErrUnauthorized
	}

	portfolio, err := s.portRepo.FindByUserID(ctx, oid)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return &model.Portfolio{
				UserID:      oid,
				Holdings:    []model.Holding{},
				TotalValue:  0,
				TotalReturn: 0,
				ReturnPct:   0,
			}, nil
		}
		return nil, err
	}
	return portfolio, nil
}

func (s *wealthService) GetPortfolioAnalytics(ctx context.Context, userID string) (*model.PortfolioAnalytics, error) {
	portfolio, err := s.GetPortfolio(ctx, userID)
	if err != nil {
		return nil, err
	}

	var totalInvested, currentValue, gainLoss float64
	categoryBreakdown := make(map[string]float64)

	for _, h := range portfolio.Holdings {
		totalInvested += h.InvestedValue
		currentValue += h.CurrentValue
		gainLoss += h.GainLoss
	}

	retPct := 0.0
	if totalInvested > 0 {
		retPct = (gainLoss / totalInvested) * 100
	}

	topHoldings := portfolio.Holdings
	if len(topHoldings) > 5 {
		topHoldings = topHoldings[:5]
	}

	return &model.PortfolioAnalytics{
		TotalInvested:     totalInvested,
		CurrentValue:      currentValue,
		TotalGainLoss:     gainLoss,
		ReturnPct:         retPct,
		CategoryBreakdown: categoryBreakdown,
		TopHoldings:       topHoldings,
	}, nil
}

func (s *wealthService) AssessRiskProfile(ctx context.Context, userID string, req *model.RiskProfileRequest) (*model.RiskProfile, error) {
	oid, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return nil, ErrUnauthorized
	}

	// Calculate score from answers
	score := 0
	for _, a := range req.Answers {
		score += a
	}

	category := "conservative"
	mix := map[string]int{"equity": 20, "debt": 70, "hybrid": 10}
	if score > 20 {
		category = "moderate"
		mix = map[string]int{"equity": 50, "debt": 40, "hybrid": 10}
	}
	if score > 30 {
		category = "aggressive"
		mix = map[string]int{"equity": 70, "debt": 20, "hybrid": 10}
	}

	rp := &model.RiskProfile{
		UserID:         oid,
		Score:          score,
		RiskCategory:   category,
		RecommendedMix: mix,
	}

	if err := s.riskRepo.Upsert(ctx, rp); err != nil {
		return nil, err
	}
	return rp, nil
}

func (s *wealthService) GetRiskProfile(ctx context.Context, userID string) (*model.RiskProfile, error) {
	oid, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return nil, ErrUnauthorized
	}
	rp, err := s.riskRepo.FindByUserID(ctx, oid)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return &model.RiskProfile{
				UserID:       oid,
				RiskCategory: "not_assessed",
			}, nil
		}
		return nil, err
	}
	return rp, nil
}
