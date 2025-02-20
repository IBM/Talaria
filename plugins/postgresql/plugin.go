package postgresql

import (
	"database/sql"

	"github.com/spf13/viper"
	"golang.org/x/exp/slog"
)

type Plugin struct {
	config config
	db     *sql.DB
}

type config struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
}

func (p *Plugin) Init(env *viper.Viper) error {
	parsedConf := make(map[string]*config)
	err := env.UnmarshalKey("plugins", &parsedConf)
	if err != nil {
		return err
	}

	p.config = *parsedConf["postgres"]

	err = p.initConnection()
	if err != nil {
		return err
	}

	slog.Info("Postgres plugin initialized")

	return nil
}

func (p *Plugin) GetSetting(key string) (string, error) {
	return "", nil
}

func (p *Plugin) SetSetting(key, value string) error {
	return nil
}
