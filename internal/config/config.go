package config

import (
	"errors"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/go-chi/jwtauth/v5"
	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/spf13/pflag"
)

const (
	defaultServerAddress          = "localhost:8080"
	defaultAccuralURL             = "http://localhost:8080"
	defaultDatabaseDSN            = "user=pqgotest dbname=pqgotest sslmode=verify-full"
	StorageContextTimeout         = time.Second
	TokenSignKey                  = "g1o2p3h4e5r"
	TokenTTL                      = 20 * time.Minute
	RegisteredOrderStatus         = "REGISTERED"
	InvalidOrderStatus            = "INVALID"
	ProcessingOrderStatus         = "PROCESSING"
	ProcessedOrderStatus          = "PROCESSED"
	AccrualOrderStatusURN         = "/api/orders/"
	GetUnprocessedOrdersFrequency = time.Second
	RetryAfterErrorDefaultTimeout = 30 * time.Second
)

var TokenAuth *jwtauth.JWTAuth = jwtauth.New(string(jwa.HS256), []byte(TokenSignKey), nil)

type Config struct {
	ServerAddress string `env:"RUN_ADDRESS" envDefault:"localhost:8080"`
	AccrualURL    string `env:"ACCRUAL_SYSTEM_ADDRESS"`
	DatabaseDSN   string `env:"DATABASE_URI"`
}

func New() *Config {
	var config Config
	err := env.Parse(&config)
	if err != nil {
		log.Fatal(err)
	}
	flagSet := pflag.FlagSet{}
	addrFlag := flagSet.StringP("-addr", "a", config.ServerAddress, "Server address: host:port")
	accrualFlag := flagSet.StringP("-acc", "r", config.AccrualURL, "Accrual service URL")
	dbDSNFlag := flagSet.StringP("-dbDsn", "d", config.DatabaseDSN, "Database DSN string")

	err = flagSet.Parse(os.Args[1:])
	if err != nil {
		log.Fatal("Error while parsing sys Args")
	}
	config.ServerAddress = *addrFlag
	config.AccrualURL = *accrualFlag
	config.DatabaseDSN = *dbDSNFlag

	err = validateConfig(&config)
	if err != nil {
		log.Fatal(err)
	}

	return &config
}

//parse url and return nil if url is valid or error
func validateURL(s string) error {
	_, err := url.ParseRequestURI(s)
	if err != nil {
		return err
	}

	u, err := url.Parse(s)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return err
	}
	return nil
}

func validateConfig(c *Config) error {
	addrErr := validateURL(c.ServerAddress)
	accuralURLerr := validateURL(c.AccrualURL)
	if addrErr != nil {
		return errors.New("wrong server address param")
	}
	if accuralURLerr != nil {
		return errors.New("wrong accrual url param")
	}
	return nil
}
