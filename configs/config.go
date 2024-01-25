package configs

import (
	"encoding/json"
	"log"

	"github.com/spf13/viper"
)

type TokenConfig struct {
	Name        string `json:"name"`
	MaxRequests int    `json:"max_requests"`
	Cooldown    int    `json:"cooldown_seconds"`
}

type Config struct {
	DBHost     string `mapstructure:"DB_HOST"`
	DBPort     string `mapstructure:"DB_PORT"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBName     string `mapstructure:"DB_NAME"`
	DBUser     string `mapstructure:"DB_USER"`

	Tokens      map[string]TokenConfig `mapstructure:"-"`
	MaxRequests int                    `mapstructure:"MAX_REQUESTS"`
	Ttl         int                    `mapstructure:"TTL_SECONDS"`
	Cooldown    int                    `mapstructure:"COOLDOWN_SECONDS"`

	WebServerPort string `mapstructure:"WEB_SERVER_PORT"`
}

func LoadConfig(path string) (*Config, error) {
	var cfg *Config
	viper.SetConfigName("rate_limiter_config")
	viper.SetConfigType("env")
	viper.SetConfigFile(path + "/.env")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()

	if err != nil {
		log.Panic(err)
	}
	err = viper.Unmarshal(&cfg)
	if err != nil {
		log.Panic(err)
	}

	var tkns_cfg []TokenConfig
	configured_tokens := viper.GetString("TOKENS")
	err = json.Unmarshal([]byte(configured_tokens), &tkns_cfg)

	if err != nil {
		log.Panic(err)
	}

	cfg.Tokens = make(map[string]TokenConfig)
	for _, tkn := range tkns_cfg {
		cfg.Tokens[tkn.Name] = tkn
	}

	return cfg, err
}
