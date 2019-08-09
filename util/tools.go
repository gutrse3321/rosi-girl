package util

import (
	"github.com/gocolly/colly"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"time"
)

type Tools struct {
}

var (
	Tool   = New()
	once   sync.Once
	logger ILogger
)

func init() {
	logger = &Logger{}
}

func New() (t *Tools) {
	once.Do(func() {
		t = &Tools{}
	})
	return t
}

func (t *Tools) ReplaceAll(s, old, newS string) (result string) {
	result = strings.Replace(s, old, newS, -1)
	return
}

func (t *Tools) CheckDirExist(path string) error {
	logger.Normal("START FILE DIR CHECK ... Please waiting for me ...")
	_, err := os.Stat(path)

	if !os.IsNotExist(err) {
		dirList, e := ioutil.ReadDir(path)
		if e != nil {
			return e
		}
		for _, v := range dirList {
			os.RemoveAll(path + v.Name())
		}
	} else {
		os.Mkdir(path, os.ModePerm)
	}
	logger.Normal("End FILE DIR CHECK!!!")
	return nil
}

func (t *Tools) SetHeader(r *colly.Request) {
	r.Headers.Set("Host", "rosi8.cc")
	r.Headers.Set("Referer", "http://rosi8.cc/")
	r.Headers.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3")
	r.Headers.Set("Accept-Encoding", "gzip, deflate")
	r.Headers.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8,ja;q=0.7")
	r.Headers.Set("Cache-Control", "max-age=0")
	r.Headers.Set("Connection", "keep-alive")
	r.Headers.Set("Cookie", "UM_distinctid=16c74324ef782-03109bcb9cf37f-7373e61-1fa400-16c74324ef835b; CNZZDATA1254468699=1470998860-1565316408-https%253A%252F%252Fwww.baidu.com%252F%7C1565337233")
	//Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/76.0.3809.100 Safari/537.36
}

func (t *Tools) ReadyGo(s int) {
	ch := t.ticker(s)
	time.Sleep(time.Duration(s) * time.Second)
	ch <- true
	close(ch)
}

func (t *Tools) ticker(s int) chan bool {
	ticker := time.NewTicker(time.Second)
	stopChan := make(chan bool)
	go func(ticker *time.Ticker) {
		num := s
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				logger.Underline(num)
				num--
			case stop := <-stopChan:
				if stop {
					logger.Normal("========= Game Start =========")
					return
				}
			}
		}
	}(ticker)
	return stopChan
}
