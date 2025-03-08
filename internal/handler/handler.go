package handler

import (
	"github.com/DenisEMPS/online-shop/internal/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) InitRoutes() *gin.Engine {
	r := gin.Default()

	auth := r.Group("/auth")
	{
		auth.POST("/sign-up", h.RegisterNewUser)
		auth.POST("/sign-in", h.LoginUser)
	}

	api := r.Group("/api", h.UserIdentity)
	{
		item := api.Group("/item")
		{
			item.GET("/:id", h.GetProductByID)
			item.POST("/", h.CreateProduct)
		}
	}

	return r
}
