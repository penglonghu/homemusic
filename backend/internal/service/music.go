package service

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/homemusic/backend/internal/dao"
	"github.com/homemusic/backend/internal/model"
	"go.uber.org/zap"
)

// MusicService 音乐扫描服务
type MusicService struct {
	musicDao *dao.MusicDao
}

// NewMusicService 创建音乐服务实例
func NewMusicService() *MusicService {
	return &MusicService{
		musicDao: dao.NewMusicDao(),
	}
}

// 支持的音频文件扩展名
var supportedExtensions = map[string]bool{
	".mp3":  true,
	".flac": true,
	".m4a":  true,
	".aac":  true,
	".ogg":  true,
	".wma":  true,
	".wav":  true,
}

// isSupportedAudioFile 检查是否为支持的音频文件
func isSupportedAudioFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	return supportedExtensions[ext]
}

// ScanMusicDirectory 扫描音乐目录
func (s *MusicService) ScanMusicDirectory(musicPath string) error {
	zap.L().Info("开始扫描音乐目录", zap.String("path", musicPath))

	// 获取所有音频文件
	var audioFiles []string
	err := filepath.Walk(musicPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			zap.L().Warn("访问文件失败", zap.String("path", path), zap.Error(err))
			return nil // 继续扫描
		}
		if !info.IsDir() && isSupportedAudioFile(info.Name()) {
			audioFiles = append(audioFiles, path)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("扫描目录失败: %v", err)
	}

	zap.L().Info("发现音频文件", zap.Int("count", len(audioFiles)))

	// 并发处理文件
	const maxWorkers = 10
	fileChan := make(chan string, len(audioFiles))
	var wg sync.WaitGroup

	// 启动工作协程
	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for filePath := range fileChan {
				s.processAudioFile(filePath)
			}
		}()
	}

	// 发送文件到通道
	for _, file := range audioFiles {
		fileChan <- file
	}
	close(fileChan)

	wg.Wait()

	zap.L().Info("音乐目录扫描完成")
	return nil
}

// processAudioFile 处理单个音频文件
func (s *MusicService) processAudioFile(filePath string) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		zap.L().Error("获取文件信息失败", zap.String("path", filePath), zap.Error(err))
		s.updateScanStatus(filePath, "error", fmt.Sprintf("获取文件信息失败: %v", err))
		return
	}

	// 检查是否已存在且未修改
	existing, err := s.musicDao.GetMusicByPath(filePath)
	if err != nil {
		zap.L().Error("查询现有记录失败", zap.String("path", filePath), zap.Error(err))
		return
	}

	if existing != nil && !fileInfo.ModTime().After(existing.ModTime) && existing.ScanStatus == "scanned" {
		// 文件未修改且已扫描，跳过
		return
	}

	// 解析元数据
	metadata, err := s.extractMetadata(filePath)
	if err != nil {
		zap.L().Warn("解析元数据失败", zap.String("path", filePath), zap.Error(err))
		s.updateScanStatus(filePath, "error", fmt.Sprintf("解析元数据失败: %v", err))
		return
	}

	// 创建或更新记录
	music := &model.Music{
		FilePath:    filePath,
		FileSize:    fileInfo.Size(),
		ModTime:     fileInfo.ModTime(),
		Title:       metadata.Title,
		Artist:      metadata.Artist,
		Album:       metadata.Album,
		Genre:       metadata.Genre,
		Year:        metadata.Year,
		Duration:    metadata.Duration,
		TrackNumber: metadata.TrackNumber,
		Bitrate:     metadata.Bitrate,
		ScanStatus:  "scanned",
		LastScanned: time.Now(),
	}

	if existing != nil {
		music.ID = existing.ID
		err = s.musicDao.UpdateMusic(music)
	} else {
		err = s.musicDao.CreateMusic(music)
	}

	if err != nil {
		zap.L().Error("保存音乐记录失败", zap.String("path", filePath), zap.Error(err))
		s.updateScanStatus(filePath, "error", fmt.Sprintf("保存记录失败: %v", err))
		return
	}

	zap.L().Debug("处理音频文件成功", zap.String("path", filePath), zap.String("title", metadata.Title))
}

// Metadata 音乐元数据
type Metadata struct {
	Title       string
	Artist      string
	Album       string
	Genre       string
	Year        int
	Duration    int
	TrackNumber int
	Bitrate     int
}

// extractMetadata 提取音频文件元数据（支持ID3v1标签解析）
func (s *MusicService) extractMetadata(filePath string) (*Metadata, error) {
	filename := filepath.Base(filePath)
	title := strings.TrimSuffix(filename, filepath.Ext(filename))

	// 尝试解析ID3v1标签
	metadata, err := s.parseID3v1(filePath)
	if err == nil && metadata.Title != "" {
		return metadata, nil
	}

	// 如果ID3v1解析失败，回退到文件名解析
	var artist string
	if strings.Contains(title, " - ") {
		parts := strings.SplitN(title, " - ", 2)
		artist = strings.TrimSpace(parts[0])
		title = strings.TrimSpace(parts[1])
	}

	return &Metadata{
		Title:       title,
		Artist:      artist,
		Album:       "",
		Genre:       "",
		Year:        0,
		Duration:    0,
		TrackNumber: 0,
		Bitrate:     128, // 默认比特率
	}, nil
}

// parseID3v1 解析ID3v1标签
func (s *MusicService) parseID3v1(filePath string) (*Metadata, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// 获取文件大小
	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}

	if stat.Size() < 128 {
		return nil, fmt.Errorf("file too small for ID3v1")
	}

	// 读取最后128字节
	_, err = file.Seek(-128, 2)
	if err != nil {
		return nil, err
	}

	buffer := make([]byte, 128)
	_, err = file.Read(buffer)
	if err != nil {
		return nil, err
	}

	// 检查TAG标识
	if string(buffer[0:3]) != "TAG" {
		return nil, fmt.Errorf("no ID3v1 tag found")
	}

	// 解析标签
	metadata := &Metadata{}

	// Title (30 bytes, offset 3)
	metadata.Title = strings.TrimRight(string(buffer[3:33]), "\x00 ")

	// Artist (30 bytes, offset 33)
	metadata.Artist = strings.TrimRight(string(buffer[33:63]), "\x00 ")

	// Album (30 bytes, offset 63)
	metadata.Album = strings.TrimRight(string(buffer[63:93]), "\x00 ")

	// Year (4 bytes, offset 93)
	yearStr := strings.TrimRight(string(buffer[93:97]), "\x00 ")
	if yearStr != "" {
		if year, err := strconv.Atoi(yearStr); err == nil {
			metadata.Year = year
		}
	}

	// Comment (28 bytes, offset 97) - 注意：ID3v1.1在comment最后有track number
	metadata.Genre = strings.TrimRight(string(buffer[97:125]), "\x00 ")

	// Genre (1 byte, offset 127)
	genreID := buffer[127]
	if int(genreID) < len(genreList) {
		metadata.Genre = genreList[genreID]
	}

	// Track number (ID3v1.1, last byte of comment)
	if buffer[125] == 0 && buffer[126] != 0 {
		metadata.TrackNumber = int(buffer[126])
	}

	return metadata, nil
}

// ID3v1 genre list
var genreList = []string{
	"Blues", "Classic Rock", "Country", "Dance", "Disco", "Funk", "Grunge", "Hip-Hop",
	"Jazz", "Metal", "New Age", "Oldies", "Other", "Pop", "R&B", "Rap", "Reggae", "Rock",
	"Techno", "Industrial", "Alternative", "Ska", "Death Metal", "Pranks", "Soundtrack",
	"Euro-Techno", "Ambient", "Trip-Hop", "Vocal", "Jazz+Funk", "Fusion", "Trance",
	"Classical", "Instrumental", "Acid", "House", "Game", "Sound Clip", "Gospel", "Noise",
	"AlternRock", "Bass", "Soul", "Punk", "Space", "Meditative", "Instrumental Pop",
	"Instrumental Rock", "Ethnic", "Gothic", "Darkwave", "Techno-Industrial", "Electronic",
	"Pop-Folk", "Eurodance", "Dream", "Southern Rock", "Comedy", "Cult", "Gangsta",
	"Top 40", "Christian Rap", "Pop/Funk", "Jungle", "Native American", "Cabaret", "New Wave",
	"Psychadelic", "Rave", "Showtunes", "Trailer", "Lo-Fi", "Tribal", "Acid Punk", "Acid Jazz",
	"Polka", "Retro", "Musical", "Rock & Roll", "Hard Rock", "Folk", "Folk-Rock", "National Folk",
	"Swing", "Fast Fusion", "Bebob", "Latin", "Revival", "Celtic", "Bluegrass", "Avantgarde",
	"Gothic Rock", "Progressive Rock", "Psychedelic Rock", "Symphonic Rock", "Slow Rock",
	"Big Band", "Chorus", "Easy Listening", "Acoustic", "Humour", "Speech", "Chanson",
	"Opera", "Chamber Music", "Sonata", "Symphony", "Booty Bass", "Primus", "Porn Groove",
	"Satire", "Slow Jam", "Club", "Tango", "Samba", "Folklore", "Ballad", "Power Ballad",
	"Rhythmic Soul", "Freestyle", "Duet", "Punk Rock", "Drum Solo", "Acapella", "Euro-House",
	"Dance Hall", "Goa", "Drum & Bass", "Club-House", "Hardcore", "Terror", "Indie", "BritPop",
	"Negerpunk", "Polsk Punk", "Beat", "Christian Gangsta Rap", "Heavy Metal", "Black Metal",
	"Crossover", "Contemporary Christian", "Christian Rock", "Merengue", "Salsa", "Thrash Metal",
	"Anime", "JPop", "Synthpop",
}

// updateScanStatus 更新扫描状态
func (s *MusicService) updateScanStatus(filePath, status, errorMsg string) {
	err := s.musicDao.UpdateScanStatus(filePath, status, time.Now(), errorMsg)
	if err != nil {
		zap.L().Error("更新扫描状态失败", zap.String("path", filePath), zap.String("status", status), zap.Error(err))
	}
}

// GetScanStatistics 获取扫描统计
func (s *MusicService) GetScanStatistics() (map[string]interface{}, error) {
	stats, err := s.musicDao.GetScanStatistics()
	if err != nil {
		return nil, err
	}

	total := int64(0)
	for _, count := range stats {
		total += count
	}

	return map[string]interface{}{
		"total":     total,
		"scanned":   stats["scanned"],
		"pending":   stats["pending"],
		"error":     stats["error"],
	}, nil
}

// GetMusics 获取音乐列表（分页）
func (s *MusicService) GetMusics(page, pageSize int) ([]*model.Music, int64, error) {
	offset := (page - 1) * pageSize
	return s.musicDao.GetAllMusics(pageSize, offset)
}

// CleanDeletedFiles 清理已删除的文件记录
func (s *MusicService) CleanDeletedFiles() error {
	// 获取所有记录
	musics, _, err := s.musicDao.GetAllMusics(0, 0)
	if err != nil {
		return err
	}

	deletedCount := 0
	for _, music := range musics {
		if _, err := os.Stat(music.FilePath); os.IsNotExist(err) {
			err = s.musicDao.DeleteMusicByPath(music.FilePath)
			if err != nil {
				zap.L().Error("删除不存在的文件记录失败", zap.String("path", music.FilePath), zap.Error(err))
			} else {
				deletedCount++
			}
		}
	}

	zap.L().Info("清理已删除文件完成", zap.Int("deleted_count", deletedCount))
	return nil
}