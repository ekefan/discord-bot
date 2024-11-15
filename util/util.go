package util

import (
	"log/slog"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	AppID int `mapstructure:"APP_ID"`
	DiscordToken string `mapstructure:"BOT_TOKEN"`
	PublicKey string `mapstructure:"PUBLIC_KEY"`
	discordBaseUrl string `mapstructure:"DISCORD_BASE_URL"`
}

// LoadConfig reads environment config from bot.env or loads them from
// the os environment variables
func LoadConfig() *Config{
	viper.SetConfigName("bot")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")

	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		slog.Error("cannot read in config variables", "error", err)
		return nil
	}
	config := Config{}
	err := viper.Unmarshal(&config)
	if err != nil {
		slog.Error("unable to decode config into struct, %v", "details", err.Error())
		os.Exit(1)
	}
	return &config
}
