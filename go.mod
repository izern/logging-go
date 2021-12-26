module github.com/izern/logging

go 1.15

require (
	github.com/spf13/viper v1.9.0
	github.com/stretchr/testify v1.7.0
	go.uber.org/atomic v1.7.0
	go.uber.org/zap v1.19.1
)

replace go.uber.org/zap v1.19.1 => github.com/uber-go/zap v1.19.1
