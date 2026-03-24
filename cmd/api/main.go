package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/user/aletheia-api/internal/blockchain"
	"github.com/user/aletheia-api/internal/storage"
	"github.com/user/aletheia-api/internal/trust"
	"github.com/user/aletheia-api/internal/verify"
	"github.com/user/aletheia-api/routes"
)

// @title           Aletheia API
// @version         1.0
// @description     Decentralized Image Authenticity System API.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

func main() {
	// Configuration (using environment variables or defaults)
	ipfsURL := getEnv("IPFS_URL", "localhost:5001")
	bigchainURL := getEnv("BIGCHAINDB_URL", "https://test.ipdb.io")            // Public BigchainDB testnet (IPDB)
	polygonURL := getEnv("POLYGON_URL", "https://rpc-amoy.polygon.technology") // Polygon Amoy Testnet
	privateKey := getEnv("PRIVATE_KEY", "7a89f52d8bb497c50230c9043b55c3e4afe2b5de656e049f155701b5408c4101")

	// Initialize services
	ipfsStorage := storage.NewIPFSStorage(ipfsURL)
	bigchainStorage, err := storage.NewBigchainDBStorage(bigchainURL, getEnv("BIGCHAINDB_SEED", "3f4103f980dabf7ad6fc0fea32f521bf9059aa47d4b51247af673ba1b080c62b"))
	if err != nil {
		log.Fatalf("Failed to initialize BigchainDB storage: %v", err)
	}

	contractAddr := getEnv("CONTRACT_ADDRESS", "0xcb9b3e944c6883f3f5a2a19cca0e608060824b23d9b0f4c3c3c3c3c3c3c3c3c3")
	polyService, err := blockchain.NewPolygonService(polygonURL, privateKey, contractAddr)
	if err != nil {
		log.Printf("Warning: Failed to initialize Polygon service: %v", err)
	}

	trustEngine := trust.NewScoringEngine()
	verifyService := verify.NewService(ipfsStorage, bigchainStorage, polyService, trustEngine)

	// Initialize Gin
	r := gin.Default()

	// Register routes
	handler := routes.NewHandler(verifyService)
	routes.RegisterRoutes(r, handler)

	// Start server
	port := getEnv("PORT", "8080")
	log.Printf("Aletheia API starting on port %s...", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
