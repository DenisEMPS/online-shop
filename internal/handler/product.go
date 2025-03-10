package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/DenisEMPS/online-shop/internal/domain"
	"github.com/DenisEMPS/online-shop/internal/filter"
	"github.com/DenisEMPS/online-shop/internal/service"
	"github.com/gin-gonic/gin"
)

const (
	defaultSortBy    = "id"
	defaultSortOrder = "ASC"
	qName            = "name"
	qPrice           = "price"
	qInStock         = "in_stock"
	qCreatedAt       = "created_at"
	qLimit           = "limit"
	qSortBy          = "sort_by"
	qSortOrder       = "sort_order"
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

func (h *Handler) GetAllProducts(c *gin.Context) {
	filterOptions, err := FilterAllProducts(c)
	if err != nil {
		return
	}

	sortOptions, err := SortAllProducts(c)
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, "invalid request params")
		return
	}

	products, err := h.service.Product.GetAll(context.Background(), filterOptions, sortOptions)
	if err != nil {
		if errors.Is(err, service.ErrProductNotExists) {
			NewErrorResponse(c, http.StatusNotFound, "products not founded")
			return
		}
		NewErrorResponse(c, http.StatusInternalServerError, "internal error")
		return
	}

	c.JSON(http.StatusOK, products)
}

func SortAllProducts(c *gin.Context) (*domain.SortOptions, error) {
	sortOptions := &domain.SortOptions{
		SortBy:    c.DefaultQuery(qSortBy, defaultSortBy),
		SortOrder: c.DefaultQuery(qSortOrder, defaultSortOrder),
	}
	if strings.ToUpper(sortOptions.SortOrder) != "DESC" && strings.ToUpper(sortOptions.SortOrder) != "ASC" {
		return nil, fmt.Errorf("invalid sort order option: %v", sortOptions.SortOrder)
	}

	return sortOptions, nil
}

func FilterAllProducts(c *gin.Context) (filter.Options, error) {
	var limit int
	queryVal, ok := c.GetQuery(qLimit)
	if ok {
		qVal, err := strconv.Atoi(queryVal)
		if err != nil {
			NewErrorResponse(c, http.StatusBadRequest, "invalid request params")
			return nil, fmt.Errorf("invalid request params")
		}
		limit = qVal
	} else {
		NewErrorResponse(c, http.StatusBadRequest, "invalid request params")
		return nil, fmt.Errorf("invalid request params")
	}

	filterOptions := filter.NewOptions(limit)

	name, ok := c.GetQuery(qName)
	if ok && name != "" {
		filterOptions.AddField(qName, filter.OperatorLike, name, filter.DataTypeStr)
	}

	price, ok := c.GetQuery(qPrice)
	if ok {
		operator := filter.OperatorEq
		value := price
		if strings.Contains(price, ":") {
			split := strings.Split(price, ":")
			operator = split[0]
			value = split[1]
		}
		if _, err := strconv.Atoi(value); err != nil {
			NewErrorResponse(c, http.StatusBadRequest, "invalid request params")
			return nil, fmt.Errorf("invalid request params")
		}
		filterOptions.AddField(qPrice, operator, value, filter.DataTypeInt)
	}

	in_stock, ok := c.GetQuery(qInStock)
	if ok {
		_, err := strconv.ParseBool(in_stock)
		if err != nil {
			NewErrorResponse(c, http.StatusBadRequest, "invalid filter params")
			return nil, fmt.Errorf("invalid request params")
		}
		filterOptions.AddField(qInStock, filter.OperatorEq, in_stock, filter.DataTypeBool)
	}

	created_at, ok := c.GetQuery(qCreatedAt)
	if ok && created_at != "" {
		var operator string
		if strings.Contains(created_at, ":") {
			operator = filter.OperatorBetween
		} else {
			operator = filter.OperatorEq
		}
		filterOptions.AddField(qCreatedAt, operator, created_at, filter.DataTypeDate)
	}
	return filterOptions, nil
}
