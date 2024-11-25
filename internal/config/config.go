package config

import "github.com/spf13/viper"

type Config struct {
	Port             uint16 `mapstructure:"PORT"`
	DBUser           string `mapstructure:"DB_USER"`
	DBPassword       string `mapstructure:"DB_PASSWORD"`
	DBHost           string `mapstructure:"DB_HOST"`
	DBPort           uint16 `mapstructure:"DB_PORT"`
	DBName           string `mapstructure:"DB_NAME"`
	ExternalAPIURL   string `mapstructure:"EXTERNAL_API_URL"`
	MaxGroupNameLen  int    `mapstructure:"MAX_GROUP_NAME_LEN"`
	MaxSongNameLen   int    `mapstructure:"MAX_SONG_NAME_LEN"`
	MaxSongLyricsLen int    `mapstructure:"MAX_SONG_LYRICS_LEN"`
	MaxSongLinkLen   int    `mapstructure:"MAX_SONG_LINK_LEN"`
}

func Load() (config Config, err error) {
	viper.SetConfigFile(".env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return
	}

	return
}
