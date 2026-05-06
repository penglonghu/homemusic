package api

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/homemusic/backend/internal/common"
	"github.com/homemusic/backend/internal/service"
	"go.uber.org/zap"
)

var musicService = service.NewMusicService()

// ScanMusic 扫描音乐
func ScanMusic(c *gin.Context) {
	// 获取音乐路径（从初始化配置中获取）
	initStatus, err := service.CheckInitStatus()
	if err != nil {
		common.ServerError(c)
		return
	}

	if !initStatus["is_initialized"].(bool) {
		common.Fail(c, 400, "系统未初始化")
		return
	}

	musicPath := initStatus["music_path"].(string)
	if musicPath == "" {
		common.Fail(c, 400, "音乐目录未配置")
		return
	}

	// 异步执行扫描
	go func() {
		err := musicService.ScanMusicDirectory(musicPath)
		if err != nil {
			common.Logger.Error("音乐扫描失败", zap.Error(err))
		}
	}()

	common.Success(c, map[string]interface{}{
		"message": "音乐扫描已启动",
	})
}

// GetScanStatistics 获取扫描统计
func GetScanStatistics(c *gin.Context) {
	stats, err := musicService.GetScanStatistics()
	if err != nil {
		common.ServerError(c)
		return
	}

	common.Success(c, stats)
}

// GetMusics 获取音乐列表
func GetMusics(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "20")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	musics, total, err := musicService.GetMusics(page, pageSize)
	if err != nil {
		common.ServerError(c)
		return
	}

	common.Success(c, map[string]interface{}{
		"musics":    musics,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
		"pages":     (total + int64(pageSize) - 1) / int64(pageSize),
	})
}

// CleanDeletedFiles 清理已删除的文件
func CleanDeletedFiles(c *gin.Context) {
	err := musicService.CleanDeletedFiles()
	if err != nil {
		common.ServerError(c)
		return
	}

	common.Success(c, map[string]interface{}{
		"message": "清理完成",
	})
}