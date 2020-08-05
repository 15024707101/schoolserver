package utils

import (
	"bufio"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"time"
)

// 检查文件或目录是否存在
// 如果由 path 指定的文件或目录存在则返回 true，否则返回 false

func IsExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

/*func IsExist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsPermission(err) {
		return false, errors.New(fmt.Sprintf("You have no permission to access %s", path))
	}
	if os.IsNotExist(err) {
		return false, nil
	}

	return false, errors.New("exist or not exist, that is a exception")

}*/

// 列出指定路径中的文件和目录
// 如果目录不存在，则返回空slice
func ScanDir(directory string) []string {
	file, err := os.Open(directory)
	if err != nil {
		return []string{}
	}
	names, err := file.Readdirnames(-1)
	if err != nil {
		return []string{}
	}
	return names
}

// 判断给定文件名是否是一个目录
// 如果文件名存在并且为目录则返回 true。如果 filename 是一个相对路径，则按照当前工作目录检查其相对路径。
func IsDir(filename string) bool {
	return isFileOrDir(filename, true)
}

// 判断给定文件名是否为一个正常的文件
// 如果文件存在且为正常的文件则返回 true
func IsFile(filename string) bool {
	return isFileOrDir(filename, false)
}

// 判断是文件还是目录，根据decideDir为true表示判断是否为目录；否则判断是否为文件
func isFileOrDir(filename string, decideDir bool) bool {
	fileInfo, err := os.Stat(filename)
	if err != nil {
		return false
	}
	isDir := fileInfo.IsDir()
	if decideDir {
		return isDir
	}
	return !isDir
}

// ReadFile 从指定的文件读取内容
func ReadFile(path string) ([]byte, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	} else {
		return b, nil
	}
}

// WriteFile 将文件写入到指定的位置
func WriteFile(path string, content []byte, perm os.FileMode) error {
	return ioutil.WriteFile(path, content, perm)
}

// AppendWriteFile 通过追加的方式将文件写入到指定的位置
func AppendWriteFile(path string, content []byte, perm os.FileMode, append bool) error {
	var f *os.File
	var err error
	if append {
		f, err = os.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, perm)
	} else {
		f, err = os.OpenFile(path, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, perm)
	}

	defer f.Close()
	if err != nil {
		return err
	} else {
		_, err := f.Write(content)
		if err != nil {
			return err
		} else {
			f.Sync()
			f.Chmod(perm)
			return nil
		}
	}
}

// BufWriteFille 通过buffio 写入文件
func BufWriteFille(path string, content []byte, perm os.FileMode) error {
	f, err := os.Create(path)
	defer f.Close()
	if err != nil {
		return err
	} else {
		w := bufio.NewWriter(f)
		w.Write(content)
		w.Flush()
		f.Chmod(perm)
		return nil
	}
}

// FileSize 获取文件大小
/*
func FileSize(path string) int64 {
	if IsExist(path) {
		if info, err := os.Stat(path); err == nil {
			return info.Size()
		}
	}
	return int64(0)
}
*/

// Json2File 将Json内容写入文件
func Json2File(content, path string, perm os.FileMode) error {
	if IsExist(path) {
		if err := os.Remove(path); err != nil {
			return err
		}
	}

	bytes := []byte(content)
	return WriteFile(path, bytes, perm)
}

// File2Json 将Json内容从文件读取并解析
func File2Json(path string, v interface{}) error {
	if IsExist(path) {
		if bytes, err := ReadFile(path); err != nil {
			return err
		} else {
			if err := json.Unmarshal(bytes, v); err != nil {
				return err
			} else {
				return nil
			}
		}
	} else {
		return fmt.Errorf("File not exists: %s\n", path)
	}
}

func WorkDir() (string, error) {
	file, err := exec.LookPath(os.Args[0])
	if err != nil {
		return "", err
	}
	p, err := filepath.Abs(file)
	if err != nil {
		return "", err
	}
	return path.Dir(strings.Replace(p, "\\", "/", -1)), err
}

func GenImageKey(in string) (out string) {
	data, err := ioutil.ReadFile(in)
	if err != nil {
		panic(err)
	}
	h := sha1.New()
	h.Write(data)
	shaBytes := h.Sum(nil)
	ext := filepath.Ext(in)
	ext = strings.ToLower(ext)
	if ext == "jpeg" {
		ext = "jpg"
	}
	out = base64.URLEncoding.EncodeToString(shaBytes) + ext
	return
}

var (
	documentRoot string
	fileCache    map[string][]byte
	fileExtType  map[string]string
	indexFile    string
)

func FileExt(path string) string {
	ext := filepath.Ext(path)
	if len(ext) == 0 {
		return ext
	}
	return ext[1:]
}

// 获取父级目录
func GetParentDirectory(directory string) string {
	return SubStr(directory, 0, strings.LastIndex(directory, "/"))
}

// 获取当前目录
func GetCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return strings.Replace(dir, "\\", "/", -1)
}

func Writeable(file string) bool {
	//err := syscall.Accept(file, syscall.O_RDWR)
	//if err != nil {
	//	return false
	//} else {
	//	return true
	//}
	_, err := os.OpenFile(file, os.O_WRONLY, 0666)
	if err != nil {
		return false
	} else {
		return true
	}

}

// 生成随机密码
func RandPwd(size int, kind int) []byte {
	//kind:0纯数字,1小写字母,2大写字母,3数字和大小写
	ikind, kinds, result := kind, [][]int{{10, 48}, {26, 97}, {26, 65}}, make([]byte, size)
	is_all := kind > 2 || kind < 0
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < size; i++ {
		if is_all { // random ikind
			ikind = rand.Intn(3)
		}
		scope, base := kinds[ikind][0], kinds[ikind][1]
		result[i] = uint8(base + rand.Intn(scope))
	}
	return result
}

func statDir(dirPath, recPath string, includeDir, isDirOnly bool) ([]string, error) {
	dir, err := os.Open(dirPath)
	if err != nil {
		return nil, err
	}
	defer dir.Close()

	fis, err := dir.Readdir(0)
	if err != nil {
		return nil, err
	}

	statList := make([]string, 0)
	for _, fi := range fis {
		if strings.Contains(fi.Name(), ".DS_Store") {
			continue
		}

		relPath := path.Join(recPath, fi.Name())
		curPath := path.Join(dirPath, fi.Name())
		if fi.IsDir() {
			if includeDir {
				statList = append(statList, relPath+"/")
			}
			s, err := statDir(curPath, relPath, includeDir, isDirOnly)
			if err != nil {
				return nil, err
			}
			statList = append(statList, s...)
		} else if !isDirOnly {
			statList = append(statList, relPath)
		}
	}
	return statList, nil
}

func StatDir(rootPath string, includeDir ...bool) ([]string, error) {
	if !IsDir(rootPath) {
		return nil, errors.New("not a directory or does not exist: " + rootPath)
	}

	isIncludeDir := false
	if len(includeDir) >= 1 {
		isIncludeDir = includeDir[0]
	}
	return statDir(rootPath, "", isIncludeDir, false)
}

// GetBytes
// 通过给定的文件名称或者url地址以及超时时间获取文件的[]byte数据.
func GetBytes(filenameOrURL string, timeout ...time.Duration) ([]byte, error) {
	if strings.Contains(filenameOrURL, "://") {
		if strings.Index(filenameOrURL, "file://") == 0 {
			filenameOrURL = filenameOrURL[len("file://"):]
		} else {
			client := http.DefaultClient
			if len(timeout) > 0 {
				client = &http.Client{Timeout: timeout[0]}
			}
			r, err := client.Get(filenameOrURL)
			if err != nil {
				return nil, err
			}
			defer r.Body.Close()
			if r.StatusCode < 200 || r.StatusCode > 299 {
				return nil, fmt.Errorf("%d: %s", r.StatusCode, http.StatusText(r.StatusCode))
			}
			return ioutil.ReadAll(r.Body)
		}
	}
	return ioutil.ReadFile(filenameOrURL)
}

// SetBytes
// 向指定的文件设置[]byte内容.
func SetBytes(filename string, data []byte) error {
	return ioutil.WriteFile(filename, data, 0660)
}

// AppendBytes
// 向指定的文件追加[]byte内容.
func AppendBytes(filename string, data []byte) error {
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0660)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(data)
	return err
}

// GetString
// 通过给定的文件名称或者url地址以及超时时间获取文件的string数据.
func GetString(filenameOrURL string, timeout ...time.Duration) (string, error) {
	bytes, err := GetBytes(filenameOrURL, timeout...)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// SetString
// 向指定的文件设置string内容.
func SetString(filename string, data string) error {
	return SetBytes(filename, []byte(data))
}

// AppendString
// 向指定的文件追加string内容.
func AppendString(filename string, data string) error {
	return AppendBytes(filename, []byte(data))
}

// Mkdir
// 创建文件夹
func Mkdir(path string) error {
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

// MkdirForFile
// 为某个文件创建目录
func MkdirForFile(path string) (err error) {
	path = filepath.Dir(path)
	return os.MkdirAll(path, os.FileMode(0777))
}

// FileTimeModified
// 返回文件的最后修改时间
// 如果有错误则返回空time.Time.
func FileTimeModified(filename string) time.Time {
	info, err := os.Stat(filename)
	if err != nil {
		return time.Time{}
	}
	return info.ModTime()
}

// Find
// 在给定的文件夹中查找某文件.
func Find(searchDirs []string, filenames ...string) (filePath string, found bool) {
	for _, dir := range searchDirs {
		for _, filename := range filenames {
			filePath = path.Join(dir, filename)
			if IsFileExists(filePath) {
				return filePath, true
			}
		}
	}
	return "", false
}

// GetPrefix
// 获取文件的前缀.
func GetPrefix(filename string) string {
	for i := len(filename) - 1; i >= 0; i-- {
		if filename[i] == '.' {
			return filename[0:i]
		}
	}
	return ""
}

// GetExt
// 获取文件的后缀.
func GetExt(filename string) string {
	for i := len(filename) - 1; i >= 0; i-- {
		if filename[i] == '.' {
			return filename[i+1:]
		}
	}
	return ""
}

// Copy
// 将文件从原地址拷贝到目的地.
func Copy(source string, dest string) (err error) {
	sourceFile, err := os.Open(source)
	if err != nil {
		return err
	}
	defer sourceFile.Close()
	destFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destFile.Close()
	_, err = io.Copy(destFile, sourceFile)
	if err == nil {
		si, err := os.Stat(source)
		if err == nil {
			err = os.Chmod(dest, si.Mode())
		}
	}
	return err
}

// DirSize
// 返回文件夹的大小
func DirSize(path string) int64 {
	var dirSize int64 = 0
	readSize := func(path string, file os.FileInfo, err error) error {
		if !file.IsDir() {
			dirSize += file.Size()
		}
		return nil
	}
	filepath.Walk(path, readSize)
	return dirSize
}
