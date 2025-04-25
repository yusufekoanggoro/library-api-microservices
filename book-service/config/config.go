package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type ConfigProvider interface {
	GetHTTPHost() string
	GetHTTPPort() string
	GetGRPCHost() string
	GetGRPCPort() string

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

	GetRedisHost() string
	GetRedisPort() string
	GetRedisPassowrd() string
	GetRedisDB() string

	GetBasicAuthUsername() string
	GetBasicAuthPassword() string
}

type EnvConfig struct {
	HTTPHost string
	HTTPPort string
	GRPCHost string
	GRPCPort string

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

	RedisHost     string
	RedisPort     string
	RedisPassword string
	RedisDB       string

	BasicAuthUsername string
	BasicAuthPassword string
}

func (e *EnvConfig) GetHTTPHost() string { return e.HTTPHost }
func (e *EnvConfig) GetHTTPPort() string { return e.HTTPPort }
func (e *EnvConfig) GetGRPCHost() string { return e.GRPCHost }
func (e *EnvConfig) GetGRPCPort() string { return e.GRPCPort }

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

func (e *EnvConfig) GetRedisHost() string     { return e.RedisHost }
func (e *EnvConfig) GetRedisPort() string     { return e.RedisPort }
func (e *EnvConfig) GetRedisPassowrd() string { return e.RedisPassword }
func (e *EnvConfig) GetRedisDB() string       { return e.RedisDB }

func (e *EnvConfig) GetBasicAuthUsername() string { return e.BasicAuthUsername }
func (e *EnvConfig) GetBasicAuthPassword() string { return e.BasicAuthPassword }

func LoadConfig() ConfigProvider {
	err := godotenv.Load()
	if err != nil {
		log.Println("Gagal membaca file .env, menggunakan environment variables yang tersedia")
	}

	return &EnvConfig{
		HTTPHost: os.Getenv("HTTP_HOST"),
		HTTPPort: os.Getenv("HTTP_PORT"),

		GRPCHost: os.Getenv("GRPC_HOST"),
		GRPCPort: os.Getenv("GRPC_PORT"),

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

		RedisHost:     os.Getenv("REDIS_HOST"),
		RedisPort:     os.Getenv("REDIS_PORT"),
		RedisPassword: os.Getenv("REDIS_PASSWORD"),
		RedisDB:       os.Getenv("REDIS_DB"),

		BasicAuthUsername: os.Getenv("BASIC_AUTH_USER"),
		BasicAuthPassword: os.Getenv("BASIC_AUTH_PASS"),
	}
}
