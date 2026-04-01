package database

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"time"
)

// APIKey special string type to validate api keys
type APIKey string

// ServerConfig server configuration
type ServerConfig struct {
	Host       string
	Port       int
	SwaggerDir string
}

// ServerConfig server configuration
type MailerConfig struct {
	Address            string
	Secret             string
	Host               string
	Port               int
	LoginLink          string
	ForgotPasswordLink string
}

// DatabaseConfig database configuration
type DatabaseConfig struct {
	User     string
	Password string
	Name     string
	Host     string
	Port     int
}

// APIKeys all api keys
type APIKeys struct {
	Mandrill APIKey
}

// Config general backend configuration
type Config struct {
	DatabaseConfig DatabaseConfig
	ServerConfig   ServerConfig
	Keys           APIKeys
	Mailer         MailerConfig
	AWSConfig      AWSConfig
	RedisConfig    RedisConfig
}

// RedisConfig related configuration
type RedisConfig struct {
	Host string
	Port string
}

// AWSConfig AWS related configuration
type AWSConfig struct {
	Region               string
	Bucket               string
	URL                  string
	PresignedURLDuration time.Duration
	ADMIN_URL            string
}

func initializeStruct(valueType reflect.Type, valuePtr reflect.Value, rootKey string) error {
	var value reflect.Value
	var err error
	value = valuePtr.Elem()
	for i := 0; i < value.NumField(); i++ {
		var field reflect.Value
		var fieldType reflect.StructField
		var key string
		var integer int
		field = value.Field(i)
		fieldType = valueType.Field(i)
		key = fmt.Sprintf("%s_%s", rootKey, fieldType.Name)
		switch fieldType.Type.Kind() {
		case reflect.String:
			field.SetString(os.Getenv(key))
		case reflect.Int:
			integer, err = strconv.Atoi(os.Getenv(key))
			if err != nil {
				continue
			}
			field.SetInt(int64(integer))
		case reflect.Struct:
			// We need this to be writeable
			valuePtr = reflect.New(fieldType.Type)
			// Call recursively now
			err = initializeStruct(fieldType.Type, valuePtr, fieldType.Name)
			if err != nil {
				return err
			}
			field.Set(valuePtr.Elem())
		case reflect.Int64:
		default:
			return fmt.Errorf("cannot determine the type of `%s'", fieldType.Name)
		}
	}
	return nil
}

// initialize config
func InitializeConfig() (*Config, error) {
	var result Config
	var err error
	err = initializeStruct(reflect.TypeOf(result), reflect.ValueOf(&result), "")
	if err != nil {
		return nil, err
	}
	return &result, nil
}
