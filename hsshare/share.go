// Author stone-bird created on 2021/8/28 11:00.
// Email 1245863260@qq.com or g1245863260@gmail.com.
// Use of  
package hsshare

import (
	"fmt"
	"github.com/stonebirdjx/writefile"
	"os"
	"regexp"
	"strconv"
	"time"
)

// Check file is exists
func CheckFile(fileName string) error {
	fileInfo, err := os.Stat(fileName)
	if err != nil {
		return err
	}
	if fileInfo.Size() == 0 {
		return fmt.Errorf("%s file size is 0", fileName)
	}
	return nil
}

// file -> file.1597742509
func BackupFile(fileNames ...string) error {
	for _, fileName := range fileNames {
		_, err := os.Stat(fileName)
		// This file exists
		if err == nil {
			newFileName := fileName + "." + strconv.FormatInt(time.Now().Unix(), 10)
			err := os.Rename(fileName, newFileName)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func CreateFileWriter(writerMap map[string]string) (map[string]*writefile.FileWriter, error) {
	fileWriterMap := map[string]*writefile.FileWriter{}
	for key, fileName := range writerMap {
		fileWriter, fileWriterErr := writefile.NewWriter(fileName)
		if fileWriterErr != nil {
			return nil, fileWriterErr
		}
		fileWriterMap[key] = fileWriter
	}
	return fileWriterMap, nil
}

func GetTagUri(line string) string {
	reg := regexp.MustCompile(`URI=".*"`)
	return reg.FindString(line)
}

func GetDuration(line string) (float64, error) {
	reg := regexp.MustCompile(`[\d|.][^,]*`)
	tsDuration := reg.FindString(line)
	tsfloat, err := strconv.ParseFloat(tsDuration, 64)
	if err != nil {
		return 0, err
	}
	return tsfloat, nil
}