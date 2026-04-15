package files

import (
	"bytes"
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/zhuofanxu/axb/errx"
)

func PathExists(path string) (bool, error) {
	fi, err := os.Stat(path)

	if err == nil {
		if fi.IsDir() {
			return true, nil
		}
		// 路径存在但不是目录，返回 true 并附带错误，让调用方明确知道路径已被占用
		return true, errors.New("path exists but is not a directory")
	} else if errors.Is(err, os.ErrNotExist) {
		return false, nil
	} else if errors.Is(err, os.ErrPermission) {
		return false, errors.New("permission denied")
	} else {
		return false, err
	}
}

// EnsureDir 确保目录存在，如果不存在则创建
func EnsureDir(dir string) error {
	if exists, err := PathExists(dir); err != nil {
		return err
	} else if !exists {
		return os.MkdirAll(dir, 0755)
	}
	return nil
}

// GenerateUniqueFilename 生成唯一的文件名（不包含原始文件名信息）
func GenerateUniqueFilename(originalFilename string) string {
	ext := filepath.Ext(originalFilename)
	timestamp := time.Now().Format("20060102150405")
	uuid := uuid.New().String()[:8]
	return fmt.Sprintf("%s_%s%s", timestamp, uuid, ext)
}

// SaveUploadedFile 保存上传的文件到指定目录
func SaveUploadedFile(file *multipart.FileHeader, uploadDir string) (string, error) {
	// 确保上传目录存在
	if err := EnsureDir(uploadDir); err != nil {
		return "", fmt.Errorf("failed to create upload directory: %w", err)
	}

	// 生成唯一文件名
	filename := GenerateUniqueFilename(file.Filename)
	filePath := filepath.Join(uploadDir, filename)

	// 打开上传的文件
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer src.Close()

	// 创建目标文件
	dst, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dst.Close()

	// 复制文件内容，失败时清理已创建的目标文件
	if _, err := io.Copy(dst, src); err != nil {
		_ = os.Remove(filePath)
		return "", fmt.Errorf("failed to copy file content: %w", err)
	}

	return filePath, nil
}

// DeleteFile 删除指定路径的文件
func DeleteFile(filePath string) error {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil // 文件不存在，认为删除成功
	}
	return os.Remove(filePath)
}

// GetFileMD5 计算文件的MD5哈希值
func GetFileMD5(file *multipart.FileHeader) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, src); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

func CloseFile(f io.Closer, operation string, logger *zap.Logger) {
	if err := f.Close(); err != nil {
		if logger != nil {
			logger.Error("failed to close file", zap.String("operation", operation), zap.Error(err))
		}
	}
}

func ValidateExcelFormat(filename string, src io.ReadSeeker) error {
	// 检查文件扩展名
	ext := strings.ToLower(filepath.Ext(filename))
	if ext == ".xls" {
		return errx.NewError(errx.CodeParamError, nil).WithMsg("不支持.xls格式，请将文件另存为.xlsx格式后重新上传")
	}

	if ext != ".xlsx" && ext != ".xlsm" {
		return errx.NewError(errx.CodeParamError, nil).WithMsg("不支持的文件格式，请上传Excel文件(.xlsx或.xlsm格式)")
	}

	// 读取文件头部字节以检测格式
	buf := make([]byte, 8)
	if _, err := src.Read(buf); err != nil {
		return errx.NewError(errx.CodeParamError, err).WithMsg("无法读取文件内容，请确认文件完整")
	}

	// 重置文件指针
	if _, err := src.Seek(0, io.SeekStart); err != nil {
		return errx.NewError(errx.CodeParamError, err).WithMsg("文件读取失败，请重新上传")
	}

	// 检查是否为老格式xls（OLE2格式）
	if bytes.HasPrefix(buf, []byte{0xD0, 0xCF, 0x11, 0xE0}) {
		return errx.NewError(errx.CodeParamError, nil).WithMsg("检测到.xls格式文件，请将文件另存为.xlsx格式后重新上传")
	}

	return nil
}
