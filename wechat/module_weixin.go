package wechat

import (
	"com.ffl/modules"

	"github.com/zezebao/flygo/logger"
)

type WXModuleConfig struct {
	modules.ModuleConfig

	AppId     string
	AppSecert string
}

type weiXinModule struct {
	modules.BaseModule
}

func (module *weiXinModule) Init(config *WXModuleConfig) {
	module.BaseModule.Init(&config.ModuleConfig)

	logger.MainLogger.Debug("---weixin module setup")

	NewWxJsSdk(config.AppId, config.AppSecert)
}

var Module = &weiXinModule{}
