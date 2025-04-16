package handler

import (
	"fmt"
	"net/http"
	"pvz-test/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func (h *Handler) RemoveLastItem(c *gin.Context) {
	role, ok := c.Get(roleCtx)
	if !ok || role != models.RoleEmployee {
		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
		return
	}

	pvzIDStr := c.Param(pvzIdParam)
	if pvzIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing pvzId parameter"})
		return
	}

	pvzID, err := uuid.Parse(pvzIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("failed to delete item: %s", err.Error())})
		return
	}

	logrus.Infof("start to delete last item from: %s", pvzID)
	err = h.services.Reception.DeleteItem(pvzID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to delete item: %s", err.Error())})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "item deleted successfully"})

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
