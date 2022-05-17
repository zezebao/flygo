package fllog

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var MainLogger *zap.Logger
var GatewayLogger *zap.Logger
var RedisLogger *zap.Logger
var LmdbLogger *zap.Logger
var HttpLogger *zap.Logger
var ErrLogger *zap.Logger

var Main *zap.SugaredLogger
var Gateway *zap.SugaredLogger
var Redis *zap.SugaredLogger
var Lmdb *zap.SugaredLogger
var Http *zap.SugaredLogger
var Err *zap.SugaredLogger

func init() {
	fmt.Println("##zap logger inited##")
	MainLogger = NewLogger("./logs/main.log", zapcore.DebugLevel, 128, 30, 7, true, "Main")
	GatewayLogger = NewLogger("./logs/gateway.log", zapcore.DebugLevel, 128, 30, 7, true, "Gateway")
	RedisLogger = NewLogger("./logs/redis.log", zapcore.DebugLevel, 128, 30, 7, true, "redis")
	LmdbLogger = NewLogger("./logs/lmdb.log", zapcore.DebugLevel, 128, 30, 7, true, "lmdb")
	HttpLogger = NewLogger("./logs/http.log", zapcore.DebugLevel, 128, 30, 7, true, "http")
	ErrLogger = NewLogger("./logs/error.log", zapcore.DebugLevel, 128, 30, 30, true, "error")

	Main = MainLogger.Sugar()
	Gateway = GatewayLogger.Sugar()
	Redis = RedisLogger.Sugar()
	Lmdb = LmdbLogger.Sugar()
	Http = HttpLogger.Sugar()
	Err = ErrLogger.Sugar()
}