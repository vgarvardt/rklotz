package svc

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/drgomesp/cargo/container"
	. "github.com/smartystreets/goconvey/convey"
)

func initNullLogger() {
	Container = container.New()
	Container.Set(DI_LOGGER, NewNullLogger())
}

func TestNewIniConfigBaseDoesNotExist(t *testing.T) {
	initNullLogger()

	Convey("Given a config path that does not exist", t, func() {
		baseConfigPath := "/tmp/must-not-exist.ini"

		Convey("When trying to crrate new ini config loader", func() {
			So(func() {
				NewIniConfig(baseConfigPath, "")
			}, ShouldPanic)
		})
	})
}

func TestNewIniConfigEnvDoesNotExist(t *testing.T) {
	initNullLogger()

	Convey("Given a config that does exist", t, func() {
		baseConfigPath := "/tmp/~rklotz-ini-base-no-env.ini"
		dataBase := []byte("debug=false\nui.email=vgarvardt@gmail.com\nui.per_page=10\n")
		err := ioutil.WriteFile(baseConfigPath, dataBase, 0644)
		defer func() {
			os.Remove(baseConfigPath)
		}()
		So(err, ShouldBeNil)

		Convey("When trying to create new ini config loader", func() {
			var config *iniConfig
			So(func() {
				config = NewIniConfig(baseConfigPath, "")
			}, ShouldNotPanic)

			Convey("Then it should return default config values", func() {
				So(config.Bool("debug"), ShouldEqual, false)
				So(config.String("ui.email"), ShouldEqual, "vgarvardt@gmail.com")
				So(config.Int("ui.per_page"), ShouldEqual, 10)
			})
		})
	})
}

func TestNewIniConfigEnvDoesExist(t *testing.T) {
	initNullLogger()

	Convey("Given a config that does exist and env config that does exist either", t, func() {
		var err error
		baseConfigPath := "/tmp/~rklotz-ini-base-env.ini"
		dataBase := []byte("debug=false\nui.email=vgarvardt@gmail.com\nui.per_page=10\nauth.name=user\n")
		err = ioutil.WriteFile(baseConfigPath, dataBase, 0644)
		defer func() {
			os.Remove(baseConfigPath)
		}()
		So(err, ShouldBeNil)

		envConfigPath := "/tmp/~rklotz-ini-env.ini"
		dataEnv := []byte("debug=true\nui.email=test@example.com\nui.per_page=5\n")
		err = ioutil.WriteFile(envConfigPath, dataEnv, 0644)
		defer func() {
			os.Remove(envConfigPath)
		}()
		So(err, ShouldBeNil)

		Convey("When trying to create new ini config loader", func() {
			var config *iniConfig
			So(func() {
				config = NewIniConfig(baseConfigPath, envConfigPath)
			}, ShouldNotPanic)

			Convey("Then it should return overridden config values", func() {
				So(config.Bool("debug"), ShouldEqual, true)
				So(config.String("ui.email"), ShouldEqual, "test@example.com")
				So(config.Int("ui.per_page"), ShouldEqual, 5)
				So(config.String("auth.name"), ShouldEqual, "user")
			})
		})
	})
}

func TestNewIniEnvConfigDoesNotExist(t *testing.T) {
	initNullLogger()

	Convey("Given a config path that does not exist", t, func() {
		baseConfigPath := "/tmp/must-not-exist.ini"

		Convey("When trying to crrate new ini config loader", func() {
			So(func() {
				NewIniEnvConfig(baseConfigPath, ENV_PREFIX)
			}, ShouldPanic)
		})
	})
}

func TestNewIniEnvConfigEnvNotSetOrEmpty(t *testing.T) {
	initNullLogger()

	Convey("Given a config that does exist", t, func() {
		configPath := "/tmp/~rklotz-ini-env-config-no-env.ini"
		dataBase := []byte("debug=false\nui.email=vgarvardt@gmail.com\nui.per_page=10\n")
		err := ioutil.WriteFile(configPath, dataBase, 0644)
		defer func() {
			os.Remove(configPath)
		}()
		So(err, ShouldBeNil)

		Convey("When trying to create new ini config loader and env not set or empty", func() {
			var config *iniEnvConfig
			So(func() {
				config = NewIniEnvConfig(configPath, ENV_PREFIX)
			}, ShouldNotPanic)

			os.Setenv(config.getEnvKey("debug"), "")
			os.Setenv(config.getEnvKey("ui.email"), "")
			os.Unsetenv(config.getEnvKey("ui.per_page"))

			Convey("Then it should return default config values", func() {
				So(config.Bool("debug"), ShouldEqual, false)
				So(config.String("ui.email"), ShouldEqual, "vgarvardt@gmail.com")
				So(config.Int("ui.per_page"), ShouldEqual, 10)
			})
		})
	})
}

func TestNewIniEnvConfigEnvIsSet(t *testing.T) {
	initNullLogger()

	Convey("Given a config that does exist", t, func() {
		var err error
		configPath := "/tmp/~rklotz-ini-env-config.ini"
		dataBase := []byte("debug=false\nui.email=vgarvardt@gmail.com\nui.per_page=10\nauth.name=user\n")
		err = ioutil.WriteFile(configPath, dataBase, 0644)
		defer func() {
			os.Remove(configPath)
		}()
		So(err, ShouldBeNil)

		Convey("When trying to create new ini config loader and env values are set", func() {
			var config *iniEnvConfig
			So(func() {
				config = NewIniEnvConfig(configPath, ENV_PREFIX)
			}, ShouldNotPanic)

			os.Setenv(config.getEnvKey("debug"), "1")
			os.Setenv(config.getEnvKey("ui.email"), "test@example.com")
			os.Setenv(config.getEnvKey("ui.per_page"), "5")

			Convey("Then it should return overridden config values", func() {
				So(config.Bool("debug"), ShouldEqual, true)
				So(config.String("ui.email"), ShouldEqual, "test@example.com")
				So(config.Int("ui.per_page"), ShouldEqual, 5)
				So(config.String("auth.name"), ShouldEqual, "user")
			})
		})
	})
}
