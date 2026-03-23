package main

import (
	"io"
	"log"
	"time"

	"aletheia-api/shared/cache"
	"aletheia-api/shared/clients"
	"aletheia-api/shared/config"
	"aletheia-api/shared/models"
	"aletheia-api/shared/utils"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load("verification-service")
	cacheClient := cache.New(cfg.RedisAddr, cfg.RedisDB)
	bigchain := clients.NewBigchainClient(cfg.BigchainDBURL, cfg.BigchainDBAPIKey)

	r := gin.New()
	r.Use(gin.Recovery(), gin.Logger())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"service": cfg.ServiceName, "status": "ok"})
	})

	r.POST("/api/v1/verify", func(c *gin.Context) {
		file, _, err := c.Request.FormFile("image")
		if err != nil {
			c.JSON(400, gin.H{"error": "image is required"})
			return
		}
		defer file.Close()

		payload, err := io.ReadAll(io.LimitReader(file, 25<<20))
		if err != nil || len(payload) == 0 {
			c.JSON(400, gin.H{"error": "invalid image"})
			return
		}

		sha, ph, err := utils.ComputeSHA256AndPHash(payload)
		if err != nil {
			c.JSON(400, gin.H{"error": "failed to hash image"})
			return
		}

		cacheKey := "verify:" + sha
		var cached models.VerifyResponse
		hit, err := cacheClient.GetJSON(c.Request.Context(), cacheKey, &cached)
		if err == nil && hit {
			c.JSON(200, cached)
			return
		}

		record, err := bigchain.FindBySHA(c.Request.Context(), sha)
		if err != nil {
			c.JSON(502, gin.H{"error": "metadata source unavailable"})
			return
		}
		if record == nil {
			result := models.VerifyResponse{TrustScore: 0, Verdict: "UNKNOWN", MatchedProofID: "", SimilarityScore: 0}
			_ = cacheClient.SetJSON(c.Request.Context(), cacheKey, result, 10*time.Minute)
			c.JSON(200, result)
			return
		}

		similarity, err := utils.SimilarityScore(ph, record.PHash)
		if err != nil {
			similarity = 0
		}
		computedMerkle, err := utils.ComputeMerkleRoot(payload)
		if err != nil {
			computedMerkle = ""
		}

		shaMatch := sha == record.SHA256
		pHashMatch := similarity >= 80
		metaMatch := computedMerkle != "" && computedMerkle == record.MerkleRoot

		score := utils.ComputeTrustScore(shaMatch, pHashMatch, metaMatch, false)
		verdict := "UNKNOWN"
		switch {
		case shaMatch && metaMatch:
			verdict = "AUTHENTIC"
		case !shaMatch && pHashMatch:
			verdict = "SIMILAR"
		case shaMatch && !metaMatch:
			verdict = "TAMPERED"
		default:
			verdict = "TAMPERED"
		}

		result := models.VerifyResponse{
			TrustScore:      score,
			Verdict:         verdict,
			MatchedProofID:  record.JobID,
			SimilarityScore: similarity,
		}
		_ = cacheClient.SetJSON(c.Request.Context(), cacheKey, result, 10*time.Minute)
		c.JSON(200, result)
	})

	log.Printf("[%s] listening on :%s", cfg.ServiceName, cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
