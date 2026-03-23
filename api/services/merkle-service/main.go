package main

import (
	"log"

	"aletheia-api/shared/config"
	"aletheia-api/shared/utils"
	"github.com/gin-gonic/gin"
)

type merkleReq struct {
	ImageBase64 string `json:"imageBase64"`
}

func main() {
	cfg := config.Load("merkle-service")
	r := gin.New()
	r.Use(gin.Recovery(), gin.Logger())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"service": cfg.ServiceName, "status": "ok"})
	})

	r.POST("/internal/merkle", func(c *gin.Context) {
		var req merkleReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": "invalid payload"})
			return
		}
		img, err := utils.DecodeBase64Image(req.ImageBase64)
		if err != nil {
			c.JSON(400, gin.H{"error": "invalid image"})
			return
		}
		root, err := utils.ComputeMerkleRoot(img)
		if err != nil {
			c.JSON(400, gin.H{"error": "merkle generation failed"})
			return
		}
		c.JSON(200, gin.H{"merkleRoot": root})
	})

	log.Printf("[%s] listening on :%s", cfg.ServiceName, cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
