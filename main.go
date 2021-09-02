package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"go.simple/structure/handlers"
	"go.simple/structure/middlewares"
	"go.simple/structure/resources"
)

func main() {
	// Load environment config
	cfg := newConfig()

	// Connect to mysql server
	db, err := resources.NewConnection(cfg.mysqlDataSourceName)
	if err != nil {
		log.Fatal("Connect to mysql error: ", err.Error())
	}

	ph := handlers.NewPaymentHandler(db)

	r := gin.Default()
	r.Use(middlewares.Auth)

	r.POST("/api/users/:user_id/transactions", ph.CreateTransaction)
	r.GET("/api/users/:user_id/transactions", ph.GetUserTransaction)
	r.PUT("/api/users/:user_id/transactions", ph.UpdateTransaction)
	r.DELETE("/api/users/:user_id/transactions", ph.DeleteTransaction)

	_ = r.Run(":" + cfg.ginPort)
}
