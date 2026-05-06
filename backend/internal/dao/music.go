package dao

import (
	"time"

	"github.com/homemusic/backend/internal/common"
	"github.com/homemusic/backend/internal/model"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// MusicDao 音乐数据访问对象
type MusicDao struct {
	db *gorm.DB
}

// NewMusicDao 创建音乐DAO实例
func NewMusicDao() *MusicDao {
	return &MusicDao{db: DB}
}

// CreateMusic 创建音乐记录
func (d *MusicDao) CreateMusic(music *model.Music) error {
	err := d.db.Create(music).Error
	if err != nil {
		common.Logger.Error("创建音乐记录失败", zap.String("file_path", music.FilePath), zap.Error(err))
		return err
	}
	common.Logger.Debug("创建音乐记录成功", zap.Uint("music_id", music.ID), zap.String("title", music.Title))
	return nil
}

// UpdateMusic 更新音乐记录
func (d *MusicDao) UpdateMusic(music *model.Music) error {
	err := d.db.Save(music).Error
	if err != nil {
		common.Logger.Error("更新音乐记录失败", zap.Uint("music_id", music.ID), zap.Error(err))
		return err
	}
	common.Logger.Debug("更新音乐记录成功", zap.Uint("music_id", music.ID))
	return nil
}

// GetMusicByPath 根据文件路径获取音乐
func (d *MusicDao) GetMusicByPath(filePath string) (*model.Music, error) {
	var music model.Music
	err := d.db.Where("file_path = ?", filePath).First(&music).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		common.Logger.Error("根据路径查询音乐失败", zap.String("file_path", filePath), zap.Error(err))
		return nil, err
	}
	return &music, nil
}

// GetMusicByID 根据ID获取音乐
func (d *MusicDao) GetMusicByID(id uint) (*model.Music, error) {
	var music model.Music
	err := d.db.First(&music, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		common.Logger.Error("根据ID查询音乐失败", zap.Uint("music_id", id), zap.Error(err))
		return nil, err
	}
	return &music, nil
}

// GetMusicsByStatus 根据扫描状态获取音乐列表
func (d *MusicDao) GetMusicsByStatus(status string, limit, offset int) ([]*model.Music, error) {
	var musics []*model.Music
	query := d.db.Where("scan_status = ?", status)
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.Find(&musics).Error
	if err != nil {
		common.Logger.Error("根据状态查询音乐列表失败", zap.String("status", status), zap.Error(err))
		return nil, err
	}
	return musics, nil
}

// GetAllMusics 获取所有音乐（分页）
func (d *MusicDao) GetAllMusics(limit, offset int) ([]*model.Music, int64, error) {
	var musics []*model.Music
	var total int64

	// 获取总数
	err := d.db.Model(&model.Music{}).Count(&total).Error
	if err != nil {
		common.Logger.Error("获取音乐总数失败", zap.Error(err))
		return nil, 0, err
	}

	// 获取分页数据
	query := d.db.Order("created_at DESC")
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	err = query.Find(&musics).Error
	if err != nil {
		common.Logger.Error("获取音乐列表失败", zap.Error(err))
		return nil, 0, err
	}

	return musics, total, nil
}

// GetArtists 获取歌手列表
func (d *MusicDao) GetArtists(limit, offset int) ([]map[string]interface{}, int64, error) {
	type artistCount struct {
		Artist string `json:"name"`
		Count  int64  `json:"count"`
	}

	var artists []artistCount
	var total int64

	countQuery := d.db.Model(&model.Music{}).
		Select("artist").
		Where("artist <> ''").
		Group("artist")
	if err := countQuery.Count(&total).Error; err != nil {
		common.Logger.Error("统计歌手数量失败", zap.Error(err))
		return nil, 0, err
	}

	query := d.db.Model(&model.Music{}).
		Select("artist as artist, count(*) as count").
		Where("artist <> ''").
		Group("artist").
		Order("count DESC, artist ASC")
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}
	if err := query.Scan(&artists).Error; err != nil {
		common.Logger.Error("获取歌手列表失败", zap.Error(err))
		return nil, 0, err
	}

	result := make([]map[string]interface{}, 0, len(artists))
	for _, item := range artists {
		result = append(result, map[string]interface{}{
			"name":  item.Artist,
			"count": item.Count,
		})
	}

	return result, total, nil
}

// GetAlbums 获取专辑列表
func (d *MusicDao) GetAlbums(limit, offset int) ([]map[string]interface{}, int64, error) {
	type albumCount struct {
		Album string `json:"name"`
		Count int64  `json:"count"`
	}

	var albums []albumCount
	var total int64

	countQuery := d.db.Model(&model.Music{}).
		Select("album").
		Where("album <> ''").
		Group("album")
	if err := countQuery.Count(&total).Error; err != nil {
		common.Logger.Error("统计专辑数量失败", zap.Error(err))
		return nil, 0, err
	}

	query := d.db.Model(&model.Music{}).
		Select("album as album, count(*) as count").
		Where("album <> ''").
		Group("album").
		Order("count DESC, album ASC")
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}
	if err := query.Scan(&albums).Error; err != nil {
		common.Logger.Error("获取专辑列表失败", zap.Error(err))
		return nil, 0, err
	}

	result := make([]map[string]interface{}, 0, len(albums))
	for _, item := range albums {
		result = append(result, map[string]interface{}{
			"name":  item.Album,
			"count": item.Count,
		})
	}

	return result, total, nil
}

// SearchMusics 全局搜索音乐
func (d *MusicDao) SearchMusics(keyword string, limit, offset int) ([]*model.Music, int64, error) {
	var musics []*model.Music
	var total int64
	pattern := "%" + keyword + "%"

	countQuery := d.db.Model(&model.Music{}).
		Where("title LIKE ? OR artist LIKE ? OR album LIKE ?", pattern, pattern, pattern)
	if err := countQuery.Count(&total).Error; err != nil {
		common.Logger.Error("统计搜索结果失败", zap.String("keyword", keyword), zap.Error(err))
		return nil, 0, err
	}

	query := d.db.Model(&model.Music{}).
		Where("title LIKE ? OR artist LIKE ? OR album LIKE ?", pattern, pattern, pattern).
		Order("created_at DESC")
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}
	if err := query.Find(&musics).Error; err != nil {
		common.Logger.Error("搜索音乐失败", zap.String("keyword", keyword), zap.Error(err))
		return nil, 0, err
	}

	return musics, total, nil
}

// RecordPlayHistory 记录播放历史
func (d *MusicDao) RecordPlayHistory(history *model.PlayHistory) error {
	err := d.db.Create(history).Error
	if err != nil {
		common.Logger.Error("记录播放历史失败", zap.Error(err))
		return err
	}
	common.Logger.Debug("记录播放历史成功", zap.Uint("user_id", history.UserID), zap.Uint("music_id", history.MusicID))
	return nil
}

// GetRecentPlayHistory 获取最近播放记录
func (d *MusicDao) GetRecentPlayHistory(userID uint, limit, offset int) ([]*model.PlayHistory, int64, error) {
	var records []*model.PlayHistory
	var total int64

	countQuery := d.db.Model(&model.PlayHistory{}).
		Where("user_id = ?", userID)
	if err := countQuery.Count(&total).Error; err != nil {
		common.Logger.Error("统计最近播放记录失败", zap.Uint("user_id", userID), zap.Error(err))
		return nil, 0, err
	}

	query := d.db.Model(&model.PlayHistory{}).
		Where("user_id = ?", userID).
		Preload("Music").
		Order("played_at DESC")
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}
	if err := query.Find(&records).Error; err != nil {
		common.Logger.Error("获取最近播放记录失败", zap.Uint("user_id", userID), zap.Error(err))
		return nil, 0, err
	}

	return records, total, nil
}

// DeleteMusicByPath 根据路径删除音乐记录
func (d *MusicDao) DeleteMusicByPath(filePath string) error {
	err := d.db.Where("file_path = ?", filePath).Delete(&model.Music{}).Error
	if err != nil {
		common.Logger.Error("删除音乐记录失败", zap.String("file_path", filePath), zap.Error(err))
		return err
	}
	common.Logger.Info("删除音乐记录成功", zap.String("file_path", filePath))
	return nil
}

// UpdateScanStatus 更新扫描状态
func (d *MusicDao) UpdateScanStatus(filePath, status string, lastScanned time.Time, errorMsg string) error {
	updateData := map[string]interface{}{
		"scan_status":  status,
		"last_scanned": lastScanned,
	}
	if errorMsg != "" {
		updateData["error_msg"] = errorMsg
	}

	err := d.db.Model(&model.Music{}).Where("file_path = ?", filePath).Updates(updateData).Error
	if err != nil {
		common.Logger.Error("更新扫描状态失败", zap.String("file_path", filePath), zap.String("status", status), zap.Error(err))
		return err
	}
	return nil
}

// GetScanStatistics 获取扫描统计信息
func (d *MusicDao) GetScanStatistics() (map[string]int64, error) {
	var stats []struct {
		Status string `json:"status"`
		Count  int64  `json:"count"`
	}

	err := d.db.Model(&model.Music{}).
		Select("scan_status as status, count(*) as count").
		Group("scan_status").
		Scan(&stats).Error
	if err != nil {
		common.Logger.Error("获取扫描统计失败", zap.Error(err))
		return nil, err
	}

	result := make(map[string]int64)
	for _, stat := range stats {
		result[stat.Status] = stat.Count
	}

	return result, nil
}
