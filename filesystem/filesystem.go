package filesystem

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

type (
	// FileSystem 文件系统
	FileSystem struct {
		dir     string
		IsExist bool
		IsDir   bool
		IsFile  bool
	}

	// FileSystemCopyFileTarget 拷贝文件目标
	FileSystemCopyFilesTarget struct {
		Src         *FileSystem
		DstFilename string
	}
)

// NewFileSystemByRelative 实例化：文件系统（相对路径）
func NewFileSystemByRelative(dir string) *FileSystem {
	ins := &FileSystem{dir: filepath.Clean(filepath.Join(FileSystem{}.GetRootPath(), dir))}
	return ins.init()
}

// NewFileSystemByAbs 实例化：文件系统（绝对路径）
func NewFileSystemByAbs(dir string) *FileSystem {
	ins := &FileSystem{dir: dir}
	return ins.init()
}

// SetDirByRelative 设置路径：相对路径
func (r *FileSystem) SetDirByRelative(dir string) *FileSystem {
	r.dir = filepath.Clean(filepath.Join(FileSystem{}.GetRootPath(), dir))

	r.init()
	return r
}

// SetDir 设置路径：绝对路径
func (r *FileSystem) SetDirByAbs(dir string) *FileSystem {
	r.dir = dir

	r.init()
	return r
}

func (r *FileSystem) Join(dir string) *FileSystem {
	r.dir = filepath.Join(r.dir, dir)

	r.init()
	return r
}

// Joins 增加若干路径
func (r *FileSystem) Joins(dir ...string) *FileSystem {
	for _, v := range dir {
		r.Join(v)
	}

	r.init()
	return r
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
func (r *FileSystem) init() *FileSystem {
	var e error
	r.IsExist, e = r.Exist() // 检查文件是否存在
	if e != nil {
		panic(fmt.Errorf("检查路径错误：%s", e.Error()))
	}
	if r.IsExist {
		e = r.CheckPathType() // 检查路径类型
		if e != nil {
			panic(fmt.Errorf("检查路径类型错误：%s", e.Error()))
		}
	}
	return r
}

// Exist 检查文件是否存在
func (r *FileSystem) Exist() (bool, error) {
	_, err := os.Stat(r.dir)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// MkDir 创建文件夹
func (r *FileSystem) MkDir() error {
	if !r.IsExist {
		if e := os.MkdirAll(r.dir, os.ModePerm); e != nil {
			return e
		}
	}

	return nil
}

// GetDir 获取当前路径
func (r *FileSystem) GetDir() string {
	return r.dir
}

// CheckPathType 判断一个路径是文件还是文件夹
func (r *FileSystem) CheckPathType() error {
	info, e := os.Stat(r.dir)
	if e != nil {
		return e
	}

	if info.IsDir() {
		r.IsDir = true
		r.IsFile = !r.IsDir
	} else {
		r.IsFile = true
		r.IsDir = !r.IsFile
	}

	return nil
}

// Delete 删除文件或文件夹
func (r *FileSystem) Delete() error {
	if r.IsExist {
		if r.IsDir {
			return r.DelDir()
		}
		if r.IsFile {
			return r.DelFile()
		}
	}
	return nil
}

// DelDir 删除文件夹
func (r *FileSystem) DelDir() error {
	err := os.RemoveAll(r.dir)
	if err != nil {
		return err
	}
	return nil
}

// DelFile 删除文件
func (r *FileSystem) DelFile() error {
	e := os.Remove("path_to_your_file")
	if e != nil {
		return e
	}
	return nil
}

// CopyFile 拷贝单文件
func (r *FileSystem) CopyFile(dstDir, dstFilename string, abs bool) error {
	var (
		err         error
		srcFile     *os.File
		srcFilename string
		dst         *FileSystem
	)

	// 如果是相对路径
	if !abs {
		dst = NewFileSystemByRelative(dstDir)
	} else {
		dst = NewFileSystemByAbs(dstDir)
	}
	// 创建目标文件夹
	if !dst.IsDir {
		dst.MkDir()
	}

	// 判断源是否是文件
	if !r.IsFile {
		return errors.New("源文件不存在")
	}

	// 打开源文件
	srcFile, err = os.Open(r.GetDir())
	if err != nil {
		return err
	}
	defer srcFile.Close()

	if dstFilename == "" {
		srcFilename = filepath.Base(r.GetDir())
		dst.Join(srcFilename)
	} else {
		dst.Join(dstFilename)
	}

	// 创建目标文件
	dstFile, err := os.Create(dst.GetDir())
	if err != nil {
		return err
	}
	defer dstFile.Close()

	fmt.Printf("拷贝文件：%s ==>  %s\n", r.GetDir(), dstDir)

	// 拷贝内容
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}

	// 确保所有内容都已写入磁盘
	err = dstFile.Sync()
	if err != nil {
		return err
	}

	return nil
}

// CopyFiles 拷贝多个文件
func (FileSystem) CopyFiles(srcFiles []*FileSystemCopyFilesTarget, dstDir string, abs bool) error {
	var (
		err error
		dst *FileSystem
	)

	if abs {
		dst = NewFileSystemByAbs(dstDir)
	} else {
		dst = NewFileSystemByRelative(dstDir)
	}

	if !dst.IsDir {
		dst.MkDir()
	}

	for _, srcFile := range srcFiles {
		// 获取源文件名
		srcFilename := filepath.Base(srcFile.Src.GetDir())

		// 拷贝文件
		if srcFile.DstFilename != "" {
			err = srcFile.Src.CopyFile(dst.GetDir(), srcFile.DstFilename, true)
		} else {
			err = srcFile.Src.CopyFile(dst.GetDir(), srcFilename, true)
		}
		if err != nil {
			return err
		}
	}

	return nil
}

// CopyDir 拷贝目录
func (r *FileSystem) CopyDir(dstDir string, abs bool) error {
	// 判断是否是目录
	if !r.IsDir {
		return errors.New("源目录不存在")
	}

	// 遍历源目录
	err := filepath.Walk(r.GetDir(), func(srcPath string, info os.FileInfo, err error) error {
		var (
			src         *FileSystem
			dst         *FileSystem
			srcFilename string
		)

		if abs {
			dst = NewFileSystemByAbs(dstDir)
		} else {
			dst = NewFileSystemByRelative(dstDir)
		}

		if !dst.IsDir {
			dst.MkDir()
		}

		if err != nil {
			return err
		}

		srcFilename = filepath.Base(srcPath)
		src = NewFileSystemByAbs(srcPath)

		if src.IsFile {
			return src.CopyFile(dst.GetDir(), srcFilename, true)
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

// WriteBytes 写入文件：bytes
func (r *FileSystem) WriteBytes(content []byte) (int64, error) {
	var written int
	// 打开文件
	file, err := os.OpenFile(r.GetDir(), os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	// 写入内容
	written, err = file.Write(content)
	if err != nil {
		return 0, err
	}

	return int64(written), nil
}

// WriteAppend 追加写入文件
func (r *FileSystem) WriteAppend(content []byte) (int64, error) {
	var written int
	// Open the file in append mode.
	file, e := os.OpenFile(r.GetDir(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if e != nil {
		return 0, e
	}
	defer file.Close()

	// 追加写入内容
	written, e = file.Write(content)
	if e != nil {
		return 0, e
	}

	return int64(written), nil
}

// WriteString 写入文件：string
func (r *FileSystem) WriteString(content string) (int64, error) {
	return r.WriteBytes([]byte(content))
}

// WriteIoReader 写入文件：IoReader
func (r *FileSystem) WriteIoReader(content io.Reader) (written int64, err error) {
	dst, err := os.Create(r.dir)
	if err != nil {
		return 0, err
	}
	defer dst.Close()

	return io.Copy(dst, content)
}
