package handler

import (
	"net/http"
	"pr-reviewer-service/internal/service"
	"strings"

	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	userServise service.UserServiceProvider
}

func NewUserHandler(userservice service.UserServiceProvider) *UserHandler {
	return &UserHandler{userServise: userservice}
}

func (h UserHandler) CreateTeam(c echo.Context) error {
	var req struct {
		TeamName string `json:team_name`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}
	if req.TeamName == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Team name is required"})
	}

	err := h.userServise.CreateTeam(c.Request().Context(), req.TeamName)

	if err != nil {
		if strings.Contains(err.Error(), "duplicatekey value violates uniq constraint") {
			return c.JSON(http.StatusConflict, map[string]string{"error": "Team name already exists"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, map[string]string{"message": "Team created successfully", "team_name": req.TeamName})
}

func (h *Handler) CreateUser(c echo.Context) error {

}
