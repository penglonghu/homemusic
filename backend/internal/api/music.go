package api

import (
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/homemusic/backend/internal/common"
	"github.com/homemusic/backend/internal/service"
	"go.uber.org/zap"
)

var musicService *service.MusicService

func SetMusicService(s *service.MusicService) {
	musicService = s
}

func getMusicService() *service.MusicService {
	if musicService == nil {
		musicService = service.NewMusicService()
	}
	return musicService
}

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
		err := getMusicService().ScanMusicDirectory(musicPath)
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
	stats, err := getMusicService().GetScanStatistics()
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

	musics, total, err := getMusicService().GetMusics(page, pageSize)
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

// BrowseMusic 浏览分类：all、artist、album
func BrowseMusic(c *gin.Context) {
	category := strings.ToLower(c.DefaultQuery("category", "all"))
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

	items, total, err := getMusicService().BrowseCategory(category, page, pageSize)
	if err != nil {
		common.ServerError(c)
		return
	}
	if items == nil {
		common.ParamError(c, "category 参数错误，支持 all/artist/album")
		return
	}

	common.Success(c, map[string]interface{}{
		"category":  category,
		"items":     items,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
		"pages":     (total + int64(pageSize) - 1) / int64(pageSize),
	})
}

// SearchMusic 全局搜索
func SearchMusic(c *gin.Context) {
	keyword := strings.TrimSpace(c.Query("keyword"))
	if keyword == "" {
		common.ParamError(c, "keyword 不能为空")
		return
	}
	if len(keyword) > 128 {
		common.ParamError(c, "keyword 长度不能超过 128")
		return
	}

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

	musics, total, err := getMusicService().SearchMusic(keyword, page, pageSize)
	if err != nil {
		common.ServerError(c)
		return
	}

	common.Success(c, map[string]interface{}{
		"keyword":   keyword,
		"musics":    musics,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
		"pages":     (total + int64(pageSize) - 1) / int64(pageSize),
	})
}

// GetRecentPlayRecords 获取最近播放记录
func GetRecentPlayRecords(c *gin.Context) {
	userIDValue, exists := c.Get("user_id")
	if !exists {
		common.Unauthorized(c)
		return
	}

	userID, ok := userIDValue.(uint)
	if !ok {
		common.ServerError(c)
		return
	}

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

	records, total, err := getMusicService().GetRecentPlays(userID, page, pageSize)
	if err != nil {
		common.ServerError(c)
		return
	}

	common.Success(c, map[string]interface{}{
		"records":   records,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
		"pages":     (total + int64(pageSize) - 1) / int64(pageSize),
	})
}

// RecordPlayHistoryRequest 记录播放历史请求
type RecordPlayHistoryRequest struct {
	MusicID uint `json:"music_id" binding:"required,gt=0"`
}

// RecordPlayHistory 记录播放历史
func RecordPlayHistory(c *gin.Context) {
	userIDValue, exists := c.Get("user_id")
	if !exists {
		common.Unauthorized(c)
		return
	}

	userID, ok := userIDValue.(uint)
	if !ok {
		common.ServerError(c)
		return
	}

	var req RecordPlayHistoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ParamError(c, "请求参数错误")
		return
	}

	if err := getMusicService().RecordPlay(userID, req.MusicID); err != nil {
		common.ServerError(c)
		return
	}

	common.Success(c, map[string]interface{}{"message": "播放历史记录成功"})
}

// CleanDeletedFiles 清理已删除的文件
func CleanDeletedFiles(c *gin.Context) {
	err := getMusicService().CleanDeletedFiles()
	if err != nil {
		common.ServerError(c)
		return
	}

	common.Success(c, map[string]interface{}{
		"message": "清理完成",
	})
}

// StreamMusic 音频流式播放
func StreamMusic(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		common.ParamError(c, "music ID 参数错误")
		return
	}

	music, err := getMusicService().GetMusicByID(uint(id))
	if err != nil {
		common.ServerError(c)
		return
	}
	if music == nil {
		common.Fail(c, common.CodeNotFound, "音乐不存在")
		return
	}

	filePath := music.FilePath
	file, err := os.Open(filePath)
	if err != nil {
		common.Logger.Error("打开音乐文件失败", zap.String("file_path", filePath), zap.Error(err))
		common.ServerError(c)
		return
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		common.Logger.Error("获取音乐文件信息失败", zap.String("file_path", filePath), zap.Error(err))
		common.ServerError(c)
		return
	}

	// 设置适当的响应头以支持音频流
	c.Header("Content-Type", "audio/mpeg") // 默认设置为audio/mpeg
	c.Header("Accept-Ranges", "bytes")
	c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Header("Pragma", "no-cache")
	c.Header("Expires", "0")

	// 使用gin的Stream方法来流式传输文件
	http.ServeContent(c.Writer, c.Request, filepath.Base(filePath), stat.ModTime(), file)
}