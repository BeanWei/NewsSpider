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

	c.OnResponse(func(resp *colly.Response) {
		html := string(resp.Body)
		homepage := strings.Split(strings.Split(html, `"hotPosts|hotPost":`)[1], `,"highProjects|focus":`)[0]
		hotNewsList := &News{}
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
			//news.Content = ""
		}
		d := c.Clone()
		d.OnResponse(func(r *colly.Response) {
			reg := regexp.MustCompile(`"content":(.*?),"cover"`)
			content := reg.FindString(string(r.Body))
			log.Println(content)
		})
		for _, v := range allnews {
			//v.Content = content
			d.Visit(v.Srcurl)
		}
		// d.Wait()
	})

	c.Visit("http://36kr.com/")
}
