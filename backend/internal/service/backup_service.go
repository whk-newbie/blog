package service

import (
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/whk-newbie/blog/internal/config"
)

// BackupService 备份服务
type BackupService interface {
	// CreateBackup 创建备份
	CreateBackup() (*BackupInfo, error)
	// ListBackups 获取备份列表
	ListBackups() ([]BackupInfo, error)
	// GetBackupPath 获取备份文件路径
	GetBackupPath(filename string) (string, error)
	// DeleteBackup 删除备份文件
	DeleteBackup(filename string) error
	// CleanupOldBackups 清理旧备份
	CleanupOldBackups(retentionCount int) error
}

// BackupInfo 备份信息
type BackupInfo struct {
	Filename    string    `json:"filename"`
	Size        int64     `json:"size"`
	CreatedAt   time.Time `json:"created_at"`
	DownloadURL string    `json:"download_url"`
}

// backupService 备份服务实现
type backupService struct {
	cfg        *config.Config
	backupDir  string
	dbHost     string
	dbPort     int
	dbUser     string
	dbPassword string
	dbName     string
}

// NewBackupService 创建备份服务
func NewBackupService(cfg *config.Config) BackupService {
	backupDir := "./backups"
	if cfg.Upload.Path != "" {
		// 使用uploads同级目录
		backupDir = filepath.Join(filepath.Dir(cfg.Upload.Path), "backups")
	}

	// 确保备份目录存在
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		log.Printf("Failed to create backup directory: %v", err)
	}

	return &backupService{
		cfg:        cfg,
		backupDir:  backupDir,
		dbHost:     cfg.Database.Host,
		dbPort:     cfg.Database.Port,
		dbUser:     cfg.Database.User,
		dbPassword: cfg.Database.Password,
		dbName:     cfg.Database.DBName,
	}
}

// CreateBackup 创建备份
func (s *backupService) CreateBackup() (*BackupInfo, error) {
	// 生成备份文件名
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("backup_%s.sql.gz", timestamp)
	backupPath := filepath.Join(s.backupDir, filename)

	// 创建临时SQL文件
	tempSQLFile := filepath.Join(s.backupDir, fmt.Sprintf("backup_%s.sql", timestamp))
	defer os.Remove(tempSQLFile) // 清理临时文件

	// 构建pg_dump命令
	// 使用环境变量传递密码，避免在命令行中暴露
	cmd := exec.Command("pg_dump",
		"-h", s.dbHost,
		"-p", fmt.Sprintf("%d", s.dbPort),
		"-U", s.dbUser,
		"-d", s.dbName,
		"-F", "plain", // 使用纯文本格式
		"-f", tempSQLFile,
	)

	// 设置环境变量
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("PGPASSWORD=%s", s.dbPassword),
	)

	// 执行备份命令
	log.Printf("Creating database backup: %s", filename)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Failed to create backup: %v, output: %s", err, string(output))
		return nil, fmt.Errorf("failed to create backup: %w", err)
	}

	// 压缩SQL文件
	if err := s.compressFile(tempSQLFile, backupPath); err != nil {
		return nil, fmt.Errorf("failed to compress backup: %w", err)
	}

	// 获取文件信息
	fileInfo, err := os.Stat(backupPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get backup file info: %w", err)
	}

	backupInfo := &BackupInfo{
		Filename:    filename,
		Size:        fileInfo.Size(),
		CreatedAt:   fileInfo.ModTime(),
		DownloadURL: fmt.Sprintf("/api/v1/admin/backups/download/%s", filename),
	}

	log.Printf("Backup created successfully: %s (size: %d bytes)", filename, fileInfo.Size())
	return backupInfo, nil
}

// compressFile 压缩文件
func (s *backupService) compressFile(srcPath, dstPath string) error {
	// 打开源文件
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer srcFile.Close()

	// 创建目标文件
	dstFile, err := os.Create(dstPath)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dstFile.Close()

	// 创建gzip writer
	gzipWriter := gzip.NewWriter(dstFile)
	defer gzipWriter.Close()

	// 复制数据
	_, err = io.Copy(gzipWriter, srcFile)
	if err != nil {
		return fmt.Errorf("failed to compress file: %w", err)
	}

	return nil
}

// ListBackups 获取备份列表
func (s *backupService) ListBackups() ([]BackupInfo, error) {
	// 读取备份目录
	entries, err := os.ReadDir(s.backupDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []BackupInfo{}, nil
		}
		return nil, fmt.Errorf("failed to read backup directory: %w", err)
	}

	var backups []BackupInfo
	for _, entry := range entries {
		// 只处理.gz文件
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".sql.gz") {
			filePath := filepath.Join(s.backupDir, entry.Name())
			fileInfo, err := os.Stat(filePath)
			if err != nil {
				log.Printf("Failed to get file info for %s: %v", entry.Name(), err)
				continue
			}

			backupInfo := BackupInfo{
				Filename:    entry.Name(),
				Size:        fileInfo.Size(),
				CreatedAt:   fileInfo.ModTime(),
				DownloadURL: fmt.Sprintf("/api/v1/admin/backups/download/%s", entry.Name()),
			}
			backups = append(backups, backupInfo)
		}
	}

	// 按创建时间倒序排序（最新的在前）
	sort.Slice(backups, func(i, j int) bool {
		return backups[i].CreatedAt.After(backups[j].CreatedAt)
	})

	return backups, nil
}

// GetBackupPath 获取备份文件路径
func (s *backupService) GetBackupPath(filename string) (string, error) {
	// 安全检查：防止路径遍历攻击
	if strings.Contains(filename, "..") || strings.Contains(filename, "/") || strings.Contains(filename, "\\") {
		return "", fmt.Errorf("invalid filename")
	}

	// 只允许.sql.gz文件
	if !strings.HasSuffix(filename, ".sql.gz") {
		return "", fmt.Errorf("invalid file type")
	}

	backupPath := filepath.Join(s.backupDir, filename)

	// 检查文件是否存在
	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		return "", fmt.Errorf("backup file not found")
	}

	return backupPath, nil
}

// DeleteBackup 删除备份文件
func (s *backupService) DeleteBackup(filename string) error {
	backupPath, err := s.GetBackupPath(filename)
	if err != nil {
		return err
	}

	if err := os.Remove(backupPath); err != nil {
		return fmt.Errorf("failed to delete backup file: %w", err)
	}

	log.Printf("Backup file deleted: %s", filename)
	return nil
}

// CleanupOldBackups 清理旧备份
func (s *backupService) CleanupOldBackups(retentionCount int) error {
	backups, err := s.ListBackups()
	if err != nil {
		return err
	}

	// 如果备份数量不超过保留数量，不需要清理
	if len(backups) <= retentionCount {
		return nil
	}

	// 删除超出保留数量的旧备份
	deletedCount := 0
	for i := retentionCount; i < len(backups); i++ {
		if err := s.DeleteBackup(backups[i].Filename); err != nil {
			log.Printf("Failed to delete old backup %s: %v", backups[i].Filename, err)
			continue
		}
		deletedCount++
	}

	log.Printf("Cleaned up %d old backup files", deletedCount)
	return nil
}
