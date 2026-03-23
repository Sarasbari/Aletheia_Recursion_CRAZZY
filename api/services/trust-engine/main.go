package main

import (
	"log"

	"aletheia-api/shared/config"
	"aletheia-api/shared/utils"
	"github.com/gin-gonic/gin"
)

type trustReq struct {
	SHAMatch       bool `json:"shaMatch"`
	PHashSimilarity bool `json:"pHashSimilarity"`
	Metadata       bool `json:"metadata"`
	Witness        bool `json:"witness"`
}

func main() {
	cfg := config.Load("trust-engine")
	r := gin.New()
	r.Use(gin.Recovery(), gin.Logger())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"service": cfg.ServiceName, "status": "ok"})
	})

	r.POST("/internal/trust/score", func(c *gin.Context) {
		var req trustReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": "invalid payload"})
			return
		}
		score := utils.ComputeTrustScore(req.SHAMatch, req.PHashSimilarity, req.Metadata, req.Witness)
		c.JSON(200, gin.H{"score": score})
	})

	log.Printf("[%s] listening on :%s", cfg.ServiceName, cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
