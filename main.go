/**
 * @Author: Tomonori
 * @Date: 2019/8/9 16:39
 * @File: main
 * @Desc:
 */
package main

import (
	"fmt"
	"github.com/gocolly/colly"
	"log"
	"os"
	"rosi/util"
	"strconv"
	"strings"
)

const (
	outputDir = "./out/"
)

var (
	err       error
	logger    util.ILogger
	pageTotal *int
)

func init() {
	logger = &util.Logger{}
}

func main() {
	if err = util.New().CheckDirExist(outputDir); err != nil {
		log.Panicln("file dir check failed:", err)
	}

	pointerWrapper := 1
	pageTotal = &pointerWrapper
	pageTotalValue := util.Scanf("Input your crawl page size:")
	if pageTotalValue != "" {
		*pageTotal, err = strconv.Atoi(pageTotalValue)
		if err != nil {
			logger.Info("please input really number type")
		}
	}

	ch := make(chan int, *pageTotal)
	i := 1
	for {
		if i > *pageTotal {
			break
		}
		go crawlHomePage(i, ch)
		i++
	}
	cashPool := 0
	chNum := *pageTotal
	for cashPool < chNum {
		cashPool += <-ch
	}
}

func crawlHomePage(pageIndex int, ch chan int) {
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.169 Safari/537.36"),
	)

	c.OnRequest(func(r *colly.Request) {
		util.New().SetHeader(r)
		logger.Info("【Visiting】", r.URL.String())
	})

	// 获取 a标签元素的详情链接
	c.OnHTML(".i20 a[href]", func(e *colly.HTMLElement) {
		d := c.Clone()
		requestDetailPage(d, e)
	})

	c.Visit(fmt.Sprintf("http://rosi8.cc/rosixiezhen/list1%d.html", pageIndex))
	ch <- 1
}

func requestDetailPage(c *colly.Collector, e *colly.HTMLElement) {
	var title *string
	link := e.Attr("href")
	nextFlag := true

	c.OnRequest(func(dr *colly.Request) {
		util.New().SetHeader(dr)
		logger.Normal("【Detail Visiting】: ", dr.URL.String())
	})

	// 获取标题名称创建漫画文件夹
	c.OnHTML(".info-bt .title", func(de *colly.HTMLElement) {
		title = &de.Text
		groupDir := fmt.Sprintf("%s%s", outputDir, *title)
		_ = os.Mkdir(groupDir, os.ModePerm)
		logger.Normal("【detail title】: ", de.Text)
	})

	// 获取图片
	c.OnHTML(".info-zi > img", func(dde *colly.HTMLElement) {
		d := c.Clone()
		link := dde.Attr("src")
		logger.Info(link)
		if link == "" {
			nextFlag = false
		}
		reduceImage(d, dde, title)
	})

	//下一页
	c.OnHTML(".pagebar > a:last-child", func(e *colly.HTMLElement) {
		next := e.Attr("href")
		logger.Info("【HEAD】Start Image Download！！", next)
		d := c.Clone()
		nextImageGiveMe(d, e, title, &nextFlag)
	})

	logger.Complate("End Request")
	c.Visit(e.Request.AbsoluteURL(link))
}

func nextImageGiveMe(c *colly.Collector, mapEl *colly.HTMLElement, title *string, nextFlag *bool) {
	if *nextFlag {
		next := mapEl.Attr("href")
		c.OnRequest(func(mr *colly.Request) {
			util.New().SetHeader(mr)
		})

		// 获取图片存储图片
		c.OnHTML(".info-zi > img", func(e *colly.HTMLElement) {
			d := c.Clone()
			link := e.Attr("src")
			logger.Info(link)
			if link == "" {
				*nextFlag = false
			}
			reduceImage(d, e, title)
		})

		c.OnHTML(".pagebar > a:last-child", func(e *colly.HTMLElement) {
			next := e.Attr("href")
			logger.Info("【Next Image Page】: ", next)
			d := c.Clone()
			nextImageGiveMe(d, e, title, nextFlag)
		})

		c.Visit(mapEl.Request.AbsoluteURL(next))
	}
}

func reduceImage(c *colly.Collector, e *colly.HTMLElement, title *string) {
	link := e.Attr("src")

	c.OnRequest(func(rr *colly.Request) {
		util.New().SetHeader(rr)
	})

	c.OnResponse(func(r *colly.Response) {
		if strings.Index(r.Headers.Get("Content-Type"), "image") > -1 {
			r.Save(fmt.Sprintf("%s%s/", outputDir, *title) + r.FileName())
			return
		}
	})

	c.Visit("http://rosi8.cc" + link)
}
