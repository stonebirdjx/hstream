// Author stone-bird created on 2021/8/29 8:45.
// Email 1245863260@qq.com or g1245863260@gmail.com.
// Use of  m3u8 check
package hshttp

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"github.com/stonebirdjx/logger"
	"hstream/hsflag"
	"hstream/hsshare"
	"net/http"
	"strings"
	"sync"
	"time"
)

type streamM3u8 struct {
	currentUrl string
	mediaType  int // 1 m3u8 ;2 ts
	islive     bool
	m3u8list   []string
	tslist     []string
	tsDuration []float64
	lastTs     string
}

func m3u8Check(res *http.Response) (streamM3u8, error) {
	finalUrl := res.Request.URL.String()
	prefixURLList := strings.Split(finalUrl, "/")
	prefixURL := strings.TrimRight(finalUrl, prefixURLList[len(prefixURLList)-1])

	isLive := true
	child := false
	var m3u8 streamM3u8
	m3u8.currentUrl = finalUrl
	bs := bufio.NewScanner(res.Body)
	for bs.Scan() {
		line := bs.Text()
		line = strings.TrimSpace(line)
		tagUrl := ""
		if strings.HasPrefix(line, "#") {
			if strings.HasPrefix(line, "#EXT") {
				// check tag has uri ?
				tagUri := hsshare.GetTagUri(line)
				if tagUri != "" {
					tagSuffix := strings.Split(tagUri, "\"")[1]
					if strings.HasPrefix(tagSuffix, "http") {
						tagUrl = tagSuffix
					} else {
						tagUrl = prefixURL + tagSuffix
					}
					res, err := get(tagUrl)
					if err != nil {
						return m3u8, fmt.Errorf("%s %s request err -> %s", tagUrl, line, err.Error())
					}
					err = otherCheck(res)
					if err != nil {
						return m3u8, fmt.Errorf("%s %s check body err -> %s", tagUrl, line, err.Error())
					}
				}

				// is live
				if line == "#EXT-X-ENDLIST" {
					isLive = false
					continue
				}

				// is ts
				if strings.HasPrefix(line, "#EXTINF") {
					tsDuration, err := hsshare.GetDuration(line)
					if err != nil {
						return m3u8, fmt.Errorf("%s %s line get duration err", finalUrl, line)
					}
					m3u8.tsDuration = append(m3u8.tsDuration, tsDuration)
					child = true
					continue
				}
			}
			continue
		}
		if strings.HasPrefix(line, "http") {
			tagUrl = line
		} else {
			tagUrl = prefixURL + line
		}
		if child {
			m3u8.tslist = append(m3u8.tslist, tagUrl)
		} else {
			m3u8.m3u8list = append(m3u8.m3u8list, tagUrl)
		}
	}
	if child {
		m3u8.mediaType = 2
		m3u8.islive = isLive
		if len(m3u8.tsDuration) != len(m3u8.tslist) {
			return m3u8, fmt.Errorf("%s EXTINF tag and media is not match", finalUrl)
		}
	} else {
		m3u8.mediaType = 1
		m3u8.islive = false
	}

	return m3u8, nil
}

func (sm streamM3u8) lowerM3u8Check() error {
	var wg sync.WaitGroup
	errChan := make(chan error, len(sm.m3u8list))
	for _, m3u8Url := range sm.m3u8list {
		wg.Add(1)
		go func(m3u8Url string) {
			err := HttpCheck(m3u8Url)
			errChan <- err
			wg.Done()
		}(m3u8Url)
	}
	wg.Wait()
	close(errChan)
	var retErr error
	for err := range errChan {
		if err != nil {
			if retErr == nil {
				retErr = err
			} else {
				retErr = errors.New(retErr.Error() + " " + err.Error())
			}
		}
	}
	return retErr
}

func (sm streamM3u8) liveCheck() error {
	ch := make(chan error)
	cxt, _ := context.WithTimeout(context.Background(), time.Second*time.Duration(*hsflag.LiveDial))
	go func() {
		err := sm.liveRequest(cxt)
		ch <- err
		close(ch)
	}()
	select {
	case <-cxt.Done():
		logger.GlobalLogger.Logger.Println(logger.Info, sm.currentUrl, "live check time end")
		return nil
	case err := <-ch:
		if err != nil {
			return err
		}
	}
	return nil
}

func (sm streamM3u8) liveRequest(cxt context.Context) error {
	for index, tsUrl := range sm.tslist {
		select {
		case <-cxt.Done():
			return nil
		default:
			if sm.lastTs == tsUrl {
				continue
			}
			err := HttpCheck(tsUrl)
			if err != nil {
				return err
			}
			time.Sleep(time.Millisecond * time.Duration(sm.tsDuration[index]*1000))
		}
	}
	select {
	case <-cxt.Done():
		return nil
	default:
		sm.tsDuration = []float64{}
		sm.tslist = []string{}
		err := sm.liveMoreRequest()
		if err != nil {
			return err
		}
		err = sm.liveRequest(cxt)
		if err != nil {
			return err
		}
	}
	return nil
}

func (sm streamM3u8) liveMoreRequest() error {
	res, err := get(sm.currentUrl)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	finalUrl := res.Request.URL.String()
	prefixURLList := strings.Split(finalUrl, "/")
	prefixURL := strings.TrimRight(finalUrl, prefixURLList[len(prefixURLList)-1])
	bs := bufio.NewScanner(res.Body)
	for bs.Scan() {
		line := bs.Text()
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "#") {
			// check tag has uri ?
			tagUri := hsshare.GetTagUri(line)
			if tagUri != "" {
				tagUrl := ""
				tagSuffix := strings.Split(tagUri, "\"")[1]
				if strings.HasPrefix(tagSuffix, "http") {
					tagUrl = tagSuffix
				} else {
					tagUrl = prefixURL + tagSuffix
				}
				res, err := get(tagUrl)
				if err != nil {
					return fmt.Errorf("%s %s request err -> %s", tagUrl, line, err.Error())
				}
				err = otherCheck(res)
				if err != nil {
					return fmt.Errorf("%s %s check body err -> %s", tagUrl, line, err.Error())
				}
			}

			if strings.HasPrefix(line, "#EXTINF") {
				tsDuration, err := hsshare.GetDuration(line)
				if err != nil {
					return fmt.Errorf("%s %s line get duration err", finalUrl, line)

				}
				sm.tsDuration = append(sm.tsDuration, tsDuration)
				continue
			}
			continue
		}
		if !strings.HasPrefix(line, "http") {
			line = prefixURL + line
		}
		sm.tslist = append(sm.tslist, line)
	}
	if len(sm.tslist) != len(sm.tsDuration) {
		return fmt.Errorf("%s EXTINF tag and media is not match", finalUrl)
	}
	return nil
}

func (sm streamM3u8) vodCheck() error {
	for _, tsUrl := range sm.tslist {
		err := HttpCheck(tsUrl)
		if err != nil {
			return err
		}
	}
	return nil
}
