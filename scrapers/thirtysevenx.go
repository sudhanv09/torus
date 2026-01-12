package scrapers

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
)

const base37xURL = "https://1337x.st"

type ThirtySevenXScraper struct{}

func (s *ThirtySevenXScraper) Name() string {
	return "1337x"
}

func init() {
	Register(&ThirtySevenXScraper{})
}

func buildSearchURL(query string) string {
	return fmt.Sprintf("%s/search/%s/1/", base37xURL, url.PathEscape(query))
}

func (s *ThirtySevenXScraper) Search(query string) ([]Torrent, error) {
	searchURL := buildSearchURL(query)

	session, err := solve(searchURL)
	if err != nil {
		return nil, err
	}

	collector := colly.NewCollector(
		colly.AllowedDomains("1337x.st", "www.1337x.st"),
	)

	if session != nil && session.UserAgent != "" {
		collector.UserAgent = session.UserAgent
	}

	if session != nil && len(session.Cookies) > 0 {
		var cookies []*http.Cookie
		for name, value := range session.Cookies {
			cookies = append(cookies, &http.Cookie{
				Name:   name,
				Value:  value,
				Domain: "1337x.st",
				Path:   "/",
			})
		}
		collector.SetCookies(base37xURL, cookies)
	}

	var torrents []Torrent

	collector.OnHTML("table.table-list tbody tr", func(e *colly.HTMLElement) {
		title := strings.TrimSpace(e.ChildText("td.coll-1.name a:nth-of-type(2)"))
		if title == "" {
			return
		}

		sizeText := strings.TrimSpace(e.ChildText("td.coll-4"))
		size := sizeText
		parts := strings.Fields(sizeText)
		if len(parts) >= 2 {
			size = parts[0] + " " + parts[1]
		}

		torrent := Torrent{
			Title:     title,
			DetailURL: strings.TrimSpace(e.Request.AbsoluteURL(e.ChildAttr("td.coll-1.name a:nth-of-type(2)", "href"))),
			Seeds:     parseInt(e.ChildText("td.coll-2")),
			Leeches:   parseInt(e.ChildText("td.coll-3")),
			Size:      size,
			Uploader:  strings.TrimSpace(e.ChildText("td.coll-5")),
			Source:    "1337x",
		}

		torrents = append(torrents, torrent)
	})

	if err := collector.Visit(searchURL); err != nil {
		return nil, err
	}

	return torrents, nil
}

func (s *ThirtySevenXScraper) GetMagnet(detailURL string) (string, error) {
	session, err := solve(detailURL)
	if err != nil {
		return "", err
	}

	collector := colly.NewCollector(
		colly.AllowedDomains("1337x.st", "www.1337x.st"),
	)

	if session != nil && session.UserAgent != "" {
		collector.UserAgent = session.UserAgent
	}

	if session != nil && len(session.Cookies) > 0 {
		var cookies []*http.Cookie
		for name, value := range session.Cookies {
			cookies = append(cookies, &http.Cookie{
				Name:   name,
				Value:  value,
				Domain: "1337x.st",
				Path:   "/",
			})
		}
		collector.SetCookies(base37xURL, cookies)
	}

	var magnet string
	collector.OnHTML("a[href^='magnet:?']", func(e *colly.HTMLElement) {
		magnet = e.Attr("href")
	})

	if err := collector.Visit(detailURL); err != nil {
		return "", err
	}

	if magnet == "" {
		return "", fmt.Errorf("magnet link not found")
	}

	return magnet, nil
}

func parseInt(val string) int {
	clean := strings.ReplaceAll(strings.TrimSpace(val), ",", "")
	n, _ := strconv.Atoi(clean)
	return n
}
