package filesystem

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/jericho-yu/http-client/httpClient"
)

type (
	FileManager struct {
		Err            error
		dstDir, srcDir string
		fileBytes      []byte
		fileSize       int64
		config         *FileManagerConfig
	}
	FileManagerConfigDriver = string
	FileManagerConfig       struct {
		Username  string
		Password  string
		AuthTitle string
		Driver    FileManagerConfigDriver
	}
)

var (
	FileManageApp FileManager
)

const (
	FileManagerConfigDriverLocal FileManagerConfigDriver = "LOCAL"
	FileManagerConfigDriverNexus FileManagerConfigDriver = "NEXUS"
	FileManagerConfigDriverOss   FileManagerConfigDriver = "OSS"
)

// NewByLocalFile 初始化：文件管理器（通过本地文件）
func (FileManager) NewByLocalFile(srcDir, dstDir string, config *FileManagerConfig) (*FileManager, error) {
	fs := FileSystemApp.NewByAbs(srcDir)
	if !fs.IsExist {
		return nil, errors.New("目标文件不存在")
	}

	fileBytes, err := fs.Read()
	if err != nil {
		return nil, err
	}

	return &FileManager{
		dstDir:    dstDir,
		srcDir:    srcDir,
		fileBytes: fileBytes,
		fileSize:  int64(len(fileBytes)),
		config:    config,
	}, nil
}

// NewByBytes 实例化：文件管理器（通过字节）
func (FileManager) NewByBytes(srcFileBytes []byte, dstDir string, config *FileManagerConfig) *FileManager {
	return &FileManager{
		dstDir:    dstDir,
		fileBytes: srcFileBytes,
		fileSize:  int64(len(srcFileBytes)),
		config:    config,
	}
}

// Upload 上传文件
func (r *FileManager) Upload(src string) (int64, error) {
	switch r.config.Driver {
	case FileManagerConfigDriverLocal:
		return r.uploadToLocal()
	case FileManagerConfigDriverNexus:
		return r.uploadToNexus()
	case FileManagerConfigDriverOss:
		return r.uploadToOss()
	}

	return 0, fmt.Errorf("不支持的驱动类型：%s", r.config.Driver)
}

// 上传到本地
func (r *FileManager) uploadToLocal() (int64, error) {
	dst := FileSystemApp.NewByAbs(r.dstDir)
	return dst.WriteBytes(r.fileBytes)
}

// 上传到nexus
func (r *FileManager) uploadToNexus() (int64, error) {
	client := httpClient.New(r.dstDir).
		SetMethod(http.MethodPut).
		SetAuthorization(r.config.Username, r.config.Password, r.config.AuthTitle).
		AddHeaders(map[string][]string{
			"Content-Length": {fmt.Sprintf("%d", r.fileSize)},
		}).
		SetBody(r.fileBytes).
		Send()

	if client.Err != nil {
		return 0, client.Err
	}

	return int64(len(r.fileBytes)), nil
}

// 上传到oss
func (r *FileManager) uploadToOss() (int64, error) {
	return 0, errors.New("暂不支持oss方式")
}
