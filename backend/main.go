package main

import (
	"github.com/homemusic/backend/cmd"
	"github.com/homemusic/backend/config"
	"github.com/homemusic/backend/internal/common"
	"github.com/homemusic/backend/internal/dao"
	"go.uber.org/zap"
)

func main() {
	// 初始化配置
	if err := config.InitConfig(); err != nil {
		panic("init config failed: " + err.Error())
	}

	// 初始化日志
	common.InitLogger()

	// 初始化数据库
	if err := dao.InitDB(); err != nil {
		common.Logger.Error("init database failed", zap.Error(err))
		panic("init database failed")
	}

	// 启动服务
	cmd.RunServer()
}
