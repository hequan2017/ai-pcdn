package plugin

import "github.com/flipped-aurora/gin-vue-admin/server/plugin/pcdn/config"

// Config 插件全局配置，由 initialize/viper.go 反序列化 config.yaml 的 pcdn: 段填充
var Config config.Config
