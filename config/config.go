package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

// LoadConfig :nodoc:
func LoadConfig() {
	err := godotenv.Load()
	if err != nil && Env() != "test" {
		logrus.Warningf("%v", err)
	}
}

// Env :nodoc:
func Env() string {
	return os.Getenv("ENV")
}

// HTTPPort :nodoc:
func HTTPPort() string {
	return os.Getenv("HTTP_PORT")
}

// MysqlHost :nodoc:
func MysqlHost() string {
	return os.Getenv("MYSQL_HOST")
}

// MysqlUser :nodoc:
func MysqlUser() string {
	return os.Getenv("MYSQL_USER")
}

// MysqlPassword :nodoc:
func MysqlPassword() string {
	return os.Getenv("MYSQL_PASSWORD")
}

// MysqlDB :nodoc:
func MysqlDB() string {
	return os.Getenv("MYSQL_DB")
}

// MysqlPort :nodoc:
func MysqlPort() int {
	port, err := strconv.Atoi(os.Getenv("MYSQL_PORT"))
	if err != nil {
		return 0
	}

	return port
}

// MysqlDSN :nodoc:
func MysqlDSN() string {
	return fmt.Sprintf(
		"mysql://%s:%s@%s:%d/%s?charset=utf8&parseTime=True&loc=Local",
		MysqlUser(),
		MysqlPassword(),
		MysqlHost(),
		MysqlPort(),
		MysqlDB(),
	)
}

// LogLevel :nodoc:
func LogLevel() string {
	return os.Getenv("LOG_LEVEL")
}
