package http

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gvu0110/bookstore_oauth-api/domain/access_token"
	"github.com/gvu0110/bookstore_oauth-api/utils/errors"
)

type AccessTokenHandler interface {
	GetByID(*gin.Context)
	CreateAccessToken(*gin.Context)
}

type accessTokenHandler struct {
	service access_token.Service
}

func NewHandler(service access_token.Service) AccessTokenHandler {
	return &accessTokenHandler{
		service: service,
	}
}

func (handler *accessTokenHandler) GetByID(c *gin.Context) {
	accessToken, err := handler.service.GetByID(strings.TrimSpace(c.Param("access_token_id")))
	if err != nil {
		c.JSON(err.StatusCode, err)
		return
	}
	c.JSON(http.StatusOK, accessToken)
}

func (handler *accessTokenHandler) CreateAccessToken(c *gin.Context) {
	var at access_token.AccessToken
	if err := c.ShouldBindJSON(&at); err != nil {
		c.JSON(errors.NewBadRequestRESTError("Invalid JSON body").StatusCode, errors.NewBadRequestRESTError("Invalid JSON body"))
		return
	}
	if err := handler.service.CreateAccessToken(at); err != nil {
		c.JSON(err.StatusCode, err)
		return
	}
	c.JSON(http.StatusCreated, at)
}
