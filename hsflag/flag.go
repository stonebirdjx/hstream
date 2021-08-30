// Author stone-bird created on 2021/8/25 21:14.
// Email 1245863260@qq.com or g1245863260@gmail.com.
// Use of hstream flag
package hsflag

import (
	"flag"
	"fmt"
	"os"
	"runtime"
)

var (
	Online   = flag.Int("u", 5, "enter the number of concurrent")
	MaxBytes = flag.Int("byte", 10485760, "max byte to read response body")
	LiveDial = flag.Int("live", 600, "max byte to read response body")
	Chunk    = flag.Bool("chunk", false, "whether check the http body data")
	Pporf    = flag.Bool("pporf", false, "whether turn on system pprof")
	Signal   = flag.Bool("signal", true, "whether turn on system signal")
	Header   = flag.String("headers", "", "set http request header")
	UrlFile  = flag.String("urlfile", "url.txt", "enter url file to batch check")
	v        = flag.Bool("v", false, "show hstream tool version")
	version  = flag.Bool("version", false, "show hstream tool help text")
	h        = flag.Bool("h", false, "show hstream tool help text")
	help     = flag.Bool("help", false, "show hstream tool help text")
)

func init() {
	flag.Parse()

	if *h || *help {
		flag.Usage()
		os.Exit(0)
	}

	if *v || *version {
		fmt.Printf("%s %s/%s %s\n", name, ver, runtime.GOOS, jx)
		os.Exit(0)
	}

}
