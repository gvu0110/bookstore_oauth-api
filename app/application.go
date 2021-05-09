package app

import (
	"github.com/gin-gonic/gin"
	"../http"
)
var (
	router = gin.Default()
)

func StartApplication() {
	atHandler := http.NewAccess