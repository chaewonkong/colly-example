package modules

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly"
	jsoniter "github.com/json-iterator/go"
)

type Site struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Author      string   `json:"author"`
	Keywords    []string `json:"keywords"`
	Image       string   `json:"image"`
	Url         string   `json:"url"`
}

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func Crawler(url string) {
	c := colly.NewCollector(

		// allow domains
		colly.AllowedDomains(strings.Replace(url, "https://", "", 1)),

		// cashe res to prevent multiple downloads
		// colly.CacheDir("./cashe"),
	)

	// collector for site details
	detailCollector := c.Clone()

	sites := make([]Site, 0, 200)

	// print url when making a request
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("visiting:", r.URL.String())
	})

	// recursively visit via anchor tag href attributes
	// c.OnHTML("a[href]", func(h *colly.HTMLElement) {
	// 	targetUrl := h.Request.AbsoluteURL(h.Attr("href"))
	// 	detailCollector.Visit(targetUrl)
	// })
	c.OnHTML("a[href]", func(h *colly.HTMLElement) {
		targetUrl := h.Attr("href")
		h.Request.Visit(targetUrl)
	})

	// // visit pages with detailCollector
	c.OnHTML("a[href]", func(h *colly.HTMLElement) {
		targetUrl := h.Request.AbsoluteURL(h.Attr("href"))
		detailCollector.Visit(targetUrl)
	})

	// extract data
	detailCollector.OnHTML("head", func(h *colly.HTMLElement) {

		site := Site{}

		site.Url = url
		site.Title = h.ChildText("title")
		site.Author = h.ChildAttr("meta[property=author]", "content")

		site.Description = h.ChildAttr("meta[property=description]", "content")
		site.Keywords = strings.Split(h.ChildAttr("meta[property=keywords]", "content"), ",")
		site.Image = h.ChildAttr("meta[property=og:image]", "content")

		sites = append(sites, site)
	})

	c.Visit(url)

	// marshal json
	res, err := json.MarshalToString(sites[:3])
	if err != nil {
		panic("omg")
	}

	fmt.Println(res)
}

func Scraper(url string) {
	c := colly.NewCollector()
	site := Site{}

	// callback
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting")
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println(err)
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println("visited", r.Request.URL)
	})

	c.OnHTML("meta", func(h *colly.HTMLElement) {
		site.Url = url
		property := h.Attr("property")
		content := h.Attr("content")
		switch property {
		case "description":
			site.Description = content
		case "author":
			site.Author = content
		case "og:image":
			site.Image = content
		case "keywords":
			site.Keywords = strings.Split(content, ",")
		default:
			break
		}
	})

	c.OnHTML("title", func(h *colly.HTMLElement) {
		site.Title = h.Text
		fmt.Println("title", site.Title)
	})

	// visit
	c.Visit(url)

	// marshal json
	res, err := json.MarshalToString(site)
	if err != nil {
		return
	}

	fmt.Println(res)
}
