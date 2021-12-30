package logging

import (
	"fmt"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"strings"
	"sync"
)

// module : log level
var moduleLevelMap sync.Map

// log level : *logger
var loggerLevelInstanceMap sync.Map

var moduleLoggerInstanceMap sync.Map

// RootModule default module
var RootModule = "root"

// default global logger instance
var logger *zap.Logger

var defaultLevel = zap.InfoLevel

var initOnce = sync.Once{}

func init() {
	zap.NewProductionConfig()
	logger, _ = zap.NewProduction()
	loggerLevelInstanceMap.Store(defaultLevel, logger)
	moduleLevelMap.Store(RootModule, defaultLevel)
}

// InitZapLoggerFromViper InitZapLoggerFromViper, should be call once
func InitZapLoggerFromViper(viper *viper.Viper, options ...zap.Option) {

	initOnce.Do(func() {
		// set custom default value
		initDefaultValue(viper)
		log := initLog(viper, options...)
		levelValue := viper.Get("logging.level")

		switch levelValue.(type) {
		case string:
			var defaultLevel zapcore.Level
			err := defaultLevel.Set(levelValue.(string))
			if err != nil {
				panic(err)
			}
			moduleLevelMap.Store(RootModule, defaultLevel)
			loggerLevelInstanceMap.Store(defaultLevel, log.WithOptions(zap.IncreaseLevel(defaultLevel)))
		case map[string]string, map[string]interface{}:
			m := viper.GetStringMapString("logging.level")
			for k, v := range m {
				err := setModuleLevel(k, v, log)
				if err != nil {
					panic(fmt.Sprintf("cannot set module {%v}log level: %v", k, err.Error()))
				}
			}
			if _, ok := moduleLevelMap.Load(RootModule); !ok {
				moduleLevelMap.Store(RootModule, zapcore.InfoLevel)
				loggerLevelInstanceMap.Store(zapcore.InfoLevel, log.WithOptions(zap.IncreaseLevel(zapcore.InfoLevel)))
			}
		default:
			moduleLevelMap.Store(RootModule, zapcore.InfoLevel)
			loggerLevelInstanceMap.Store(zapcore.InfoLevel, log.WithOptions(zap.IncreaseLevel(zapcore.InfoLevel)))
		}
		logger = GetLogger(RootModule)
		zap.ReplaceGlobals(logger)
		refresh()
	})
}

// GetLevel get module log level
func GetLevel(module string) zapcore.Level {
	// Closest module
	var matchModule = RootModule
	// Have the most identical prefixes
	var maxLength = 0
	moduleLevelMap.Range(func(key, value interface{}) bool {
		if key == module {
			matchModule = key.(string)
			return false
		}
		prefix := fmt.Sprintf("%v.", key.(string))
		if strings.HasPrefix(module, prefix) {
			length := len(prefix)
			if length > maxLength {
				matchModule = key.(string)
				maxLength = length
			}
		}
		return true
	})
	value, _ := moduleLevelMap.Load(matchModule)
	return value.(zapcore.Level)
}

// GetLogger get module logger
func GetLogger(module string) *zap.Logger {

	cache, ok := moduleLoggerInstanceMap.Load(module)
	if ok {
		return cache.(*zap.Logger)
	} else {
		level := GetLevel(module)
		value, _ := loggerLevelInstanceMap.Load(level)
		log := &*value.(*zap.Logger)
		moduleLoggerInstanceMap.Store(module, log)
		return log
	}
}

func refresh() {
	moduleLoggerInstanceMap.Range(func(key, value interface{}) bool {
		module := key.(string)
		oldLog := value.(*zap.Logger)
		newLevel := GetLevel(module)
		newLog, _ := loggerLevelInstanceMap.Load(newLevel)
		*oldLog = *newLog.(*zap.Logger)
		return true
	})

}

func initDefaultValue(viper *viper.Viper) {
	viper.Set("logging.level", "INFO")
	viper.Set("logging.encoding", "json")
	if !viper.IsSet("logging.output") {
		output := make(map[string]interface{})
		output["console"] = map[string]interface{}{"async": true}
		viper.Set("logging.output", output)
	}
}

// setModuleLevel
func setModuleLevel(module string, logLevel string, log *zap.Logger) error {
	var level zapcore.Level
	err := level.Set(logLevel)
	if err == nil {
		moduleLevelMap.Store(module, level)
		loggerLevelInstanceMap.Store(level, log.WithOptions(zap.IncreaseLevel(level)))
	}

	return err
}

func initLog(viper *viper.Viper, options ...zap.Option) *zap.Logger {

	var encoderConfig *zapcore.EncoderConfig
	encoder := viper.Get("logging.encoder")
	if encoder != nil {
		encoderConfig = &zapcore.EncoderConfig{}
		err := viper.UnmarshalKey("logging.encoder", encoderConfig)
		if err != nil {
			panic(fmt.Sprintf("UnmarshalKey logging.encoder failed %v", err.Error()))
		}
		encoderConfig.LineEnding = zapcore.DefaultLineEnding
		encoderConfig.FunctionKey = zapcore.OmitKey
		encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder       // ISO8601 UTC
		encoderConfig.EncodeDuration = zapcore.NanosDurationEncoder //
		encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder     // ShortCallerEncoder
		encoderConfig.EncodeName = zapcore.FullNameEncoder
		encoderConfig.EncodeLevel = zapcore.LowercaseLevelEncoder
	} else {
		encoderConfigTmp := zap.NewProductionEncoderConfig()
		encoderConfig = &encoderConfigTmp
	}
	var enc zapcore.Encoder
	encoding := viper.GetString("logging.encoding")
	if encoding == "console" {
		enc = zapcore.NewConsoleEncoder(*encoderConfig)
	} else {
		enc = zapcore.NewJSONEncoder(*encoderConfig)
	}

	var writers = make([]zapcore.WriteSyncer, 0)
	hasConsole := viper.IsSet("logging.output.console")
	if hasConsole {
		async := viper.GetBool("logging.output.console.async")
		if async {
			writers = append(writers, zapcore.Lock(os.Stdout))
		} else {
			writers = append(writers, zapcore.AddSync(os.Stdout))
		}
	}
	if viper.IsSet("logging.output.file") {
		filePath := viper.GetString("logging.output.file.path")
		async := viper.GetBool("logging.output.file.async")
		logFileLocation, err := openFileSafe(filePath)
		if err != nil {
			panic(fmt.Sprintf("logging.output.file.path set error. %v", err))
		}
		if async {
			writers = append(writers, zapcore.Lock(logFileLocation))
		} else {
			writers = append(writers, zapcore.AddSync(logFileLocation))
		}
	}

	core := zapcore.NewCore(
		enc,
		zapcore.NewMultiWriteSyncer(writers...), // 打印到控制台和文件
		zapcore.DebugLevel,                      // 日志级别
	)
	log := zap.New(core, zap.AddCaller()).WithOptions(options...)
	return log
}
