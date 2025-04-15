package handler

import (
	"fmt"
	"net/http"
	"pvz-test/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *Handler) CreateReception(c *gin.Context) {
	role, ok := c.Get(roleCtx)
	if !ok || role != models.RoleEmployee {
		c.JSON(http.StatusForbidden, fmt.Errorf("role verification error"))
		return
	}

	var recReq models.CreateReceptionRequest

	if err := c.BindJSON(&recReq); err != nil {
		newErrorResponse(c, http.StatusForbidden, `binding error`)
		return
	}

	reception, err := h.services.CreateReception(recReq.PvzID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, reception)
}

func (h *Handler) CloseReception(c *gin.Context) {
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

	reception, err := h.services.Reception.CloseActiveReception(pvzID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Errorf("reception close error: %s", err.Error()))
		return
	}
	if (reception == models.Reception{}) {
		c.JSON(http.StatusInternalServerError, fmt.Errorf("no any active reception for pvz"))
		return
	}

	c.JSON(http.StatusOK, reception)
}
