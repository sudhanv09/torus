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

type TorrentRow struct {
	Title     string
	DetailURL string
	Seeds     int
	Leeches   int
	Date      string
	Size      string
	Uploader  string
}

func buildSearchURL(query string) string {
	return fmt.Sprintf("%s/search/%s/1/", base37xURL, url.PathEscape(query))
}

func Search(query string) ([]TorrentRow, error) {
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

	var rows []TorrentRow

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

		row := TorrentRow{
			Title:     title,
			DetailURL: strings.TrimSpace(e.Request.AbsoluteURL(e.ChildAttr("td.coll-1.name a:nth-of-type(2)", "href"))),
			Seeds:     parseInt(e.ChildText("td.coll-2")),
			Leeches:   parseInt(e.ChildText("td.coll-3")),
			Date:      strings.TrimSpace(e.ChildText("td.coll-date")),
			Size:      size,
			Uploader:  strings.TrimSpace(e.ChildText("td.coll-5")),
		}

		rows = append(rows, row)
	})

	if err := collector.Visit(searchURL); err != nil {
		return nil, err
	}

	return rows, nil
}

func parseInt(val string) int {
	clean := strings.ReplaceAll(strings.TrimSpace(val), ",", "")
	n, _ := strconv.Atoi(clean)
	return n
}
