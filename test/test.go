package main

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
	"github.com/PuerkitoBio/goquery"
)

// NewsStruct 定义新闻的结构
type NewsStruct struct {
	Srcurl             string `json:"srcurl"`
	Title              string `json:"title"`
	Coverlink          string `json:"coverlink"`
	AuthorName         string `json:"authorname"`
	AuthorAvatar       string `json:"authoravatar"`
	AuthorIntroduction string `json:"authorintroduction"`
	Publishtime        string `json:"publishtime"`
	Content            string `json:"content"`
}

// News 爬取news结构
type News []struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Coverlink string `json:"avatar_url"`
	Author    struct {
		Name         string `json:"name"`
		Avatar       string `json:"avatar_url"`
		Introduction string `json:"introduction"`
	} `json:"user"`
	Publishtime string `json:"published_at"`
	Content     string `json:"content"`
}

// PostData 请求数据
type PostData struct {
	Code	string
	Page	int
	Time	int64
}

// PostResult 请求返回数据
type PostResult struct {
	Data		string 		`json:"data"`
	Dateline	string		`json:"last_dateline"`
	Msg			string		`json:"msg"`
	Result		int			`json:"result"`
	TotalPage	int			`json:"total_page"`
}

func main() {

	allnews := []NewsStruct{}

	c := colly.NewCollector(
		colly.Async(true),
	)
	extensions.RandomUserAgent(c)
	extensions.Referrer(c)
	//设置请求超时
	c.SetRequestTimeout(200 * time.Second)
	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/63.0.3239.108 Safari/537.36"
	for page := 1; page <= 10; page++ {
		dateline := time.Now().Unix()
		err := c.Post("https://www.huxiu.com/v2_action/article_list", 
			map[string]PostData{
				"huxiu_hash_code": "dc6ad039c0702be4185b5918168c08c3",
				"page": page,
				"last_dateline": dateline
			})
		if err != nil {
			log.Fatalf(err)
		}
	}
	
	c.OnResponse(func(resp *colly.Response) {
		postresult := &PostResult{}
		if err := json.Unmarshal([]byte(resp.Body), postresult); err != nil {
			log.Fatalf(err)
		}
		newsdata := postresult.Data
		doc, err := goquery.NewDocumentFromReader(newsdata)
		if err != nil {
			log.Fatalf(err)
		}
		//TODO:根据test.md解析出news信息，然后进一步爬取
	})
		
	c.OnError(func(r *colly.Response, err error) {
		log.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
		r.Request.Retry()
	})

	
}
