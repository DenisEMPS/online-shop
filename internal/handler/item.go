package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/DenisEMPS/online-shop/internal/domain"
	"github.com/DenisEMPS/online-shop/internal/service"
	"github.com/gin-gonic/gin"
)

func (h *Handler) CreateProduct(c *gin.Context) {
	var product domain.CreateProduct

	err := c.Bind(&product)
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, "invalid request params")
		return
	}

	id, err := h.service.Product.Create(&product)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "internal error")
		return
	}

	c.JSON(http.StatusOK, map[string]int64{
		"id": id,
	})
}

func (h *Handler) GetProductByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, "invalid request")
		return
	}

	product, err := h.service.Product.GetByID(int64(id))
	if err != nil {
		if errors.Is(err, service.ErrProductNotExists) {
			NewErrorResponse(c, http.StatusNotFound, "product not founded")
			return
		}
		NewErrorResponse(c, http.StatusInternalServerError, "internal error")
		return
	}

	c.JSON(http.StatusOK, product)
}
