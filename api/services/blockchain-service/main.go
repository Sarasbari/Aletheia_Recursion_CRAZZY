package main

import (
	"log"

	"aletheia-api/shared/config"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load("blockchain-service")
	r := gin.New()
	r.Use(gin.Recovery(), gin.Logger())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"service": cfg.ServiceName,
			"status":  "ok",
			"note":    "anchoring is primarily handled asynchronously by storage-worker",
		})
	})

	log.Printf("[%s] listening on :%s", cfg.ServiceName, cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
