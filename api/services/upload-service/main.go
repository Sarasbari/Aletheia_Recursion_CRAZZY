package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"aletheia-api/shared/config"
	"aletheia-api/shared/models"
	"aletheia-api/shared/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func main() {
	httpClient := &http.Client{Timeout: 20 * time.Second}
	defer producer.Close()

	r := gin.New()
	r.Use(gin.Recovery(), gin.Logger())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"service": cfg.ServiceName, "status": "ok"})
	})

	r.POST("/api/v1/images/upload", func(c *gin.Context) {
		file, header, err := c.Request.FormFile("image")
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

		job := models.UploadJob{
			JobID:       uuid.NewString(),
			FileName:    header.Filename,
			ImageBase64: utils.EncodeBase64Image(payload),
			UploadedAt:  time.Now().UTC(),
		}

		body, _ := json.Marshal(job)
		req, err := http.NewRequestWithContext(c.Request.Context(), http.MethodPost, cfg.TaskAPIURL+"/tasks/upload", bytes.NewReader(body))
		if err != nil {
			c.JSON(500, gin.H{"error": "failed to create task request"})
			return
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := httpClient.Do(req)
		if err != nil {
			c.JSON(502, gin.H{"error": "failed to enqueue celery task"})
			return
		}
		defer resp.Body.Close()
		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			c.JSON(502, gin.H{"error": "task api rejected upload job"})
			return
		}

		var queueResp map[string]any
		_ = json.NewDecoder(resp.Body).Decode(&queueResp)

		c.JSON(202, gin.H{
			"jobId":   job.JobID,
			"taskId":  queueResp["taskId"],
			"status":  "QUEUED_IN_CELERY",
			"message": "image accepted for async processing",
		})
	})

	log.Printf("[%s] listening on :%s", cfg.ServiceName, cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
