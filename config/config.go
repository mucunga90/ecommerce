package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DB    DBConfig
	Redis RedisConfig
	SMS   SMSConfig
	Email EmailConfig
	OIDC  OIDCConfig
}

type DBConfig struct {
	ConnectionStr string
}

type RedisConfig struct {
	Addr     string
	Password string
}

type SMSConfig struct {
	AfricaSTKey  string
	AfricaSTUser string
	SMTPHost     string
	SMTPPort     int
	SMTPUser     string
	SMTPPass     string
}

type EmailConfig struct {
	EmailFrom     string
	EmailPassword string
	AdminEmail    string
	SMTPHost      string
	SMTPPort      int
}

type OIDCConfig struct {
	Issuer   string
	Audience string
	JWKSURL  string
}

func LoadConfig() *Config {
	// Load .env file if exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	smsSmtpPort, err := strconv.Atoi(getEnv("AFRICASTALKING_SMTP_PORT", "587"))
	if err != nil {
		log.Printf("Invalid AFRICASTALKING_SMTP_PORT, using default 587")
		smsSmtpPort = 587
	}

	emailSmtpPort, err := strconv.Atoi(getEnv("EMAIL_SMTP_PORT", "587"))
	if err != nil {
		log.Printf("Invalid EMAIL_SMTP_PORT, using default 587")
		emailSmtpPort = 587
	}

	return &Config{
		DB: DBConfig{
			ConnectionStr: getEnv("DB_CONNECTION_URL", "localhost"),
		},
		Redis: RedisConfig{
			Addr:     getEnv("REDIS_ADDR", "localhost:6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
		},
		Email: EmailConfig{
			EmailFrom:     getEnv("EMAIL_FROM", ""),
			EmailPassword: getEnv("EMAIL_PASSWORD", ""),
			SMTPHost:      getEnv("EMAIL_SMTP_HOST", ""),
			SMTPPort:      emailSmtpPort,
			AdminEmail:    getEnv("ADMIN_EMAIL", ""),
		},
		SMS: SMSConfig{
			AfricaSTKey:  getEnv("AFRICASTALKING_API_KEY", ""),
			AfricaSTUser: getEnv("AFRICASTALKING_USERNAME", "sandbox"),
			SMTPHost:     getEnv("AFRICASTALKING_SMTP_HOST", ""),
			SMTPPort:     smsSmtpPort,
			SMTPUser:     getEnv("AFRICASTALKING_SMTP_USER", ""),
			SMTPPass:     getEnv("AFRICASTALKING_SMTP_PASS", ""),
		},
		OIDC: OIDCConfig{
			Issuer:   getEnv("OIDC_ISSUER", ""),
			Audience: getEnv("OIDC_AUDIENCE", ""),
			JWKSURL:  getEnv("OIDC_JWKS_URL", ""),
		},
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
