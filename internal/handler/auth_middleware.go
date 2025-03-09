package handler

import (
	"errors"
	"fmt"
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

func getUserID(c *gin.Context) (int64, error) {
	idParam, ok := c.Get("user_id")
	if !ok {
		NewErrorResponse(c, http.StatusUnauthorized, "invalid request params")
		return 0, fmt.Errorf("failed to get user id from context")
	}

	id, ok := idParam.(int64)
	if !ok {
		NewErrorResponse(c, http.StatusUnauthorized, "invalid request params")
		return 0, fmt.Errorf("failed to cast id in int %v", id)
	}

	return id, nil
}
