// Author stone-bird created on 2021/8/27 9:00.
// Email 1245863260@qq.com or g1245863260@gmail.com.
// Use of  
package hscore

import (
	"github.com/stonebirdjx/logger"
	"github.com/stonebirdjx/readfile"
	"github.com/stonebirdjx/writefile"
	"hstream/hscore/hshttp"
	"hstream/hsshare"
	"net/url"
	"strings"
)

type RuntimeBase struct {
	UrlFile       string
	SuccessFile   string
	FailureFile   string
	FileWriterMap map[string]*writefile.FileWriter
}

type RuntimeApi interface {
	FileCheck() error
	Writer() error
	Core(fileReader *readfile.FileReader)
	WriterClose()
}

func (rb *RuntimeBase) FileCheck() error {
	if err := hsshare.CheckFile(rb.UrlFile); err != nil {
		return err
	}

	if err := hsshare.BackupFile(rb.FailureFile, rb.SuccessFile); err != nil {
		return err
	}
	return nil
}

func (rb *RuntimeBase) Writer() error {
	var writerFileMap = map[string]string{}
	var err error
	writerFileMap["succ"] = rb.SuccessFile
	writerFileMap["fail"] = rb.FailureFile
	rb.FileWriterMap, err = hsshare.CreateFileWriter(writerFileMap)
	if err != nil {
		return err
	}
	return nil
}

func (rb *RuntimeBase) Core(fileReader *readfile.FileReader) {
	for line := range fileReader.ContentChan {
		content := strings.TrimSpace(line)
		if content == "" || strings.HasPrefix(content, "#") {
			continue
		}
		u, err := url.Parse(content)
		if err != nil {
			rb.FileWriterMap["fail"].WriteString(line + " is not net/url")
			continue
		}
		switch u.Scheme {
		case "http", "https":
			err := hshttp.HttpCheck(content)
			if err != nil {
				rb.FileWriterMap["fail"].WriteString(line + " check server err -> " + err.Error())
				logger.GlobalLogger.Logger.Println(logger.Error, content, "server err ->", err)
			} else {
				rb.FileWriterMap["succ"].WriteString(line)
				logger.GlobalLogger.Logger.Println(logger.Info, content, "server is ok")
			}
		default:
			rb.FileWriterMap["fail"].WriteString(line + " is not http or https protocol url")
			logger.GlobalLogger.Logger.Println(logger.Error, line, "is not http or https protocol url")
		}
	}
}

func (rb *RuntimeBase) WriterClose() {
	for key := range rb.FileWriterMap {
		err := rb.FileWriterMap[key].File.Close()
		if err != nil {
			logger.GlobalLogger.Logger.Println(logger.Error, key, "close file err")
		} else {
			logger.GlobalLogger.Logger.Println(logger.Info, key, "close file success")
		}
	}
}
