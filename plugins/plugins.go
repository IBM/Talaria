package plugins

import (
	"github.com/ibm/opentalaria/protocol"
	"github.com/spf13/viper"
)

type PluginInterface interface {
	Init(env *viper.Viper) error

	// settings
	GetSetting(key string) (string, error)
	SetSetting(key, value string) error

	// topics
	AddTopic(topic protocol.CreatableTopic) error
}
