package server

import (
	"verve_project/controller"

	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {

	engine := gin.New()

	app := engine.Group("/api")

	go controller.LogUniqueRequestsEveryMinute()

	app.GET("/verve/accept", controller.HandleVerveAccept)

	return engine
}
