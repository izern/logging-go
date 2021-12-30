package example

import (
	"github.com/izern/logging"
	"github.com/spf13/viper"
	"go.uber.org/zap/zapcore"
	"testing"
)

var l = logging.GetLogger("example1")

func init() {

}

func TestExample1(t *testing.T) {
	log := logging.GetLogger("example1")
	log.Info("111")

	viper.Set("logging.level", "ERROR")

	logging.InitZapLoggerFromViper(viper.GetViper())

	log = logging.GetLogger("example1")
	enabled := log.Core().Enabled(zapcore.InfoLevel)
	println(enabled)
}

func TestRefresh(t *testing.T) {

	logger1 := logging.GetLogger("example2")
	enabled := logger1.Core().Enabled(zapcore.InfoLevel)
	println(enabled)

	viper.Set("logging.level", "ERROR")
	logging.InitZapLoggerFromViper(viper.GetViper())

	enabled = logger1.Core().Enabled(zapcore.InfoLevel)
	println(enabled)

}
