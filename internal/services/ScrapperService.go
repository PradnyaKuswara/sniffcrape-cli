package services

import (
	"io"
	"net/http"

	customeErrors "github.com/PradnyaKuswara/sniffcrape/pkg/errors"
	"github.com/gocolly/colly"
)

type ScrapperService struct {
}

type ScrapeResultResponse struct {
	Title       string `json:"title"`
	Url         string `json:"url"`
	Description string `json:"description"`
	Keywords    string `json:"keywords"`
	Canonical   string `json:"canonical"`
	Robots      string `json:"robots"`
	Viewport    string `json:"viewport"`
	Charset     string `json:"charset"`
	Author      string `json:"author"`

	OgTitle string `json:"og_title"`
	OgDesc  string `json:"og_description"`
	OgImage string `json:"og_image"`
	OgUrl   string `json:"og_url"`

	TwTitle string `json:"twitter_title"`
	TwDesc  string `json:"twitter_description"`
	TwImage string `json:"twitter_image"`
	TwCard  string `json:"twitter_card"`

	H1 []string `json:"h1"`
	H2 []string `json:"h2"`
	H3 []string `json:"h3"`

	Images   []string `json:"images"`
	Links    []string `json:"links"`
	Favicons []string `json:"favicons"`
	Scripts  []string `json:"scripts"`
}

func NewScrapperService() *ScrapperService {
	return &ScrapperService{}
}

func (s *ScrapperService) Scrape(url string) (string, error) {
	resp, err := http.Get(url)

	if err != nil {
		return "", customeErrors.ErrInternalServer
	}

	if resp.StatusCode != http.StatusOK {
		return "", customeErrors.ErrDataNotFound
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", customeErrors.ErrInternalServer
	}
	defer resp.Body.Close()
	return string(body), nil
}

func (s *ScrapperService) ScrapeColly(url string) (*ScrapeResultResponse, error) {
	// Placeholder for future implementation using Colly
	c := colly.NewCollector()

	var (
		title, description, keywords, canonical string
		ogTitle, ogDesc, ogImage, ogUrl         string
		twTitle, twDesc, twImage, twCard        string
		robots, viewport, charset, author       string

		h1s, h2s, h3s []string
		imageUrls     []string
		links         []string
		favicons      []string
		scripts       []string
	)

	// Title
	c.OnHTML("title", func(e *colly.HTMLElement) {
		title = e.Text
	})

	// Meta Tags
	c.OnHTML(`meta[name="description"]`, func(e *colly.HTMLElement) {
		description = e.Attr("content")
	})
	c.OnHTML(`meta[name="keywords"]`, func(e *colly.HTMLElement) {
		keywords = e.Attr("content")
	})
	c.OnHTML(`meta[name="robots"]`, func(e *colly.HTMLElement) {
		robots = e.Attr("content")
	})
	c.OnHTML(`meta[name="viewport"]`, func(e *colly.HTMLElement) {
		viewport = e.Attr("content")
	})
	c.OnHTML(`meta[charset]`, func(e *colly.HTMLElement) {
		charset = e.Attr("charset")
	})
	c.OnHTML(`meta[name="author"]`, func(e *colly.HTMLElement) {
		author = e.Attr("content")
	})

	// Canonical
	c.OnHTML(`link[rel="canonical"]`, func(e *colly.HTMLElement) {
		canonical = e.Attr("href")
	})

	// Open Graph Tags
	c.OnHTML(`meta[property="og:title"]`, func(e *colly.HTMLElement) {
		ogTitle = e.Attr("content")
	})
	c.OnHTML(`meta[property="og:description"]`, func(e *colly.HTMLElement) {
		ogDesc = e.Attr("content")
	})
	c.OnHTML(`meta[property="og:image"]`, func(e *colly.HTMLElement) {
		ogImage = e.Attr("content")
	})
	c.OnHTML(`meta[property="og:url"]`, func(e *colly.HTMLElement) {
		ogUrl = e.Attr("content")
	})

	// Twitter Card Tags
	c.OnHTML(`meta[name="twitter:title"]`, func(e *colly.HTMLElement) {
		twTitle = e.Attr("content")
	})
	c.OnHTML(`meta[name="twitter:description"]`, func(e *colly.HTMLElement) {
		twDesc = e.Attr("content")
	})
	c.OnHTML(`meta[name="twitter:image"]`, func(e *colly.HTMLElement) {
		twImage = e.Attr("content")
	})
	c.OnHTML(`meta[name="twitter:card"]`, func(e *colly.HTMLElement) {
		twCard = e.Attr("content")
	})

	// Headings
	c.OnHTML("h1", func(e *colly.HTMLElement) {
		h1s = append(h1s, e.Text)
	})
	c.OnHTML("h2", func(e *colly.HTMLElement) {
		h2s = append(h2s, e.Text)
	})
	c.OnHTML("h3", func(e *colly.HTMLElement) {
		h3s = append(h3s, e.Text)
	})

	// Images
	c.OnHTML("img", func(e *colly.HTMLElement) {
		src := e.Request.AbsoluteURL(e.Attr("src"))
		if src != "" {
			imageUrls = append(imageUrls, src)
		}
	})

	// Links
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		href := e.Request.AbsoluteURL(e.Attr("href"))
		if href != "" {
			links = append(links, href)
		}
	})

	// Favicons
	c.OnHTML(`link[rel="icon"], link[rel="shortcut icon"], link[rel="apple-touch-icon"]`, func(e *colly.HTMLElement) {
		icon := e.Request.AbsoluteURL(e.Attr("href"))
		if icon != "" {
			favicons = append(favicons, icon)
		}
	})

	// External scripts
	c.OnHTML("script[src]", func(e *colly.HTMLElement) {
		src := e.Request.AbsoluteURL(e.Attr("src"))
		if src != "" {
			scripts = append(scripts, src)
		}
	})

	// Start scraping
	if err := c.Visit(url); err != nil {
		return nil, err
	}

	return &ScrapeResultResponse{
		Title:       title,
		Url:         url,
		Description: description,
		Keywords:    keywords,
		Canonical:   canonical,
		Robots:      robots,
		Viewport:    viewport,
		Charset:     charset,
		Author:      author,

		OgTitle: ogTitle,
		OgDesc:  ogDesc,
		OgImage: ogImage,
		OgUrl:   ogUrl,

		TwTitle: twTitle,
		TwDesc:  twDesc,
		TwImage: twImage,
		TwCard:  twCard,

		H1:     h1s,
		H2:     h2s,
		H3:     h3s,
		Images: imageUrls,
		Links:  links,

		Favicons: favicons,
		Scripts:  scripts,
	}, nil
}
