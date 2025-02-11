package config

import (
	"log/slog"
	"strings"

	"github.com/google/uuid"
	"github.com/ibm/opentalaria/plugins"
	"github.com/ibm/opentalaria/plugins/postgresql"
	"github.com/ibm/opentalaria/utils"
	"github.com/spf13/viper"
)

type Config struct {
	OTProfile       OTProfile
	LogLevel        slog.Level
	LogFormat       string
	DebugServerPort int

	Broker  *Broker
	Cluster *Cluster

	Env *viper.Viper

	Plugin plugins.PluginInterface
}

type Cluster struct {
	ClusterID string
}

func NewConfig(confFilename string) (*Config, error) {
	config := Config{}

	// init viper
	env := viper.New()

	env.AutomaticEnv()
	env.SetEnvPrefix("ot")
	env.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	env.SetConfigType("yaml")
	env.SetConfigFile(confFilename)
	env.AddConfigPath(".")

	// set defaults for configuration properties
	setDefaults(env)

	env.ReadInConfig()

	config.Env = env

	config.loadProfile()
	config.loadLogLevel()
	config.LogFormat = env.GetString("log.format")
	config.DebugServerPort = env.GetInt("debug.server.port")

	broker, err := NewBroker(env)
	if err != nil {
		return &Config{}, err
	}

	config.Broker = broker

	clusterId := env.GetString("cluster.id")
	if clusterId == "" {
		uid, err := uuid.NewV6()
		if err != nil {
			return &Config{}, err
		}

		clusterId = uid.String()
	}

	config.Cluster = &Cluster{
		ClusterID: clusterId,
	}

	// load the plugin. For now the plugins are statically defined in the plugins package,
	// but this will change in the future.
	pluginsConf := make(map[string]interface{})
	err = env.UnmarshalKey("plugins", &pluginsConf)
	if err != nil {
		return &Config{}, err
	}

	k := utils.MapKeys(pluginsConf)

	// In the first implementation of the plugin system, we will load only the first defined plugin.
	// Later this will be changed to implement a plugin/middleware chain where multiple plugins will be loaded
	// and chained together.
	if len(k) > 0 {
		switch k[0] {
		case "postgres":
			plugin := postgresql.Plugin{}
			err := plugin.Init(env)
			if err != nil {
				return &Config{}, err
			}

			config.Plugin = &plugin
		}
	}

	return &config, nil
}

// setDefaults sets the default values for properties that are not set.
func setDefaults(env *viper.Viper) {
	env.SetDefault("log.level", "warn")
	env.SetDefault("log.format", "text")
	env.SetDefault("debug.server.port", 9090)
	env.SetDefault("broker.id", -1)
	env.SetDefault("reserved.broker.max.id", 1000)
}

/**
 * Unit test helpers
 */

// MockCluster generates a mock object used for unit testing.
func MockCluster() *Cluster {
	return &Cluster{
		ClusterID: "abc",
	}
}

// MockConfig generates a mock object used for unit testing.
func MockConfig() *Config {
	config := Config{}

	config.Cluster = MockCluster()
	config.Broker = MockBroker()

	return &config
}
