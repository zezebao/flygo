package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var MainLogger *zap.Logger
var GatewayLogger *zap.Logger
var RedisLogger *zap.Logger
var LmdbLogger *zap.Logger
var HttpLogger *zap.Logger
var ErrLogger *zap.Logger

func init() {
	MainLogger = NewLogger("./logs/main.log", zapcore.InfoLevel, 128, 30, 7, true, "Main")
	GatewayLogger = NewLogger("./logs/gateway.log", zapcore.DebugLevel, 128, 30, 7, true, "Gateway")
	RedisLogger = NewLogger("./logs/redis.log", zapcore.InfoLevel, 128, 30, 7, true, "redis")
	LmdbLogger = NewLogger("./logs/lmdb.log", zapcore.InfoLevel, 128, 30, 7, true, "lmdb")
	HttpLogger = NewLogger("./logs/http.log", zapcore.DebugLevel, 128, 30, 7, true, "http")
	ErrLogger = NewLogger("./logs/error.log", zapcore.DebugLevel, 128, 30, 30, true, "error")
}