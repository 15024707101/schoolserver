package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

var (
	successNum int = 0
	failNum    int = 0
	returnNum  int = 0
	allNum     int = 0
	oldString      = ""
	newString      = ""
)

func main() {
	initFileName()
	ReSetPhotoNames()
	returnNum = allNum - failNum - successNum
	fmt.Printf("该文件夹下共发现 %v个文件和目录,%v个重命名成功，%v个重命名失败，%v个没匹配到【%v】", allNum, successNum, failNum, returnNum,oldString)
}

func ReSetPhotoNames() {
	photoFolder := `./`

	files, _ := ioutil.ReadDir(photoFolder)
	allNum = len(files) - 2
	for _, file := range files {

		if file.IsDir() {
			continue
		} else {
			fileName := file.Name()

			index := strings.Index(fileName, oldString)
			if index == -1 {
				continue
			}

			newFileName := ""
			dotIndex := strings.LastIndex(fileName, ".")
			if dotIndex != -1 && dotIndex != 0 {

				newFileName = fileName[:dotIndex]
				newFileName = strings.Replace(newFileName, oldString, newString, -1)

				newFileName += fileName[dotIndex:]
			}
			err := os.Rename(photoFolder+fileName, photoFolder+newFileName)
			if err != nil {
				fmt.Println("reName Error", err)
				failNum++
				continue
			}
			successNum++
		}
	}
}

func initFileName() {
	// 读取一个文件的内容
	file, err := os.Open("./重命名工具配置.txt")
	if err != nil {
		fmt.Println("open file err:", err.Error())
		return
	}

	// 处理结束后关闭文件
	defer file.Close()

	// 使用bufio读取
	r := bufio.NewReader(file)

	data, _ := r.ReadBytes('\n')
	oldString=string(data)
	data, _ = r.ReadBytes('\n')
	newString=string(data)

	oldString = strings.Replace(oldString, " ", "", -1)
	oldString = strings.Replace(oldString, "\n", "", -1)
	oldString = strings.Replace(oldString, "\r", "", -1)

	newString = strings.Replace(newString, " ", "", -1)
	newString = strings.Replace(newString, "\n", "", -1)
	newString = strings.Replace(newString, "\r", "", -1)

}
