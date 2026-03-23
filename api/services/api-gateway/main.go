package main

import (
	"log"
	"net/http/httputil"
	"net/url"

	"aletheia-api/shared/config"
	"aletheia-api/shared/middleware"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load("api-gateway")

	r := gin.New()
	r.Use(gin.Recovery(), gin.Logger())
	r.Use(middleware.APIKey(cfg.APIKey))
	r.Use(middleware.RateLimit(120))

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"service": cfg.ServiceName, "status": "ok"})
	})

	uploadProxy := mustProxy(cfg.UploadServiceURL)
	verifyProxy := mustProxy(cfg.VerifyServiceURL)

	r.Any("/api/v1/images/upload", gin.WrapH(uploadProxy))
	r.Any("/api/v1/verify", gin.WrapH(verifyProxy))

	log.Printf("[%s] listening on :%s", cfg.ServiceName, cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}

func mustProxy(rawURL string) *httputil.ReverseProxy {
	u, err := url.Parse(rawURL)
	if err != nil {
		panic(err)
	}
	return httputil.NewSingleHostReverseProxy(u)
}
