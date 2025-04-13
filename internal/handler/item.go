package handler

import (
	"fmt"
	"net/http"
	"pvz-test/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *Handler) RemoveLastItem(c *gin.Context) {
	role, ok := c.Get(roleCtx)
	if !ok || role != models.RoleEmployee {
		c.JSON(http.StatusForbidden, fmt.Errorf("role verification error"))
		return
	}

	pvzIDStr := c.Query(pvzIdParam)
	pvzID, err := uuid.Parse(pvzIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, fmt.Errorf("%s parse error", pvzIdParam))
		return
	}

	err = h.services.Reception.DeleteItem(pvzID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Errorf("reception close error: %s", err.Error()))
		return
	}
	c.JSON(http.StatusOK, "OK")

}
func (h *Handler) AddItem(c *gin.Context) {
	role, ok := c.Get(roleCtx)
	if !ok || role != models.RoleEmployee {
		c.JSON(http.StatusForbidden, fmt.Errorf("role verification error"))
		return
	}

	var req models.AddProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	item, err := h.services.Reception.AddItem(req.PvzID, req.Type)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, item)

}
