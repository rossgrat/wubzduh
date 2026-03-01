package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Server  ServerConfig  `mapstructure:"server"`
	DB      DBConfig      `mapstructure:"db"`
	Spotify SpotifyConfig `mapstructure:"spotify"`
	Log     LogConfig     `mapstructure:"log"`
}

type ServerConfig struct {
	Port        string `mapstructure:"port"`
	TemplateDir string `mapstructure:"template_dir"`
	StaticDir   string `mapstructure:"static_dir"`
}

type DBConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`
}

func (c DBConfig) ConnStr() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		c.User, c.Password, c.Host, c.Port, c.Name,
	)
}

type SpotifyConfig struct {
	ClientID     string `mapstructure:"client_id"`
	ClientSecret string `mapstructure:"client_secret"`
}

type LogConfig struct {
	Path       string `mapstructure:"path"`
	MaxSizeMB  int    `mapstructure:"max_size_mb"`
	MaxBackups int    `mapstructure:"max_backups"`
}

func Load(path string) (*Config, error) {
	viper.SetConfigFile(path)

	viper.BindEnv("db.user", "DB_USER")
	viper.BindEnv("db.password", "DB_PASSWORD")
	viper.BindEnv("spotify.client_id", "SPOTIFY_CLIENT_ID")
	viper.BindEnv("spotify.client_secret", "SPOTIFY_CLIENT_SECRET")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("reading config file: %w", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshaling config: %w", err)
	}

	return &cfg, nil
}
