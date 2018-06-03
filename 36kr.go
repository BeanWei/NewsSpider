package main

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/gocolly/colly"
)

// NewsStruct 定义新闻的结构
type NewsStruct struct {
	Srcurl             string
	Title              string
	Coverlink          string
	AuthorName         string
	AuthorAvatar       string
	AuthorIntroduction string
	Publishtime        string
	Content            string
}

// OneNews 爬取的每条news结构
type OneNews struct {
	Srcurl      string   `json:"id"`
	Title       string   `json:"title"`
	Coverlink   string   `json:"avatar_url"`
	Author      []Author `json:"user"`
	Publishtime string   `json:"published_at"`
	Content     string   `json:"content"`
}

// Author 爬取的每条news中的作者信息
type Author struct {
	Name         string `json:"name"`
	Avatar       string `json:"avatar_url"`
	Introduction string `json:"introduction"`
}

//定义channel
var (
	chanNewslist chan string //负责传输所有的焦点新闻链接
	chanOk       chan string //负责收集焦点新闻爬取状态
)

// func news_36kr() {

// }

func main() {
	html, err := html("http://36kr.com/")
	if err != nil {
		log.Fatalln("无法获取文章页面")
	}
	hotNewsList := strings.Split(strings.Split(html, `"hotPosts|hotPost":`)[1], `,"highProjects|focus":`)[0]
	channewslist := make(chan string, 10)
	//string转成[]map[string]interface{}
	var dat []map[string]interface{}
	if err := json.Unmarshal([]byte(hotNewsList), &dat); err != nil {
		log.Fatalln("无法序列化数据")
	}
	for _, v := range dat {
		chanNewslist <- v
	}
	chanOk := make(chan string, 10)
	for onenews := range chanNewslist {
		go newsdetail(onenews)
		<-chanOk
	}

}

//html 爬取页面通用函数
func html(url string) (html string, err error) {
	c := colly.NewCollector()
	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/63.0.3239.108 Safari/537.36"
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Host", "36kr.com")
		r.Headers.Set("Referer", "http://36kr.com/")
		r.Headers.Set("X-Tingyun-Id", "Dio1ZtdC5G4;r=996634586")
	})
	c.OnResponse(func(resp *colly.Response) {
		html = string(resp.Body)
		// reg := regexp.MustCompile(`^hotPost":(.*?)*,"highProjects`)
		// hotlist = reg.FindString(string(resp.Body))
		//hotlist = strings.Split(strings.Split(string(resp.Body), `"hotPosts|hotPost":`)[1], `,"highProjects|focus":`)[0]
	})

	c.OnError(func(resp *colly.Response, errHttp error) {
		err = errHttp
	})
	err = c.Visit(url)
	return
}

//newsdetail 获取文章详情
func newsdetail(onenews *OneNews) (news []NewsStruct, err error) {
	srcurl := fmt.Sprintf("http://36kr.com/p/%s.html", onenews["id"])
	title := onenews["title"]
	cover := onenews["cover"]
	authorname := onenews["user"]["name"]
	authoravatar := onenews["user"]["avatar_url"]
	authorintroduction := onenews["user"]["introduction"]
	publishtime := onenews["published_at"]
	content, err := newscontent(srcurl)
	if err != nil {
		log.Fatalln("文章内容解析失败")
		return nil, err
	}
	n := NewsStruct{
		Srcurl:             srcurl,
		Title:              title,
		Coverlink:          cover,
		AuthorName:         authorname,
		AuthorAvatar:       authoravatar,
		AuthorIntroduction: authorintroduction,
		Publishtime:        publishtime,
		Content:            content,
	}
	news = append(news, n)
	fmt.Println(news)
	return
}

//newscontent 解析文章正文
func newscontent(url string) (content string, err error) {
	newshtml, err := html(url)
	if err != nil {
		log.Fatalln("无法获取文章页面")
		panic(err)
	}
	reg := regexp.MustCompile(`"content":(.*?),"cover"`)
	content = reg.FindString(newshtml)
	return
}
