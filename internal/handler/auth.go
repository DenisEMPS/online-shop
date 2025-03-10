package handler

import (
	"errors"
	"net/http"

	"github.com/DenisEMPS/online-shop/internal/domain"
	"github.com/DenisEMPS/online-shop/internal/service"
	"github.com/gin-gonic/gin"
)

func (h *Handler) RegisterNewUser(c *gin.Context) {
	var input domain.UserCreate

	if err := c.ShouldBindJSON(&input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, "invalid request")
		return
	}

	id, err := h.service.Auth.Register(&input)
	if err != nil {
		if errors.Is(err, service.ErrUserExists) {
			NewErrorResponse(c, http.StatusConflict, "user already exists")
			return
		}
		NewErrorResponse(c, http.StatusInternalServerError, "internal error")
		return
	}

	c.JSON(http.StatusCreated, map[string]int64{
		"id": id,
	})
}

func (h *Handler) LoginUser(c *gin.Context) {
	var input domain.UserLogin

	if err := c.ShouldBindJSON(&input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, "invalid request")
		return
	}

	token, err := h.service.Auth.Login(&input)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			NewErrorResponse(c, http.StatusUnauthorized, "user not found")
			return
		} else if errors.Is(err, service.ErrInvalidCredentials) {
			NewErrorResponse(c, http.StatusUnauthorized, "wrong password")
			return
		} else if errors.Is(err, service.ErrInvalidToken) {
			NewErrorResponse(c, http.StatusUnauthorized, "invalid session token")
			return
		} else if errors.Is(err, service.ErrTokenInvalidSigningMethod) {
			NewErrorResponse(c, http.StatusUnauthorized, "invalid session token")
			return
		} else {
			NewErrorResponse(c, http.StatusInternalServerError, "internal error")
			return
		}
	}
	c.JSON(http.StatusOK, map[string]string{
		"token": token,
	})
}
