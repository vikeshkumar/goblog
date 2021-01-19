package config

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"time"
)

//Configuration interface wrapping some of viper.Viper
type Configuration interface {
	GetString(property string) string
	GetDuration(key string) time.Duration
	GetInt(key string) int
	GetBoolean(key string) bool
}

type cfg struct {
	v *viper.Viper
}

var c *cfg

// GetString return string value for provided key
func (c *cfg) GetString(key string) string {
	return c.v.GetString(key)
}

// GetString return string value for provided key
func (c *cfg) GetInt(key string) int {
	return c.v.GetInt(key)
}

// GetDuration returns the value associated with the key as a duration.
func (c *cfg) GetDuration(key string) time.Duration {
	return c.v.GetDuration(key)
}

// GetDuration returns the value associated with the key as a duration.
func (c *cfg) GetBoolean(key string) bool {
	return c.v.GetBool(key)
}

// Initializes the internal structures
func Init(file string, searchDir string, configType string) (Configuration, error) {
	log.Infof("file: %v, searchDir: %v, configType: %v", file, searchDir, configType)
	v := viper.New()
	v.SetConfigFile(file)
	v.SetConfigType(configType)
	v.AddConfigPath(searchDir)
	e := v.ReadInConfig()
	if e == nil {
		c = &cfg{v}
		return c, nil
	}
	log.Errorf("file: %v, searchDir: %v, configType: %v. %v", file, searchDir, configType, e)
	return nil, e
}

func Get() Configuration {
	return c
}
