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
	Server   ServerConfig   `json:"server"`
	Database DatabaseConfig `json:"database"`
	Mqtt     MqttConfig     `json:"mqtt"`
}

type ServerConfig struct {
	Host            string        `json:"host"`
	Port            int           `json:"port"`
	ReleaseMode     string        `json:"release_mode"`
	LogLevel        string        `json:"log_level"`
	ReadTimeout     time.Duration `json:"read_timeout"`
	WriteTimeout    time.Duration `json:"write_timeout"`
	ShutdownTimeout time.Duration `json:"shutdown_timeout"`
}

type DatabaseConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Database string `json:"database"`
	User     string `json:"user"`
	Password string `json:"password"`
	SSLMode  bool   `json:"ssl_mode"`
	TimeZone string `json:"time_zone"`
}

type MqttConfig struct {
	Host                 string        `json:"host"`
	Port                 int           `json:"port"`
	ClientId             string        `json:"client_id"`
	Username             string        `json:"username"`
	Password             string        `json:"password"`
	BaseTopic            string        `json:"base_topic"`
	AutoReconnect        bool          `json:"auto_reconnect"`
	MaxReconnectInterval time.Duration `json:"max_reconnect_interval"`
	CleanSession         bool          `json:"clean_session"`
}

func LoadEnvFile() {
	if err := godotenv.Load(); err != nil {
		dir, err := os.Getwd()
		if err != nil {
			return
		}

		rootDir := filepath.Dir(dir)
		envPath := filepath.Join(rootDir, ".env")
		_ = godotenv.Load(envPath)
	}
}

func Load() (*Config, error) {
	LoadEnvFile()

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
			Host:                 getEnv("MQTT_HOST", "localhost"),
			Port:                 getEnvAsInt("MQTT_PORT", 1883),
			ClientId:             getEnv("MQTT_CLIENT_ID", "client"),
			Username:             getEnv("MQTT_USERNAME", ""),
			Password:             getEnv("MQTT_PASSWORD", ""),
			BaseTopic:            getEnv("MQTT_BASE_TOPIC", "gps_no"),
			AutoReconnect:        getEnvAsBool("MQTT_AUTO_RECONNECT", true),
			MaxReconnectInterval: getEnvAsDuration("MQTT_MAX_RECONNECT", 1*time.Second),
			CleanSession:         getEnvAsBool("MQTT_CLEAN_SESSION", true),
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
