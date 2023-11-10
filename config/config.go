package config

import (
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// LoadConfig :nodoc:
func LoadConfig() {
	viper.AddConfigPath(".")
	viper.SetConfigName("config")

	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)

	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil && Env() != "test" {
		logrus.Warningf("%v", err)
	}
}

// Env :nodoc:
func Env() string {
	return viper.GetString("env")
}

// HTTPPort :nodoc:
func HTTPPort() string {
	return viper.GetString("http_port")
}

// MysqlHost :nodoc:
func MysqlHost() string {
	return viper.GetString("mysql.host")
}

// MysqlUser :nodoc:
func MysqlUser() string {
	return viper.GetString("mysql.user")
}

// MysqlPassword :nodoc:
func MysqlPassword() string {
	return viper.GetString("mysql.password")
}

// MysqlDB :nodoc:
func MysqlDB() string {
	return viper.GetString("mysql.database")
}

// MysqlPort :nodoc:
func MysqlPort() int {
	return viper.GetInt("mysql.port")
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
	return viper.GetString("log_level")
}

// JwtSecret :nodoc:
func JwtSecret() string {
	return viper.GetString("jwt_secret")
}

// RedisHost :nodoc:
func RedisHost() string {
	return viper.GetString("redis_host")
}

// RedisDB :nodoc:
func RedisDB() int {
	return viper.GetInt("redis_db")
}
