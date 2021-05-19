package app

import (
	"github.com/gin-gonic/gin"
	"github.com/gvu0110/bookstore_oauth-api/http"
	"github.com/gvu0110/bookstore_oauth-api/repository/db"
	"github.com/gvu0110/bookstore_oauth-api/repository/rest"
	"github.com/gvu0110/bookstore_oauth-api/services/access_token"
)

var (
	router = gin.Default()
)

func StartApplication() {
	atService := access_token.NewService(rest.NewRepository(), db.NewRepository())
	atHandler := http.NewHandler(atService)

	router.GET("/oauth/access_token/:access_token_id", atHandler.GetByID)
	router.POST("/oauth/access_token", atHandler.CreateAccessToken)

	router.Run(":8080")
}
