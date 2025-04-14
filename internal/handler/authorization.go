package handler

import (
	"net/http"
	"pvz-test/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *Handler) Login(c *gin.Context) {
	var input models.LoginRequest

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusUnauthorized, `binding error`)
		return
	}

	token, err := h.services.Authorization.Login(input)
	if err != nil {
		logrus.Info(err)
		newErrorResponse(c, http.StatusUnauthorized, `Unauthorized`)
		return
	}

	logrus.Info("user logged: ", input.Email)
	c.JSON(http.StatusOK, token)
}

func (h *Handler) DummyLogin(c *gin.Context) {
	var req models.DummyLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		newErrorResponse(c, http.StatusBadRequest, `binding error`)
		return
	}

	if err := h.validate.Struct(&req); err != nil {
		newErrorResponse(c, http.StatusBadRequest, `role validate error`)
		return
	}

	token, err := h.services.DummyLogin(req.Role)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, token)
}

func (h *Handler) Register(c *gin.Context) {
	var input models.RegisterRequest

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.validate.Struct(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.services.Authorization.Register(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	logrus.Info("user created: ", input.Email)
	c.JSON(http.StatusCreated, models.UserResponse{
		ID:    user.ID,
		Email: user.Email,
		Role:  user.Role,
	})
}
