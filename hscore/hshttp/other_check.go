// Author stone-bird created on 2021/8/28 17:53.
// Email 1245863260@qq.com or g1245863260@gmail.com.
// Use of other check
package hshttp

import (
	"fmt"
	"github.com/stonebirdjx/logger"
	"hstream/hsflag"
	"io"
	"net/http"
	"time"
)

func otherCheck(res *http.Response) error {
	start := float64(time.Now().UnixNano())
	finalUrl := res.Request.URL.String()
	total := res.ContentLength
	buffer := make([]byte, *hsflag.MaxBytes)
	accept := 0
	for {
		bytes, readBufferErr := res.Body.Read(buffer)
		if bytes > 0 {
			accept = accept + bytes
			end := float64(time.Now().UnixNano())
			logger.GlobalLogger.Logger.Printf("%s %s checking body total:%d recv:%d waste-time:%.3fms\n",
				logger.Info,
				finalUrl,
				total,
				accept,
				(end-start)/1e6,
			)
		}
		if readBufferErr != nil {
			if readBufferErr == io.EOF {
				break
			} else {
				return fmt.Errorf("read bytes err total:%d,accpet:%d", total, accept)
			}
		}
	}
	return nil
}
