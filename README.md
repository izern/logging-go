# logging

+ wrap for [zap](https://github.com/uber-go/zap)
+ init from viper
+ log level of each module is set independently

## Installation

`go get -u github.com/izern/logging`

## Quick Start

first, init [viper](https://github.com/spf13/viper), example yaml

```yml

logging:
  #  DEBUG
  #  INFO
  #  WARN
  #  ERROR
  #  PANIC
  #  FATAL
  level:
    root: INFO    # default INFO
    module1.child1: DEBUG
    module1.child2: INFO
    module2: ERROR
  encoding: json # json or console, default json. only encoding, console is plan text encoding
  encoder:
    TimeKey: time
    LevelKey: level
    NameKey: logger
    CallerKey: caller
    MessageKey: msg
    StacktraceKey: stacktrace
  output: # default is console
    console: # console or file. console is output on TTY
      async: false # async output,default false
    file:
      path: /tmp/golang/log.log
      async: false # async output,default false
```

and init viper

```golang
  gopath, b := os.LookupEnv("GOPATH")
  assert.Truef(t, b, "no GOPATH ENV")
  viper.SetConfigFile(gopath + "/src/github.com/izern/logging-go/config.yaml")
  err := viper.ReadInConfig()
```

then init logging, and get logger

```golang
  InitZapLoggerFromViper(viper.GetViper())
  log := GetLogger("module1.child1")
  log = GetLogger("module1.child1.child")
  log = GetLogger("test")
```

enjoy~ beginning your APP
