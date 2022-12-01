package config

import (
	"errors"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/spf13/pflag"
)

const (
	defaultServerAddress = "localhost:8080"
	defaultAccuralURL    = ""
	defaultDatabaseDSN   = "user=pqgotest dbname=pqgotest sslmode=verify-full"
)

const StorageContextTimeout = time.Second

type Config struct {
	ServerAddress string `env:"SERVER_ADDRESS" envDefault:"localhost:8080"`
	AccrualURL    string `env:"ACCRUAL_URL"`
	DatabaseDSN   string `env:"DATABASE_DSN"`
}

func New() *Config {
	var config Config
	err := env.Parse(&config)
	if err != nil {
		log.Fatal(err)
	}
	flagSet := pflag.FlagSet{}
	addrFlag := flagSet.StringP("-addr", "a", config.ServerAddress, "Server address: host:port")
	accrualFlag := flagSet.StringP("-acc", "f", config.AccrualURL, ""+"Accrual service URL")
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
		return errors.New("wrong accural url param")
	}
	return nil
}
