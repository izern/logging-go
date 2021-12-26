package main

import (
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
	"os"
	"testing"
)

func init() {

}

// test init by default
func TestInitZapLoggerFromViper1(t *testing.T) {
	InitZapLoggerFromViper(viper.GetViper())
	assert.NotNil(t, logger)
	logger.Sugar().Info("TestInitZapLoggerFromViper")
	logger.Sugar().Info("TestInitZapLoggerFromViper")
}

// test init by config.yaml
func TestInitZapLoggerFromViper2(t *testing.T) {
	gopath, b := os.LookupEnv("GOPATH")
	assert.Truef(t, b, "no GOPATH ENV")
	viper.SetConfigFile(gopath + "/src/github.com/izern/logging-go/config.yaml")
	err := viper.ReadInConfig()
	assert.NoError(t, err)
	InitZapLoggerFromViper(viper.GetViper())
	assert.NotNil(t, logger)
	logger.Sugar().Info("TestInitZapLoggerFromViper")
}

// test by
func TestGetLevel(t *testing.T) {
	TestInitZapLoggerFromViper2(t)
	level := GetLevel("test")
	assert.Equal(t, level, zapcore.InfoLevel, level.String())

	level = GetLevel("module1.child1")
	assert.Equal(t, level, zapcore.DebugLevel, level.String())

	level = GetLevel("module1.child1.child")
	assert.Equal(t, level, zapcore.DebugLevel, level.String())

	level = GetLevel("module")
	assert.Equal(t, level, zapcore.InfoLevel, level.String())

	level = GetLevel("module2")
	assert.Equal(t, level, zapcore.ErrorLevel, level.String())
}

func TestGetLogger(t *testing.T) {
	TestInitZapLoggerFromViper2(t)
	log := GetLogger("test")
	assert.NotNil(t, log)
	assert.True(t, log.Core().Enabled(zapcore.InfoLevel))
	assert.False(t, log.Core().Enabled(zapcore.DebugLevel))
	log.Info("INFO must be show")
	log.Debug("DEBUG must be hide")

	log = GetLogger("module1.child1")
	assert.NotNil(t, log)
	assert.True(t, log.Core().Enabled(zapcore.DebugLevel))
	log.Debug("DEBUG must be show")

	log = GetLogger("module1.child1.child")
	assert.NotNil(t, log)
	assert.True(t, log.Core().Enabled(zapcore.DebugLevel))

	log = GetLogger("module2")
	assert.NotNil(t, log)
	assert.True(t, log.Core().Enabled(zapcore.ErrorLevel))
	assert.False(t, log.Core().Enabled(zapcore.WarnLevel))

	log2 := GetLogger("module2")
	assert.True(t, log == log2)

	log3 := GetLogger("module3")
	assert.False(t, log == log3)

}
