package svc

import (
	"fmt"
	"os"
	"strconv"

	log "github.com/Sirupsen/logrus"
	"gopkg.in/ini.v1"
)

const ENV_PREFIX = "RKLOTZ_"

type Config interface {
	String(key string) string
	Int(key string) int
	Bool(key string) bool
}

func NewIniConfig(baseConfigPath, envConfigPath string) *iniConfig {
	logger := Container.MustGet(DI_LOGGER).(*log.Logger)

	var config = &iniConfig{}
	var err error
	logger.WithField("path", baseConfigPath).Info("Loading base config")
	if config.iniFile, err = ini.Load(baseConfigPath); err != nil {
		panic(err)
	}

	logger.WithField("path", envConfigPath).Info("Loading env config")
	if _, err := os.Stat(envConfigPath); os.IsNotExist(err) {
		logger.WithField("path", envConfigPath).Warn("Env config not found")
	} else {
		if err := config.iniFile.Append(envConfigPath); err != nil {
			logger.WithField("err", err).Fatal("Failed to append env config")
		}
	}

	return config
}

func NewIniEnvConfig(configPath, envPrefix string) *iniEnvConfig {
	logger := Container.MustGet(DI_LOGGER).(*log.Logger)

	logger.WithFields(log.Fields{"path": configPath, "prefix": envPrefix}).Info("Loading ini config for env loader")

	var config = &iniEnvConfig{envPrefix: envPrefix}
	config.iniConfig = NewIniConfig(configPath, "")

	return config
}

type iniConfig struct {
	iniFile *ini.File
}

func (config *iniConfig) configKey(key string) *ini.Key {
	return config.iniFile.Section("").Key(key)
}

func (config *iniConfig) String(key string) string {
	return config.configKey(key).String()
}

func (config *iniConfig) Int(key string) int {
	if val, err := config.configKey(key).Int(); err != nil {
		panic(err)
	} else {
		return val
	}
}

func (config *iniConfig) Bool(key string) bool {
	if val, err := config.configKey(key).Bool(); err != nil {
		panic(err)
	} else {
		return val
	}
}

type iniEnvConfig struct {
	*iniConfig
	envPrefix string
}

func (config *iniEnvConfig) getEnvKey(key string) string {
	return fmt.Sprintf("%s%s", config.envPrefix, key)
}

func (config *iniEnvConfig) String(key string) string {
	envValue := os.Getenv(config.getEnvKey(key))
	if envValue == "" {
		return config.iniConfig.String(key)
	}
	return envValue
}

func (config *iniEnvConfig) Int(key string) int {
	envValue := os.Getenv(config.getEnvKey(key))
	if envValue == "" {
		return config.iniConfig.Int(key)
	} else {
		if val, err := strconv.Atoi(envValue); err != nil {
			panic(err)
		} else {
			return val
		}
	}
}

func (config *iniEnvConfig) Bool(key string) bool {
	envValue := os.Getenv(config.getEnvKey(key))
	if envValue == "" {
		return config.iniConfig.Bool(key)
	} else {
		if val, err := strconv.ParseBool(envValue); err != nil {
			panic(err)
		} else {
			return val
		}
	}
}
