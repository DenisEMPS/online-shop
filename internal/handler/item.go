package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/DenisEMPS/online-shop/internal/domain"
	"github.com/gin-gonic/gin"
)

func (h *Handler) GetItemByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, "invalid request params")
		return
	}

	item, err := h.service.Item.GetByID(int64(id))
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "internal error")
		fmt.Println(err)
		return
	}

	c.JSON(http.StatusOK, item)
}

func (h *Handler) CreateItem(c *gin.Context) {
	var item domain.CreateItem

	err := c.Bind(&item)
	if err != nil {
		fmt.Println(err)
		NewErrorResponse(c, http.StatusBadRequest, "invalid request params")
		return
	}

	id, err := h.service.Item.Create(&item)
	if err != nil {
		fmt.Println(err)
		NewErrorResponse(c, http.StatusInternalServerError, "internal error")
		return
	}

	c.JSON(http.StatusOK, map[string]int64{
		"id": id,
	})
}
