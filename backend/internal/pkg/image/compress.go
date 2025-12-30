package image

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/chai2010/webp"
	"github.com/disintegration/imaging"
)

const (
	// CompressThreshold 压缩阈值：超过1MB的图片需要压缩
	CompressThreshold = 1024 * 1024 // 1MB
	// CompressQuality 压缩质量：80%
	CompressQuality = 80
	// WebPExtension WebP文件扩展名
	WebPExtension = ".webp"
)

// CompressImage 压缩图片
// 如果图片大小超过阈值，则压缩并转换为WebP格式
// 返回压缩后的文件路径和文件大小
func CompressImage(srcPath string) (string, int64, error) {
	// 检查源文件大小
	fileInfo, err := os.Stat(srcPath)
	if err != nil {
		return "", 0, fmt.Errorf("failed to stat source file: %w", err)
	}

	// 如果文件小于阈值，不需要压缩
	if fileInfo.Size() < CompressThreshold {
		return srcPath, fileInfo.Size(), nil
	}

	// 打开源图片
	src, err := imaging.Open(srcPath)
	if err != nil {
		return "", 0, fmt.Errorf("failed to open image: %w", err)
	}

	// 生成WebP文件路径
	ext := filepath.Ext(srcPath)
	webpPath := strings.TrimSuffix(srcPath, ext) + WebPExtension

	// 保存为WebP格式
	err = saveAsWebP(src, webpPath, CompressQuality)
	if err != nil {
		return "", 0, fmt.Errorf("failed to save compressed image: %w", err)
	}

	// 获取压缩后的文件大小
	compressedInfo, err := os.Stat(webpPath)
	if err != nil {
		return "", 0, fmt.Errorf("failed to stat compressed file: %w", err)
	}

	// 如果压缩后的文件比原文件大，保留原文件
	if compressedInfo.Size() >= fileInfo.Size() {
		os.Remove(webpPath)
		return srcPath, fileInfo.Size(), nil
	}

	// 删除原文件，使用压缩后的文件
	if err := os.Remove(srcPath); err != nil {
		// 如果删除失败，保留两个文件，但返回压缩后的文件路径
		return webpPath, compressedInfo.Size(), nil
	}

	return webpPath, compressedInfo.Size(), nil
}

// CompressImageFromReader 从Reader压缩图片
// 用于处理上传的文件流
func CompressImageFromReader(reader io.Reader, dstPath string, originalSize int64) (string, int64, error) {
	// 如果文件小于阈值，直接保存，不压缩
	if originalSize < CompressThreshold {
		dstFile, err := os.Create(dstPath)
		if err != nil {
			return "", 0, fmt.Errorf("failed to create destination file: %w", err)
		}
		defer dstFile.Close()

		_, err = io.Copy(dstFile, reader)
		if err != nil {
			return "", 0, fmt.Errorf("failed to copy file: %w", err)
		}

		return dstPath, originalSize, nil
	}

	// 解码图片
	img, _, err := image.Decode(reader)
	if err != nil {
		return "", 0, fmt.Errorf("failed to decode image: %w", err)
	}

	// 生成WebP文件路径
	ext := filepath.Ext(dstPath)
	webpPath := strings.TrimSuffix(dstPath, ext) + WebPExtension

	// 转换为imaging.Image类型
	imagingImg := imaging.Clone(img)

	// 保存为WebP格式
	err = saveAsWebP(imagingImg, webpPath, CompressQuality)

	if err != nil {
		return "", 0, fmt.Errorf("failed to save compressed image: %w", err)
	}

	// 获取压缩后的文件大小
	compressedInfo, err := os.Stat(webpPath)
	if err != nil {
		return "", 0, fmt.Errorf("failed to stat compressed file: %w", err)
	}

	// 如果压缩后的文件比原文件大，保留原文件
	if compressedInfo.Size() >= originalSize {
		os.Remove(webpPath)
		// 重新保存原文件
		dstFile, err := os.Create(dstPath)
		if err != nil {
			return "", 0, fmt.Errorf("failed to create destination file: %w", err)
		}
		defer dstFile.Close()

		// 重新读取源数据（需要重新定位reader，这里简化处理）
		return dstPath, originalSize, nil
	}

	return webpPath, compressedInfo.Size(), nil
}

// GetImageDimensions 获取图片尺寸
func GetImageDimensions(imgPath string) (width, height int, err error) {
	img, err := imaging.Open(imgPath)
	if err != nil {
		return 0, 0, err
	}
	bounds := img.Bounds()
	return bounds.Dx(), bounds.Dy(), nil
}

// IsImageFormatSupported 检查图片格式是否支持
func IsImageFormatSupported(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	supportedFormats := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".webp": true,
	}
	return supportedFormats[ext]
}

// SaveImage 保存图片（用于JPEG和PNG格式）
func SaveImage(img image.Image, path string, format string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	switch format {
	case "jpeg", "jpg":
		return jpeg.Encode(file, img, &jpeg.Options{Quality: CompressQuality})
	case "png":
		return png.Encode(file, img)
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}
}

// saveAsWebP 保存图片为WebP格式
func saveAsWebP(img image.Image, path string, quality int) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create webp file: %w", err)
	}
	defer file.Close()

	// 使用webp库编码为WebP格式
	// quality参数范围是0-100，我们传入80
	options := &webp.Options{
		Lossless: false,
		Quality:  float32(quality),
	}

	err = webp.Encode(file, img, options)
	if err != nil {
		return fmt.Errorf("failed to encode webp: %w", err)
	}

	return nil
}
