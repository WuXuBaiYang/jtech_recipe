package common

import "server/tool"

// DebugMode 调试状态
const DebugMode = true

// ServerHost 服务基本地址
var ServerHost = tool.If[string](DebugMode, devServerHost, productionServerHost)

// 开发地址 api.jtech.live
// const devServerHost = "api.jtech.live"
const devServerHost = "35.73.34.38"

// 生产地址
const productionServerHost = "127.0.0.1"
