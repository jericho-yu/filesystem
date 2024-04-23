package filesystem

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

// FileSystem 文件系统
type FileSystem struct {
	dir     string
	IsExist bool
	IsDir   bool
	IsFile  bool
}

// NewFileSystemByRelative 实例化：文件系统（相对路径）
func NewFileSystemByRelative(dir string) (*FileSystem, error) {
	ins := &FileSystem{dir: filepath.Clean(filepath.Join(FileSystem{}.GetRootPath(), dir))}
	return ins.init()
}

// NewFileSystemByAbs 实例化：文件系统（绝对路径）
func NewFileSystemByAbs(dir string) (*FileSystem, error) {
	ins := &FileSystem{dir: dir}
	return ins.init()
}

func (FileSystem) GetRootPath() string {
	rootPath, _ := filepath.Abs(".")
	return rootPath
}

// GetCurrentPath 最终方案-全兼容
func (FileSystem) GetCurrentPath(paths ...string) string {
	dir := getGoBuildPath()
	if strings.Contains(dir, getTmpDir()) {
		return getGoRunPath()
	}
	return dir
}

// 获取系统临时目录，兼容go run
func getTmpDir() string {
	dir := os.Getenv("TEMP")
	if dir == "" {
		dir = os.Getenv("TMP")
	}
	res, _ := filepath.EvalSymlinks(dir)
	return res
}

// 获取当前执行文件绝对路径
func getGoBuildPath() string {
	exePath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	res, _ := filepath.EvalSymlinks(filepath.Dir(exePath))
	return res
}

// 获取当前执行文件绝对路径（go run）
func getGoRunPath() string {
	var abPath string
	_, filename, _, ok := runtime.Caller(0)
	if ok {
		abPath = path.Dir(filename)
	}
	return abPath
}

// 初始化
func (receiver *FileSystem) init() (*FileSystem, error) {
	var e error
	receiver.IsExist, e = receiver.Exist() // 检查文件是否存在
	if e != nil {
		panic(fmt.Errorf("检查路径错误：%s", e.Error()))
	}
	if receiver.IsExist {
		e = receiver.CheckPathType() // 检查路径类型
		if e != nil {
			panic(fmt.Errorf("检查路径类型错误：%s", e.Error()))
		}
	}
	return receiver, nil
}

// Exist 检查文件是否存在
func (receiver *FileSystem) Exist() (bool, error) {
	_, err := os.Stat(receiver.dir)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// MkDir 创建文件夹
func (receiver *FileSystem) MkDir() error {
	if !receiver.IsExist {
		if e := os.MkdirAll(receiver.dir, os.ModePerm); e != nil {
			return e
		}
	}

	return nil
}

// GetDir 获取当前路径
func (receiver *FileSystem) GetDir() string {
	return receiver.dir
}

// CheckPathType 判断一个路径是文件还是文件夹
func (receiver *FileSystem) CheckPathType() error {
	info, e := os.Stat(receiver.dir)
	if e != nil {
		return e
	}

	if info.IsDir() {
		receiver.IsDir = true
		receiver.IsFile = !receiver.IsDir
	} else {
		receiver.IsFile = true
		receiver.IsDir = !receiver.IsFile
	}

	return nil
}

// Delete 删除文件或文件夹
func (receiver *FileSystem) Delete() error {
	if receiver.IsExist {
		if receiver.IsDir {
			return receiver.DelDir()
		}
		if receiver.IsFile {
			return receiver.DelFile()
		}
	}
	return nil
}

// DelDir 删除文件夹
func (receiver *FileSystem) DelDir() error {
	err := os.RemoveAll(receiver.dir)
	if err != nil {
		return err
	}
	return nil
}

// DelFile 删除文件
func (receiver *FileSystem) DelFile() error {
	e := os.Remove("path_to_your_file")
	if e != nil {
		return e
	}
	return nil
}
