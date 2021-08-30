// Author stone-bird created on 2021/8/28 17:46.
// Email 1245863260@qq.com or g1245863260@gmail.com.
// Use of  
package hshttp

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/stonebirdjx/logger"
	"hstream/hsflag"
	"net/http"
	"strings"
)

func get(url string) (*http.Response, error) {
	c := &http.Client{}
	c.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	header := strings.TrimSpace(*hsflag.Header)
	if header != "" {
		var mp map[string]interface{}
		err := json.Unmarshal([]byte(header), &mp)
		if err != nil {
			return nil, err
		}
		if _, ok := mp["User-Agent"]; !ok {
			req.Header.Set("User-Agent", hsflag.UserAgent)
		}
		for k, v := range mp {
			req.Header.Set(k, fmt.Sprint(v))
		}
	} else {
		req.Header.Set("User-Agent", hsflag.UserAgent)
	}
	res, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func HttpCheck(url string) error {
	logger.GlobalLogger.Logger.Println(logger.Info, "now checking url", url)
	res, err := get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	statusCode := res.StatusCode
	finalUrl := res.Request.URL.String()
	if !(statusCode >= 200 && statusCode < 300) {
		return fmt.Errorf("%s return code %d", finalUrl, statusCode)
	}
	contentType := res.Header.Get("Content-Type")
	switch contentType {
	case "application/x-mpegURL", "audio/x-mpegurl", "application/vnd.apple.mpegurl":
		m3u8, err := m3u8Check(res)
		if err != nil {
			return fmt.Errorf("%s check err -> %s", finalUrl, err)
		}
		switch m3u8.mediaType {
		case 1:
			logger.GlobalLogger.Logger.Println(logger.Info, finalUrl, "is parent m3u8")
			err = m3u8.lowerM3u8Check()
		case 2:
			if m3u8.islive {
				logger.GlobalLogger.Logger.Println(logger.Info, finalUrl, "is live m3u8")
				err = m3u8.liveCheck()
			} else {
				logger.GlobalLogger.Logger.Println(logger.Info, finalUrl, "is vod m3u8")
				err = m3u8.vodCheck()
			}
		}
	default:
		if *hsflag.Chunk {
			err = otherCheck(res)
		}
	}
	if err != nil {
		return fmt.Errorf("%s check err -> %s", finalUrl, err)
	}
	return nil
}
