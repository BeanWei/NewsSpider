package main

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/gocolly/colly"
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

//定义channel
// var (
// 	chanNewslist chan string //负责传输所有的焦点新闻链接
// 	chanOk       chan string //负责收集焦点新闻爬取状态
// )

// func news_36kr() {

// }

func main() {
	html, err := html("http://36kr.com/")
	if err != nil {
		log.Fatal(err)
	}
	homepage := strings.Split(strings.Split(html, `"hotPosts|hotPost":`)[1], `,"highProjects|focus":`)[0]
	hotNewsList := &News{}
	//string转成[]map[string]interface{}
	//var dat []map[string]interface{}
	if err := json.Unmarshal([]byte(homepage), hotNewsList); err != nil {
		log.Fatal(err)
	}
	allnews := []NewsStruct{}
	for _, onenews := range *hotNewsList {
		news := NewsStruct{}
		id := onenews.ID
		news.Srcurl = fmt.Sprintf("http://36kr.com/p/%s.html", id)
		news.Title = onenews.Title
		news.Coverlink = onenews.Coverlink
		news.AuthorName = onenews.Author.Name
		news.AuthorAvatar = onenews.Author.Avatar
		news.AuthorIntroduction = onenews.Author.Introduction
		news.Publishtime = onenews.Publishtime
		news.Content, err = newsdetail(id)
		if err != nil {
			log.Fatal(err)
			continue
		}
		allnews = append(allnews, news)
	}
	fmt.Println(allnews)

}

//html 爬取页面通用函数
func html(url string) (html string, err error) {
	c := colly.NewCollector()
	//设置请求超时
	c.SetRequestTimeout(200 * time.Second)
	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/63.0.3239.108 Safari/537.36"
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Host", "36kr.com")
		r.Headers.Set("Referer", "http://36kr.com/")
		r.Headers.Set("X-Tingyun-Id", "Dio1ZtdC5G4;r=996634586")
	})
	c.OnResponse(func(resp *colly.Response) {
		html = string(resp.Body)
	})

	c.OnError(func(resp *colly.Response, errHttp error) {
		err = errHttp
	})
	err = c.Visit(url)
	return
}

//newsdetail 获取文章正文详情
func newsdetail(id string) (content string, err error) {
	srcurl := fmt.Sprintf("http://36kr.com/p/%s.html", id)
	newshtml, err := html(srcurl)
	if err != nil {
		log.Fatal(err)
	}
	reg := regexp.MustCompile(`"content":(.*?),"cover"`)
	content = reg.FindString(newshtml)
	fmt.Println(content)
	return
}
