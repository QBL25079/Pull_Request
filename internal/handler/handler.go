package handler

import (
	"database/sql"
	"net/http"
	"pr-reviewer-service/internal/domain"
	"pr-reviewer-service/internal/service"
	"strings"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	userService service.UserServiceProvider
}

func NewHandler(userService service.UserServiceProvider) *Handler {
	return &Handler{
		userService: userService,
	}
}

func (h *Handler) CreateTeam(c echo.Context) error {
	var req struct {
		TeamName string `json:"team_name"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}
	if req.TeamName == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Team name is required"})
	}

	err := h.userService.CreateTeam(c.Request().Context(), req.TeamName)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") ||
			strings.Contains(err.Error(), "uniq_team_name") { // зависит от имени индекса
			return c.JSON(http.StatusConflict, map[string]string{"error": "Team name already exists"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create team"})
	}

	return c.JSON(http.StatusCreated, map[string]string{
		"message":   "Team created successfully",
		"team_name": req.TeamName,
	})
}

// Создание пользователя
func (h *Handler) CreateUser(c echo.Context) error {
	var user domain.User
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	err := h.userService.CreateUser(c.Request().Context(), user)
	if err != nil {
		if strings.Contains(err.Error(), "violates foreign key constraint") {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Team does not exist"})
		}
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return c.JSON(http.StatusConflict, map[string]string{"error": "User ID already exists"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create user"})
	}

	return c.JSON(http.StatusCreated, map[string]string{
		"message": "User created successfully",
		"user_id": user.UserID,
	})
}

func (h *Handler) GetUserByID(c echo.Context) error {
	userID := c.Param("userID")
	user, err := h.userService.GetUserByID(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get user"})
	}
	if user == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
	}
	return c.JSON(http.StatusOK, user)
}

func (h *Handler) UpdateTeamName(c echo.Context) error {
	oldName := c.Param("teamName")

	var req struct {
		NewName string `json:"newname"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}
	if req.NewName == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "New name is required"})
	}

	err := h.userService.UpdateTeamName(c.Request().Context(), oldName, req.NewName)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Team not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update team name"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Team name updated successfully"})
}

func (h *Handler) CreatePullRequest(c echo.Context) error {
	var pr domain.PullRequest
	if err := c.Bind(&pr); err != nil { // Исправлено: &pr, а не &req
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	err := h.userService.CreatePullRequest(c.Request().Context(), pr)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create pull request"})
	}

	return c.JSON(http.StatusCreated, map[string]string{
		"message":         "Pull Request created successfully",
		"pull_request_id": pr.PullRequestID,
	})
}

func (h *Handler) GetPullRequestByID(c echo.Context) error {
	prID := c.Param("prID")

	pr, err := h.userService.GetPullRequestByID(c.Request().Context(), prID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get pull request"})
	}
	if pr == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Pull Request not found"})
	}

	return c.JSON(http.StatusOK, pr)
}
