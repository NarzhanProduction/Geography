package config

import "os"

type Config struct {
	Address    string
	Port       string
	JWTKey     string
	JWTttl     string
	OpenAPIKey string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBMaxConns int
}

func InitConfig() Config {
	return Config{
		Address:    os.Getenv("ADDRESS"),
		Port:       os.Getenv("PORT"),
		JWTKey:     os.Getenv("JWT_KEY"),
		JWTttl:     os.Getenv("JWT_TTL"),
		OpenAPIKey: os.Getenv("OPENAI_API_KEY"),
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
		DBMaxConns: 10,
	}
}

func (c *Config) DSN() string {
	return "postgres://" + c.DBUser + ":" + c.DBPassword + "@" + c.DBHost + ":" + c.DBPort + "/" + c.DBName + "?sslmode=disable"
}
