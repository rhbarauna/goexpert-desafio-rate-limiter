package configs

import "github.com/spf13/viper"

type conf struct {
	// DBDriver          string `mapstructure:"DB_DRIVER"`
	// DBHost            string `mapstructure:"DB_HOST"`
	// DBPort            string `mapstructure:"DB_PORT"`
	// DBUser            string `mapstructure:"DB_USER"`
	// DBPassword        string `mapstructure:"DB_PASSWORD"`
	// DBName            string `mapstructure:"DB_NAME"`

	// GRPCServerPort    string `mapstructure:"GRPC_SERVER_PORT"`
	// GraphQLServerPort string `mapstructure:"GRAPHQL_SERVER_PORT"`
	RedisHost     string `mapstructure:"REDIS_HOST"`
	RedisPort     string `mapstructure:"REDIS_PORT"`
	RedisPassword string `mapstructure:"REDIS_PASSWORD"`
	RedisDatabase int    `mapstructure:"REDIS_DATABASE"`

	IpRequestLimit       int    `mapstructure:"IP_REQUEST_LIMIT"`
	IpDefaultWindow      int    `mapstructure:"IP_DEFAULT_WINDOW_SECONDS"`
	IpDefaultCooldown    int    `mapstructure:"IP_DEFAULT_COOLDOWN_SECONDS"`
	TokenDefaultCooldown int    `mapstructure:"TOKEN_DEFAULT_COOLDOWN_SECONDS"`
	WebServerPort        string `mapstructure:"WEB_SERVER_PORT"`
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
	return cfg, err
}
