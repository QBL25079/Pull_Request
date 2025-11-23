package handler

import (
	"net/http"
	"pr-reviewer-service/internal/domain"
	"pr-reviewer-service/internal/service"

	"github.com/labstack/echo/v4"
)

type TeamHandler struct {
	service *service.TeamService
}

func NewTeamHandler(s *service.TeamService) *TeamHandler {
	return &TeamHandler{service: s}
}

func (h *TeamHandler) CreateTeam(c echo.Context) error {
	t := new(domain.Team)
	if err := c.Bind(t); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid payload"})
	}
	if t.TeamName == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "team_name is required"})
	}

	if err := h.service.CreateTeam(t); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, map[string]string{"message": "team created"})
}

func (h *TeamHandler) ListTeams(c echo.Context) error {
	teams, err := h.service.ListTeams()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, teams)
}
