package config

import (
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Environment string         `mapstructure:"environment"`
	HTTP        HTTPConfig     `mapstructure:"http"`
	Postgres    PostgresConfig `mapstructure:"postgres"`
	Redis       RedisConfig    `mapstructure:"redis"`
}

type HTTPConfig struct {
	Port         string `mapstructure:"port"`
	ReadTimeout  int    `mapstructure:"read_timeout"`  // in seconds
	WriteTimeout int    `mapstructure:"write_timeout"` // in seconds
	IdleTimeout  int    `mapstructure:"idle_timeout"`  // in seconds
}

type PostgresConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	SSLMode  string `mapstructure:"sslmode"`
	TimeZone string `mapstructure:"timezone"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

func Load(configPath string) (*Config, error) {
	v := viper.New()
	v.SetEnvPrefix("APP") // Prefix for env vars (e.g., APP_ENVIRONMENT, APP_HTTP_PORT)
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	if configPath != "" {
		v.SetConfigFile(configPath)
		v.SetConfigType("yaml")
		_ = v.ReadInConfig() // Non-fatal if missing
	}

	setDefaults(v)
	bindEnvVars(v)

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func setDefaults(v *viper.Viper) {
	v.SetDefault("environment", "development")

	// HTTP defaults
	v.SetDefault("http.port", "8080")
	v.SetDefault("http.read_timeout", 15)
	v.SetDefault("http.write_timeout", 15)
	v.SetDefault("http.idle_timeout", 60)

	// Postgres defaults
	v.SetDefault("postgres.host", "localhost")
	v.SetDefault("postgres.port", 5432)
	v.SetDefault("postgres.user", "postgres")
	v.SetDefault("postgres.password", "postgres")
	v.SetDefault("postgres.dbname", "mydb")
	v.SetDefault("postgres.sslmode", "disable")
	v.SetDefault("postgres.timezone", "UTC")

	// Redis defaults
	v.SetDefault("redis.host", "localhost")
	v.SetDefault("redis.port", 6379)
	v.SetDefault("redis.password", "")
	v.SetDefault("redis.db", 0)
}

// bindEnvVars binds environment variables for all config fields.
func bindEnvVars(v *viper.Viper) {
	keys := []string{
		"environment",
		"http.port",
		"http.read_timeout",
		"http.write_timeout",
		"http.idle_timeout",
		"postgres.host",
		"postgres.port",
		"postgres.user",
		"postgres.password",
		"postgres.dbname",
		"postgres.sslmode",
		"postgres.timezone",
		"redis.host",
		"redis.port",
		"redis.password",
		"redis.db",
	}
	for _, key := range keys {
		_ = v.BindEnv(key)
	}
}
