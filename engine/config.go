package engine

type QualityProfile struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`

	PreferredResolutions []string `json:"preferred_resolutions"`
	MinResolution        string   `json:"min_resolution"`

	PreferredSources []string `json:"preferred_sources"`
	PreferredCodecs  []string `json:"preferred_codecs"`

	BlockedWords    []string `json:"blocked_words"`
	PreferredGroups []string `json:"preferred_groups"`

	MinSeeders int `json:"min_seeders"`
}

func DefaultQualityProfile() *QualityProfile {
	return &QualityProfile{
		Name:                 "Default",
		PreferredResolutions: []string{"1080p", "2160p"},
		MinResolution:        "720p",
		PreferredSources:     []string{"BluRay", "WEB-DL", "Remux"},
		PreferredCodecs:      []string{"x265", "HEVC"},
		BlockedWords:         []string{"CAM", "HDTS", "TELESYNC", "HDCAM", "TS", "TC"},
		PreferredGroups:      []string{},
		MinSeeders:           3,
	}
}

var (
	resolutionScores = map[string]int{
		"2160p": 100,
		"1080p": 80,
		"720p":  50,
		"480p":  10,
	}

	sourceScores = map[string]int{
		"BluRay": 100,
		"Remux":  95,
		"WEB-DL": 80,
		"WEBRip": 70,
		"HDTV":   50,
		"DVDRip": 30,
		"BDRip":  60,
	}

	codecScores = map[string]int{
		"x265": 20,
		"HEVC": 20,
		"AV1":  25,
		"x264": 10,
	}
)
