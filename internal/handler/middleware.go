package handler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/DenisEMPS/online-shop/internal/service"
	"github.com/gin-gonic/gin"
)

const (
	authorizationHeader = "Authorization"
)

func (h *Handler) UserIdentity(c *gin.Context) {
	token := c.GetHeader(authorizationHeader)
	if token == "" {
		NewErrorResponse(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	tokenParts := strings.Split(token, " ")
	if len(tokenParts) != 2 || tokenParts[0] != `Bearer` {
		NewErrorResponse(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	id, err := h.service.Auth.ParseToken(tokenParts[1])
	if err != nil {
		if errors.Is(err, service.ErrInvalidToken) {
			NewErrorResponse(c, http.StatusUnauthorized, "unauthorized")
			return
		} else if errors.Is(err, service.ErrTokenInvalidSigningMethod) {
			NewErrorResponse(c, http.StatusUnauthorized, "unauthorized")
			return
		} else {
			NewErrorResponse(c, http.StatusInternalServerError, "internal error")
			return
		}
	}

	c.Set("user_id", id)
}
