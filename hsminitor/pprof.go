// Author stone-bird created on 2021/8/28 10:24.
// Email 1245863260@qq.com or g1245863260@gmail.com.
// Use of pprof
package hsminitor

import (
	"github.com/stonebirdjx/logger"
	"hstream/hsflag"
	"os"
	"runtime/pprof"
	"runtime/trace"
	"sync"
)

var CpuFile *os.File
var HeapFile *os.File
var TraceFile *os.File
var err error
var cpuOnce sync.Once
var heapOnce sync.Once
var traceOnce sync.Once

func CpuProfile() {
	switch {
	case CpuFile != nil:
		return
	default:
		cpuFileFunc := func() {
			CpuFile, err = os.OpenFile(hsflag.CpuProfile, os.O_RDWR|os.O_CREATE, 0644)
			if err != nil {
				logger.GlobalLogger.Logger.Fatal(logger.Error, err)
			}
			logger.GlobalLogger.Logger.Println(logger.Info, "cpu profile started")
			cpuError := pprof.StartCPUProfile(CpuFile)
			if cpuError != nil {
				logger.GlobalLogger.Logger.Fatal(logger.Error, cpuError)
			}
		}
		cpuOnce.Do(cpuFileFunc)
	}
}

// 生成堆内存报告
func HeapProfile() {
	switch {
	case HeapFile != nil:
		return
	default:
		heapFileFunc := func() {
			HeapFile, err = os.OpenFile(hsflag.HeapProfile, os.O_RDWR|os.O_CREATE, 0644)
			if err != nil {
				logger.GlobalLogger.Logger.Fatal(logger.Error, err)
			}
			logger.GlobalLogger.Logger.Println(logger.Info, "heap profile started")
			heapError := pprof.WriteHeapProfile(HeapFile)
			if heapError != nil {
				logger.GlobalLogger.Logger.Fatal(logger.Error, heapError)
			}
		}
		heapOnce.Do(heapFileFunc)
	}
}

func TraceProfile() {
	switch {
	case TraceFile != nil:
		return
	default:
		traceFileFunc := func() {
			TraceFile, err = os.OpenFile(hsflag.TraceProfile, os.O_RDWR|os.O_CREATE, 0644)
			if err != nil {
				logger.GlobalLogger.Logger.Fatal(logger.Error, err)
			}
			logger.GlobalLogger.Logger.Println(logger.Info, "trace profile started")
			traceError := trace.Start(TraceFile)
			if traceError != nil {
				logger.GlobalLogger.Logger.Fatal(logger.Error, traceError)
			}
		}
		traceOnce.Do(traceFileFunc)
	}
}
