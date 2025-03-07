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

	api := r.Group("/api")
	{
		item := api.Group("/item")
		{
			item.GET("/:id", h.GetItemByID)
			item.POST("/", h.CreateItem)
		}
	}

	return r
}
