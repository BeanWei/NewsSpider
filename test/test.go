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

func main() {

	allnews := []NewsStruct{}

	c := colly.NewCollector()
	//设置请求超时
	c.SetRequestTimeout(200 * time.Second)
	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/63.0.3239.108 Safari/537.36"
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Host", "36kr.com")
		r.Headers.Set("Referer", "http://36kr.com/")
		r.Headers.Set("X-Tingyun-Id", "Dio1ZtdC5G4;r=996634586")
		log.Println("Visiting: ", r.URL.String())
	})
	c.OnError(func(r *colly.Response, err error) {
		log.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
		r.Request.Retry()
	})

	detailCollector := c.Clone()
	detailCollector.OnHTML("html", func(e *colly.HTMLElement) {

		reg := regexp.MustCompile(`"content":(.*?),"cover"`)
		content := reg.FindString(string(e.Text))

		newsdetail := NewsStruct{}
		newsdetail.Srcurl = e.Request.URL.String()
		newsdetail.Title = e.Request.Ctx.Get("title")
		newsdetail.Coverlink = e.Request.Ctx.Get("coverlink")
		newsdetail.AuthorName = e.Request.Ctx.Get("authorname")
		newsdetail.AuthorAvatar = e.Request.Ctx.Get("authoravatar")
		newsdetail.AuthorIntroduction = e.Request.Ctx.Get("authorintroduction")
		newsdetail.Publishtime = e.Request.Ctx.Get("publishtime")
		newsdetail.Content = content
		allnews = append(allnews, newsdetail)

		fmt.Println(len(allnews))
	})

	c.OnResponse(func(resp *colly.Response) {
		html := string(resp.Body)
		homepage := strings.Split(strings.Split(html, `"hotPosts|hotPost":`)[1], `,"highProjects|focus":`)[0]
		hotNewsList := &News{}
		if err := json.Unmarshal([]byte(homepage), hotNewsList); err != nil {
			log.Fatal(err)
		}
		for _, onenews := range *hotNewsList {
			id := onenews.ID
			url := fmt.Sprintf("http://36kr.com/p/%s.html", id)
			ctx := colly.NewContext()
			ctx.Put("title", onenews.Title)
			ctx.Put("coverlink", onenews.Coverlink)
			ctx.Put("authorname", onenews.Author.Name)
			ctx.Put("authoravatar", onenews.Author.Avatar)
			ctx.Put("authorintroduction", onenews.Author.Introduction)
			ctx.Put("publishtime", onenews.Publishtime)
			detailCollector.Request("GET", url, nil, ctx, nil)
			log.Println("Visiting: ", url)
		}
	})

	c.Visit("http://36kr.com/")
}
