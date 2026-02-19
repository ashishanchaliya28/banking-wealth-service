package handler

import (
	"errors"

	"github.com/banking-superapp/wealth-service/model"
	"github.com/banking-superapp/wealth-service/service"
	"github.com/gofiber/fiber/v2"
)

type WealthHandler struct{ svc service.WealthService }

func NewWealthHandler(svc service.WealthService) *WealthHandler { return &WealthHandler{svc: svc} }

func (h *WealthHandler) GetCatalogue(c *fiber.Ctx) error {
	category := c.Query("category")
	schemes, err := h.svc.GetCatalogue(c.Context(), category)
	if err != nil {
		return respond(c, fiber.StatusInternalServerError, nil, err.Error())
	}
	return respond(c, fiber.StatusOK, schemes, "")
}

func (h *WealthHandler) CreateSIP(c *fiber.Ctx) error {
	userID := c.Get("X-User-ID")
	var req model.CreateSIPRequest
	if err := c.BodyParser(&req); err != nil {
		return respond(c, fiber.StatusBadRequest, nil, "invalid request body")
	}
	sip, err := h.svc.CreateSIP(c.Context(), userID, &req)
	if err != nil {
		if errors.Is(err, service.ErrSchemeNotFound) {
			return respond(c, fiber.StatusNotFound, nil, err.Error())
		}
		return respond(c, fiber.StatusInternalServerError, nil, err.Error())
	}
	return respond(c, fiber.StatusCreated, sip, "")
}

func (h *WealthHandler) GetPortfolio(c *fiber.Ctx) error {
	userID := c.Get("X-User-ID")
	portfolio, err := h.svc.GetPortfolio(c.Context(), userID)
	if err != nil {
		return respond(c, fiber.StatusInternalServerError, nil, err.Error())
	}
	return respond(c, fiber.StatusOK, portfolio, "")
}

func (h *WealthHandler) GetPortfolioAnalytics(c *fiber.Ctx) error {
	userID := c.Get("X-User-ID")
	analytics, err := h.svc.GetPortfolioAnalytics(c.Context(), userID)
	if err != nil {
		return respond(c, fiber.StatusInternalServerError, nil, err.Error())
	}
	return respond(c, fiber.StatusOK, analytics, "")
}

func (h *WealthHandler) AssessRiskProfile(c *fiber.Ctx) error {
	userID := c.Get("X-User-ID")
	var req model.RiskProfileRequest
	if err := c.BodyParser(&req); err != nil {
		return respond(c, fiber.StatusBadRequest, nil, "invalid request body")
	}
	rp, err := h.svc.AssessRiskProfile(c.Context(), userID, &req)
	if err != nil {
		return respond(c, fiber.StatusInternalServerError, nil, err.Error())
	}
	return respond(c, fiber.StatusOK, rp, "")
}

func (h *WealthHandler) GetRiskProfile(c *fiber.Ctx) error {
	userID := c.Get("X-User-ID")
	rp, err := h.svc.GetRiskProfile(c.Context(), userID)
	if err != nil {
		return respond(c, fiber.StatusInternalServerError, nil, err.Error())
	}
	return respond(c, fiber.StatusOK, rp, "")
}

func respond(c *fiber.Ctx, status int, data interface{}, errMsg string) error {
	if errMsg != "" {
		return c.Status(status).JSON(fiber.Map{"success": false, "error": errMsg})
	}
	return c.Status(status).JSON(fiber.Map{"success": true, "data": data})
}
