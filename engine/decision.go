package engine

import (
	"log"
	ptn "github.com/middelink/go-parse-torrent-name"
)

func ParseTorrentName(name string) *ptn.TorrentInfo {
	torrent, err := ptn.Parse(name)
	if err != nil {
		log.Fatalf("Error parsing torrent name: %v", err)
	}
	return torrent
}

func BestTorrent(query string) {}