package engine

import (
	"sort"
	"strings"

	"sudhanv09/torus/scrapers"

	ptn "github.com/middelink/go-parse-torrent-name"
)

type ScoredTorrent struct {
	Torrent   scrapers.Torrent
	Parsed    *ptn.TorrentInfo
	Score     int
	Rejected  bool
	RejectMsg string
}

func ParseTorrentName(name string) *ptn.TorrentInfo {
	torrent, err := ptn.Parse(name)
	if err != nil {
		return nil
	}
	return torrent
}

func ScoreTorrents(torrents []scrapers.Torrent, profile *QualityProfile) []ScoredTorrent {
	if profile == nil {
		profile = DefaultQualityProfile()
	}

	blockedWords := profile.BlockedWords
	preferredGroups := profile.PreferredGroups
	preferredResolutions := profile.PreferredResolutions
	preferredSources := profile.PreferredSources
	preferredCodecs := profile.PreferredCodecs

	var scored []ScoredTorrent

	for _, t := range torrents {
		parsed := ParseTorrentName(t.Title)
		if parsed == nil {
			continue
		}

		st := ScoredTorrent{
			Torrent: t,
			Parsed:  parsed,
		}

		// Check minimum seeders
		if t.Seeds < profile.MinSeeders {
			st.Rejected = true
			st.RejectMsg = "below minimum seeders"
			scored = append(scored, st)
			continue
		}

		// Check blocked words
		titleLower := strings.ToLower(t.Title)
		blocked := false
		for _, word := range blockedWords {
			if strings.Contains(titleLower, strings.ToLower(word)) {
				st.Rejected = true
				st.RejectMsg = "contains blocked word: " + word
				blocked = true
				break
			}
		}
		if blocked {
			scored = append(scored, st)
			continue
		}

		// Check minimum resolution
		if !meetsMinResolution(parsed.Resolution, profile.MinResolution) {
			st.Rejected = true
			st.RejectMsg = "below minimum resolution"
			scored = append(scored, st)
			continue
		}

		score := 0

		// Resolution score
		if resScore, ok := resolutionScores[parsed.Resolution]; ok {
			score += resScore
			// Bonus if it's one of the preferred resolutions
			for _, pr := range preferredResolutions {
				if parsed.Resolution == pr {
					score += 20
					break
				}
			}
		}

		// Source score
		if srcScore, ok := sourceScores[parsed.Quality]; ok {
			score += srcScore
			// Bonus if it's one of the preferred sources
			for _, ps := range preferredSources {
				if parsed.Quality == ps {
					score += 15
					break
				}
			}
		}

		// Codec score
		if parsed.Codec != "" {
			if codecScore, ok := codecScores[parsed.Codec]; ok {
				score += codecScore
			}
			// Bonus if it's one of the preferred codecs
			for _, pc := range preferredCodecs {
				if strings.EqualFold(parsed.Codec, pc) {
					score += 10
					break
				}
			}
		}

		// Seeder bonus (capped)
		seederBonus := t.Seeds / 2
		if seederBonus > 30 {
			seederBonus = 30
		}
		score += seederBonus

		// Preferred group bonus
		if parsed.Group != "" {
			for _, grp := range preferredGroups {
				if strings.EqualFold(parsed.Group, grp) {
					score += 25
					break
				}
			}
		}

		st.Score = score
		scored = append(scored, st)
	}

	// Sort by score descending, rejected items at the end
	sort.Slice(scored, func(i, j int) bool {
		if scored[i].Rejected != scored[j].Rejected {
			return !scored[i].Rejected
		}
		return scored[i].Score > scored[j].Score
	})

	return scored
}

func BestTorrent(torrents []scrapers.Torrent, profile *QualityProfile) *ScoredTorrent {
	scored := ScoreTorrents(torrents, profile)
	for _, st := range scored {
		if !st.Rejected {
			return &st
		}
	}
	return nil
}

func meetsMinResolution(actual, minimum string) bool {
	order := map[string]int{"2160p": 4, "1080p": 3, "720p": 2, "480p": 1}
	actualRank, okA := order[actual]
	minRank, okM := order[minimum]
	if !okA || !okM {
		return true // If unknown, allow it
	}
	return actualRank >= minRank
}
