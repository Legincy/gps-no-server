package config

import (
	"github.com/joho/godotenv"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Mqtt     MqttConfig
}

type ServerConfig struct {
	Host            string
	Port            int
	ReleaseMode     string
	LogLevel        string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ShutdownTimeout time.Duration
}

type DatabaseConfig struct {
	Host     string
	Port     int
	Database string
	User     string
	Password string
	SSLMode  bool
	TimeZone string
}

type MqttConfig struct {
	Host      string
	Port      int
	ClientId  string
	Username  string
	Password  string
	BaseTopic string
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		dir, err := os.Getwd()
		if err != nil {
			rootDir := filepath.Dir(dir)
			envPath := filepath.Join(rootDir, ".env")
			_ = godotenv.Load(envPath)
		}
	}

	config := &Config{
		Server: ServerConfig{
			Host:            getEnv("SERVER_HOST", "localhost"),
			Port:            getEnvAsInt("SERVER_PORT", 8080),
			ReleaseMode:     getEnv("SERVER_RELEASE_MODE", "release"),
			LogLevel:        getEnv("SERVER_LOG_LEVEL", "info"),
			ReadTimeout:     getEnvAsDuration("SERVER_READ_TIMEOUT", 5*time.Second),
			WriteTimeout:    getEnvAsDuration("SERVER_WRITE_TIMEOUT", 5*time.Second),
			ShutdownTimeout: getEnvAsDuration("SERVER_SHUTDOWN_TIMEOUT", 5*time.Second),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnvAsInt("DB_PORT", 5432),
			User:     getEnv("DB_USERNAME", ""),
			Password: getEnv("DB_PASSWORD", ""),
			Database: getEnv("DB_DATABASE", "gps_no"),
			SSLMode:  getEnvAsBool("DB_SSL_MODE", false),
			TimeZone: getEnv("DB_TIME_ZONE", "UTC"),
		},
		Mqtt: MqttConfig{
			Host:      getEnv("MQTT_HOST", "localhost"),
			Port:      getEnvAsInt("MQTT_PORT", 1883),
			ClientId:  getEnv("MQTT_CLIENT_ID", "client"),
			Username:  getEnv("MQTT_USERNAME", ""),
			Password:  getEnv("MQTT_PASSWORD", ""),
			BaseTopic: getEnv("MQTT_BASE_TOPIC", "gps_no"),
		},
	}

	return config, nil
}

//nolint:unused
func getEnvAsStringArray(key string, fallback []string) []string {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return fallback
	}

	values := make([]string, 0)
	for _, v := range strings.Split(valueStr, ",") {
		v = strings.TrimSpace(v)
		if v != "" {
			values = append(values, v)
		}
	}

	return values
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return fallback
}

func getEnvAsInt(key string, fallback int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}

	return fallback
}

func getEnvAsDuration(key string, fallback time.Duration) time.Duration {
	valueStr := getEnv(key, "")
	if value, err := time.ParseDuration(valueStr); err == nil {
		return value
	}

	return fallback
}

func getEnvAsBool(key string, fallback bool) bool {
	valueStr := getEnv(key, "")
	if value, err := strconv.ParseBool(valueStr); err == nil {
		return value
	}

	return fallback
}
