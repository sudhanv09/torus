package scrapers

type Torrent struct {
	Title     string
	DetailURL string
	Magnet    string
	Seeds     int
	Leeches   int
	Size      string
	Uploader  string
	Source    string
}

type Scraper interface {
	Name() string
	Search(query string) ([]Torrent, error)
	GetMagnet(detailURL string) (string, error)
}

var Registry = make(map[string]Scraper)

func Register(s Scraper) {
	Registry[s.Name()] = s
}

func SearchAll(query string) ([]Torrent, error) {
	var all []Torrent
	for _, s := range Registry {
		results, err := s.Search(query)
		if err != nil {
			// Log error and continue with other scrapers
			continue
		}
		all = append(all, results...)
	}
	return all, nil
}
