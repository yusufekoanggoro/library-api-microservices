package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type ConfigProvider interface {
	GetHTTPPort() string
	GetHTTPHost() string

	GetBookGRPCHost() string
	GetBookGRPCPort() string

	GetDBReadHost() string
	GetDBReadPort() string
	GetDBReadUser() string
	GetDBReadPassword() string
	GetDBReadName() string

	GetDBWriteHost() string
	GetDBWritePort() string
	GetDBWriteUser() string
	GetDBWritePassword() string
	GetDBWriteName() string

	GetDBSSLMode() string

	GetBasicAuthUsername() string
	GetBasicAuthPassword() string
}

type EnvConfig struct {
	DBReadHost     string
	DBReadPort     string
	DBReadUser     string
	DBReadPassword string
	DBReadName     string

	DBWriteHost     string
	DBWritePort     string
	DBWriteUser     string
	DBWritePassword string
	DBWriteName     string

	DBSSLMode string

	HTTPHost string
	HTTPPort string

	BookGRPCHost string
	BookGRPCPort string

	BasicAuthUsername string
	BasicAuthPassword string
}

func (e *EnvConfig) GetHTTPHost() string { return e.HTTPHost }
func (e *EnvConfig) GetHTTPPort() string { return e.HTTPPort }

func (e *EnvConfig) GetBookGRPCHost() string { return e.BookGRPCHost }
func (e *EnvConfig) GetBookGRPCPort() string { return e.BookGRPCPort }

func (e *EnvConfig) GetDBReadHost() string     { return e.DBReadHost }
func (e *EnvConfig) GetDBReadPort() string     { return e.DBReadPort }
func (e *EnvConfig) GetDBReadUser() string     { return e.DBReadUser }
func (e *EnvConfig) GetDBReadPassword() string { return e.DBReadPassword }
func (e *EnvConfig) GetDBReadName() string     { return e.DBReadName }

func (e *EnvConfig) GetDBWriteHost() string     { return e.DBWriteHost }
func (e *EnvConfig) GetDBWritePort() string     { return e.DBWritePort }
func (e *EnvConfig) GetDBWriteUser() string     { return e.DBWriteUser }
func (e *EnvConfig) GetDBWritePassword() string { return e.DBWritePassword }
func (e *EnvConfig) GetDBWriteName() string     { return e.DBWriteName }

func (e *EnvConfig) GetDBSSLMode() string { return e.DBSSLMode }

func (e *EnvConfig) GetBasicAuthUsername() string { return e.BasicAuthUsername }
func (e *EnvConfig) GetBasicAuthPassword() string { return e.BasicAuthPassword }

func LoadConfig() ConfigProvider {
	if err := godotenv.Load(); err != nil {
		log.Println("Failed to read .env file, using available environment variables.")
	} else {
		log.Println("Successfully loaded .env file.")
	}

	return &EnvConfig{
		HTTPHost: os.Getenv("HTTP_HOST"),
		HTTPPort: os.Getenv("HTTP_PORT"),

		BookGRPCHost: os.Getenv("BOOK_GRPC_HOST"),
		BookGRPCPort: os.Getenv("BOOK_GRPC_PORT"),

		DBReadHost:     os.Getenv("DB_READ_HOST"),
		DBReadPort:     os.Getenv("DB_READ_PORT"),
		DBReadUser:     os.Getenv("DB_READ_USER"),
		DBReadPassword: os.Getenv("DB_READ_PASSWORD"),
		DBReadName:     os.Getenv("DB_READ_NAME"),

		DBWriteHost:     os.Getenv("DB_WRITE_HOST"),
		DBWritePort:     os.Getenv("DB_WRITE_PORT"),
		DBWriteUser:     os.Getenv("DB_WRITE_USER"),
		DBWritePassword: os.Getenv("DB_WRITE_PASSWORD"),
		DBWriteName:     os.Getenv("DB_WRITE_NAME"),

		DBSSLMode: os.Getenv("DB_SSLMODE"),

		BasicAuthUsername: os.Getenv("BASIC_AUTH_USER"),
		BasicAuthPassword: os.Getenv("BASIC_AUTH_PASS"),
	}
}
