package config

import "os"

type Config struct {
	Port           string
	AllowedOrigins string
	AdminUser      string
	AdminPass      string
	DBType         string // "memory" or "oracle"
	OracleDBDSN    string
	OCIBucketName  string
	OCINamespace   string
	OCIRegion      string
}

func Load() *Config {
	return &Config{
		Port:           getEnv("PORT", "8080"),
		AllowedOrigins: getEnv("ALLOWED_ORIGINS", "http://localhost:5173"),
		AdminUser:      getEnv("ADMIN_USER", "admin"),
		AdminPass:      getEnv("ADMIN_PASS", ""),
		DBType:         getEnv("DB_TYPE", "memory"),
		OracleDBDSN:    getEnv("ORACLE_DB_DSN", ""),
		OCIBucketName:  getEnv("OCI_BUCKET_NAME", "tech-blog-images"),
		OCINamespace:   getEnv("OCI_BUCKET_NAMESPACE", "axva0xxfvkwr"),
		OCIRegion:      getEnv("OCI_REGION", "mx-queretaro-1"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
