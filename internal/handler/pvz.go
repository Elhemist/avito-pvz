package handler

import (
	"fmt"
	"log"
	"net/http"
	"pvz-test/internal/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *Handler) CreatePVZ(c *gin.Context) {
	var req models.PVZRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат запроса"})
		return
	}

	if err := h.validate.Struct(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный город"})
		return
	}

	userRole, _ := c.Get(roleCtx)
	if userRole != models.RoleModerator {
		c.JSON(http.StatusConflict, gin.H{"error": "Доступ запрещен. Только модератор может создавать ПВЗ"})
		return
	}

	newPVZ := &models.PVZ{
		ID:               uuid.New(),
		RegistrationDate: time.Now(),
		City:             req.City,
	}

	c.JSON(http.StatusCreated, newPVZ)
}

func (h *Handler) GetPVZList(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, fmt.Errorf("unauthorized"))
		return
	}
	u, ok := user.(models.User)
	if !ok || (u.Role != "employee" && u.Role != "moderator") {
		c.JSON(http.StatusForbidden, fmt.Errorf("forbidden"))
		return
	}

	var q models.GetPVZListQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		c.JSON(http.StatusBadRequest, fmt.Errorf("invalid query parameters"))
		return
	}
	if q.Page < 1 {
		q.Page = 1
	}
	if q.Limit < 1 || q.Limit > 30 {
		q.Limit = 10
	}
	offset := (q.Page - 1) * q.Limit

	pvzList, err := h.services.Pvz.GetFilteredPVZ(q.StartDate, q.EndDate, q.Limit, offset)
	if err != nil {
		log.Println("failed to fetch pvz list:", err)
		c.JSON(http.StatusInternalServerError, fmt.Errorf("failed to fetch pvz list"))
		return
	}

	c.JSON(http.StatusOK, pvzList)
}
