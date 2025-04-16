package handler

import (
	"fmt"
	"log"
	"net/http"
	"pvz-test/internal/models"

	"github.com/gin-gonic/gin"
)

func (h *Handler) CreatePVZ(c *gin.Context) {
	var req models.PVZRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Wrong request format"})
		return
	}

	if err := h.validate.Struct(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Wrong city format"})
		return
	}

	userRole, _ := c.Get(roleCtx)
	if userRole != models.RoleModerator {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only moderators can create PVZ"})
		return
	}

	newPVZ, err := h.services.CreatePvz(req.City)
	if err != nil {
		log.Println("failed to create pvz:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Pvz creation error"})
		return
	}

	c.JSON(http.StatusCreated, newPVZ)
}

func (h *Handler) GetPVZList(c *gin.Context) {
	userRole, _ := c.Get(roleCtx)
	if userRole != models.RoleModerator && userRole != models.RoleEmployee {
		c.JSON(http.StatusUnauthorized, fmt.Errorf("unauthorized"))
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
