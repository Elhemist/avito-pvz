package handler

import (
	"pvz-test/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

type Handler struct {
	services *service.Service
	validate *validator.Validate
}

const pvzIdParam = "pvz_id"

func NewHandler(services *service.Service) *Handler {
	return &Handler{
		services: services,
		validate: validator.New(),
	}
}

func (h *Handler) InitRoutes() *gin.Engine {

	router := gin.New()

	api := router.Group("/api")
	{
		api.POST("/register", h.Register)
		api.POST("/login", h.Login)
		api.POST("/dummyLogin", h.DummyLogin)

		api.Use(h.JWTMiddleware())
		{
			api.POST("/pvz", h.CreatePVZ)
			api.GET("/pvz", h.GetPVZList)

			api.POST("/receptions", h.CreateReception)
			api.POST("/pvz/:pvzId/close_last_reception", h.CloseReception)

			api.POST("/products", h.AddItem)
			api.POST("/pvz/:pvzId/delete_last_product", h.RemoveLastItem)
		}
	}

	return router
}
