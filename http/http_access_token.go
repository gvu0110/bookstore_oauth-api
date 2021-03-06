package http

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	atDomain "github.com/gvu0110/bookstore_oauth-api/domain/access_token"
	"github.com/gvu0110/bookstore_oauth-api/services/access_token"
	"github.com/gvu0110/bookstore_utils-go/rest_errors"
)

type AccessTokenHandler interface {
	GetByID(*gin.Context)
	CreateAccessToken(*gin.Context)
	DeleteAccessToken(*gin.Context)
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
		c.String(err.StatusCode(), err.Message())
		return
	}
	c.JSON(http.StatusOK, accessToken)
}

func (handler *accessTokenHandler) CreateAccessToken(c *gin.Context) {
	var request atDomain.AccessTokenRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(rest_errors.NewBadRequestRESTError("Invalid JSON body").StatusCode(), rest_errors.NewBadRequestRESTError("Invalid JSON body"))
		return
	}

	at, err := handler.service.CreateAccessToken(request)
	if err != nil {
		c.JSON(err.StatusCode(), err)
		return
	}
	c.JSON(http.StatusCreated, at)
}

func (handler *accessTokenHandler) DeleteAccessToken(c *gin.Context) {
	accessToken := strings.TrimSpace(c.Param("access_token_id"))
	var request atDomain.AccessTokenRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(rest_errors.NewBadRequestRESTError("Invalid JSON body").StatusCode(), rest_errors.NewBadRequestRESTError("Invalid JSON body"))
		return
	}

	if err := handler.service.DeleteAccessToken(accessToken, request); err != nil {
		c.String(err.StatusCode(), err.Message())
		return
	}
	responseBody := `"status":"deleted"`
	c.JSON(http.StatusOK, responseBody)
}
