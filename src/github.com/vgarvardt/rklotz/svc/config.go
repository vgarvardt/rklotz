package svc

import (
	"os"

	"gopkg.in/ini.v1"
	log "github.com/Sirupsen/logrus"
)

type Config interface {
	String(key string) string
	Int(key string) int
	Bool(key string) bool
}

type iniConfig struct {
	iniFile *ini.File
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
