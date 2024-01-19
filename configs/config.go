package configs

import (
	"encoding/json"

	"github.com/spf13/viper"
)

type TokenConfig struct {
	Name        string `json:"name"`
	MaxRequests int    `json:"max_requests"`
	Cooldown    int    `json:"cooldown_seconds"`
}

type conf struct {
	RedisHost     string `mapstructure:"REDIS_HOST"`
	RedisPort     string `mapstructure:"REDIS_PORT"`
	RedisPassword string `mapstructure:"REDIS_PASSWORD"`
	RedisDatabase int    `mapstructure:"REDIS_DATABASE"`

	Tokens      map[string]TokenConfig `mapstructure:"-"`
	MaxRequests int                    `mapstructure:"MAX_REQUESTS"`
	Ttl         int                    `mapstructure:"TTL_SECONDS"`
	Cooldown    int                    `mapstructure:"COOLDOWN_SECONDS"`

	WebServerPort string `mapstructure:"WEB_SERVER_PORT"`
}

func LoadConfig(path string) (*conf, error) {
	var cfg *conf
	viper.SetConfigName("rate_limiter_config")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	err = viper.Unmarshal(&cfg)
	if err != nil {
		panic(err)
	}

	var tkns_cfg []TokenConfig
	configured_tokens := viper.GetString("TOKENS")
	err = json.Unmarshal([]byte(configured_tokens), &tkns_cfg)

	if err != nil {
		panic(err)
	}

	cfg.Tokens = make(map[string]TokenConfig)
	for _, tkn := range tkns_cfg {
		cfg.Tokens[tkn.Name] = tkn
	}

	return cfg, err
}
