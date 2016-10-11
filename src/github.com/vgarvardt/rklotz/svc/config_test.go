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

func TestNewIniConfigDoesNotExist(t *testing.T) {
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

func TestNewIniEnvConfigDoesNotExist(t *testing.T) {
	initNullLogger()

	Convey("Given a config that does exist", t, func() {
		baseConfigPath := "/tmp/~rklotz-base.ini"
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

func TestNewIniEnvConfigDoesExist(t *testing.T) {
	initNullLogger()

	Convey("Given a config that does exist and env config that does exist either", t, func() {
		var err error
		baseConfigPath := "/tmp/~rklotz-base.ini"
		dataBase := []byte("debug=false\nui.email=vgarvardt@gmail.com\nui.per_page=10\nauth.name=user\n")
		err = ioutil.WriteFile(baseConfigPath, dataBase, 0644)
		defer func() {
			os.Remove(baseConfigPath)
		}()
		So(err, ShouldBeNil)

		envConfigPath := "/tmp/~rklotz-env.ini"
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
