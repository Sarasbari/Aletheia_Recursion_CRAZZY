package config

import (
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	ServiceName string
	Port        string

	APIKey           string
	UploadServiceURL string
	VerifyServiceURL string

	TaskAPIURL       string
	CeleryBrokerURL  string
	CeleryBackendURL string

	RedisAddr string
	RedisDB   int

	IPFSEndpoint string

	BigchainDBURL    string
	BigchainDBAPIKey string

	PolygonRPCURL          string
	PolygonPrivateKey      string
	PolygonContractAddress string
	PolygonChainID         string

	CassandraHosts []string
	CassandraKeyspace string
}

func Load(serviceName string) Config {
	_ = godotenv.Load()

	return Config{
		ServiceName: value("SERVICE_NAME", serviceName),
		Port:        value("PORT", defaultPort(serviceName)),

		APIKey:           os.Getenv("API_KEY"),
		UploadServiceURL: value("UPLOAD_SERVICE_URL", "http://localhost:8081"),
		VerifyServiceURL: value("VERIFY_SERVICE_URL", "http://localhost:8082"),
		TaskAPIURL:       value("TASK_API_URL", "http://localhost:8090"),
		CeleryBrokerURL:  value("CELERY_BROKER_URL", "redis://localhost:6379/1"),
		CeleryBackendURL: value("CELERY_BACKEND_URL", "redis://localhost:6379/2"),

		RedisAddr: value("REDIS_ADDR", "localhost:6379"),
		RedisDB:   intValue("REDIS_DB", 0),

		IPFSEndpoint: value("IPFS_ENDPOINT", ""),

		BigchainDBURL:    value("BIGCHAINDB_URL", ""),
		BigchainDBAPIKey: value("BIGCHAINDB_API_KEY", ""),

		PolygonRPCURL:          value("POLYGON_RPC_URL", ""),
		PolygonPrivateKey:      value("POLYGON_PRIVATE_KEY", ""),
		PolygonContractAddress: value("POLYGON_CONTRACT_ADDRESS", ""),
		PolygonChainID:         value("POLYGON_CHAIN_ID", "137"),

		CassandraHosts:   splitCSV(value("CASSANDRA_HOSTS", "localhost")),
		CassandraKeyspace: value("CASSANDRA_KEYSPACE", "aletheia"),
	}
}

func defaultPort(serviceName string) string {
	switch serviceName {
	case "api-gateway":
		return "8080"
	case "upload-service":
		return "8081"
	case "verification-service":
		return "8082"
	case "hash-service":
		return "8083"
	case "merkle-service":
		return "8084"
	case "trust-engine":
		return "8085"
	case "storage-service":
		return "8086"
	case "blockchain-service":
		return "8087"
	default:
		return "8080"
	}
}

func value(key, fallback string) string {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return fallback
	}
	return v
}

func intValue(key string, fallback int) int {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(v)
	if err != nil {
		return fallback
	}
	return parsed
}

func splitCSV(v string) []string {
	parts := strings.Split(v, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		trimmed := strings.TrimSpace(p)
		if trimmed != "" {
			out = append(out, trimmed)
		}
	}
	if len(out) == 0 {
		return []string{"localhost:9092"}
	}
	return out
}
