// Author stone-bird created on 2021/8/25 21:08.
// Email 1245863260@qq.com or g1245863260@gmail.com.
// Use of http stream
package main

import (
	"github.com/stonebirdjx/logger"
	"github.com/stonebirdjx/readfile"
	"hstream/hscore"
	"hstream/hsflag"
	"hstream/hsminitor"
	"os"
	"runtime/pprof"
	"runtime/trace"
	"sync"
)

func init() {
	logger.NewLogger(hsflag.LogName, hsflag.LogPrefix)
}

func main() {
	defer logger.GlobalLogger.LogFile.Close()
	logger.GlobalLogger.Logger.Println(logger.Info, logger.LogBegin)
	// signal
	if *hsflag.Signal {
		logger.GlobalLogger.Logger.Println(logger.Info, "turn on system signal monitoring")
		go hsminitor.HsSignal()
	} else {
		logger.GlobalLogger.Logger.Println(logger.Info, "turn off system signal monitoring")
	}

	// pporf
	if *hsflag.Pporf {
		logger.GlobalLogger.Logger.Println(logger.Info, "turn on system pprof monitoring")
		hsminitor.HeapProfile()
		hsminitor.CpuProfile()
		hsminitor.TraceProfile()
		defer hsminitor.HeapFile.Close()
		defer hsminitor.CpuFile.Close()
		defer hsminitor.TraceFile.Close()
		defer pprof.StopCPUProfile()
		defer trace.Stop()
	} else {
		logger.GlobalLogger.Logger.Println(logger.Info, "turn off system pprof monitoring")
	}

	var api hscore.RuntimeApi
	contentFile :=  *hsflag.UrlFile
	api = &hscore.RuntimeBase{
		UrlFile:       contentFile,
		SuccessFile:   hsflag.SuccessFile,
		FailureFile:   hsflag.FailureFile,
		FileWriterMap: nil,
	}
	if err := api.FileCheck(); err != nil {
		logger.GlobalLogger.Logger.Println(logger.Error, "file check err ->", err)
		os.Exit(1)
	}

	if err := api.Writer(); err != nil {
		logger.GlobalLogger.Logger.Println(logger.Error, "make write file err ->", err)
		os.Exit(1)
	}
	var wg sync.WaitGroup
	fileReader, err := readfile.NewReader(contentFile)
	if err != nil {
		logger.GlobalLogger.Logger.Println(logger.Error, contentFile, "make file reader fail")
		os.Exit(1)
	}
	go fileReader.Scanner()
	for i := 0; i < *hsflag.Online; i++ {
		wg.Add(1)
		go func() {
			api.Core(fileReader)
			wg.Done()
		}()
	}
	wg.Wait()
	api.WriterClose()
	logger.GlobalLogger.Logger.Println(logger.Info, logger.LogEnd)
}
