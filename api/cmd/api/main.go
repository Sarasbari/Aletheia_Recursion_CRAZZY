package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"aletheia-api/internal/blockchain"
	"aletheia-api/internal/storage"
	"aletheia-api/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	logger := log.New(os.Stdout, "[aletheia-api] ", log.LstdFlags|log.Lshortfile)

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery(), gin.Logger())

	repo := storage.NewInMemoryRepository()
	ipfs := storage.NewIPFSService(os.Getenv("IPFS_ENDPOINT"))
	bigchain := storage.NewBigchainHTTPClient(os.Getenv("BIGCHAINDB_URL"), os.Getenv("BIGCHAINDB_API_KEY"))
	polygon := blockchain.NewPolygonClient(
		os.Getenv("POLYGON_RPC_URL"),
		os.Getenv("POLYGON_PRIVATE_KEY"),
		os.Getenv("POLYGON_CONTRACT_ADDRESS"),
		os.Getenv("POLYGON_CHAIN_ID"),
	)

	routes.Register(r, repo, ipfs, bigchain, polygon, logger)

	srv := &http.Server{
		Addr:              ":8080",
		Handler:           r,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	go func() {
		logger.Printf("server starting on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("server failed: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	logger.Println("shutting down server")
	if err := srv.Shutdown(ctx); err != nil {
		logger.Printf("graceful shutdown failed: %v", err)
	}
}
