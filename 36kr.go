package main

import (
	"regexp"
	"fmt"
	"strings"
	"log"

	"github.com/gocolly/colly"
)

/*定义新闻的结构*/
type NewsStruct struct {
	Srcurl      string
	Title       string
	Coverlink   string
	Author      []NewsAutor
	Publishtime string
	Content     string
}

/*定义作者的结构*/
type NewsAutor struct {
	Name				string
	Avatar				string
	Introduction		string
}

//定义channel
var (
	chan_newslist chan string		//负责传输所有的焦点新闻链接
	chan_ok chan string 			//负责收集焦点新闻爬取状态
)

// func news_36kr() {

// }


func main() {
	html, err := html("http://36kr.com/")
	if err != nil {
		log.Fatalln("无法获取文章页面")
	}
	hot_news_list := strings.Split(strings.Split(html, `"hotPosts|hotPost":`)[1], `,"highProjects|focus":`)[0]
	chan_newslist := make(chan string, 10)
	for _, v := range hot_news_list {
		chan_newslist <- v
	}
	chan_ok := make(chan string, 10)
	for onenews := range chan_newslist {
		go newsdetail(onenews)
		<-chan_ok
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
func newsdetail(onenews string) (news *NewsStruct, err error) {
	// news := make(*NewsStruct, 6)
	// author := make(*NewsAutor, 1)
	author := NewsAutor{}
	news.Srcurl := fmt.Sprintf("http://36kr.com/p/%s.html", onenews["id"])
	news.Title = onenews["title"]	
	news.Coverlink = onenews["cover"]
	author.Name = onenews["user"]["name"]
	author.Avatar = onenews["user"]["avatar_url"]
	author.Introduction = onenews["user"]["introduction"]
	news.Publishtime = onenews["published_at"]
	news.Content, err := newscontent(news.Srcurl)
	if err != nil {
		log.Fatalln("文章内容解析失败")
		return nil, err
	}
	fmt.Println("news")
	return
}

//newscontent 解析文章正文
func newscontent(url string) (content string, err error) {
	newshtml, err := html(url)
	if err != nil {
		log.Fatalln("无法获取文章页面")
		return nil, err 
	}
	reg := regexp.MustCompile(`"content":(.*?),"cover"`)
	content = reg.FindAllString(newshtml)
	return
}


