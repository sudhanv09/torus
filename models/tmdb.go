package models

// Search Results
type TMDBResults struct {
	Adult            bool    `json:"adult"`
	BackdropPath     string  `json:"backdrop_path"`
	ID               int     `json:"id"`
	Title            string  `json:"title"`
	OriginalLanguage string  `json:"original_language"`
	OriginalTitle    string  `json:"original_title"`
	Overview         string  `json:"overview"`
	PosterPath       string  `json:"poster_path"`
	MediaType        string  `json:"media_type"`     // "movie", "tv", or "person"
	GenreIds         []int   `json:"genre_ids"`
	Popularity       float64 `json:"popularity"`
	ReleaseDate      string  `json:"release_date"`
	Video            bool    `json:"video"`
	VoteAverage      float64 `json:"vote_average"`
	VoteCount        int     `json:"vote_count"`
}

type TMDBResponse struct {
	Page         int         `json:"page"`
	Results      []TMDBResults `json:"results"`
	TotalPages   int         `json:"total_pages"`
	TotalResults int         `json:"total_results"`
}

// TMDB Tv Series Details
type TMDBShowDetails struct {
	Adult             bool                    `json:"adult"`
	BackdropPath      string                  `json:"backdrop_path"`
	CreatedBy         []TMDBCreatedBy         `json:"created_by"`
	EpisodeRunTime    []int                   `json:"episode_run_time"`
	FirstAirDate      string                  `json:"first_air_date"`
	Genres            []TMDBGenre             `json:"genres"`
	Homepage          string                  `json:"homepage"`
	ID                int                     `json:"id"`
	InProduction      bool                    `json:"in_production"`
	Languages         []string                `json:"languages"`
	LastAirDate       string                  `json:"last_air_date"`
	LastEpisodeToAir  TMDBEpisode             `json:"last_episode_to_air"`
	Name              string                  `json:"name"`
	NextEpisodeToAir  string                  `json:"next_episode_to_air"`
	Networks          []TMDBNetwork           `json:"networks"`
	NumberOfEpisodes  int                     `json:"number_of_episodes"`
	NumberOfSeasons   int                     `json:"number_of_seasons"`
	OriginCountry     []string                `json:"origin_country"`
	OriginalLanguage  string                  `json:"original_language"`
	OriginalName      string                  `json:"original_name"`
	Overview          string                  `json:"overview"`
	Popularity        float64                 `json:"popularity"`
	PosterPath        string                  `json:"poster_path"`
	ProductionCompany []TMDBProductionCompany `json:"production_companies"`
	ProductionCountry []TMDBProductionCountry `json:"production_countries"`
	Seasons           []TMDBSeason            `json:"seasons"`
	SpokenLanguages   []TMDBSpokenLanguage    `json:"spoken_languages"`
	Status            string                  `json:"status"`
	Tagline           string                  `json:"tagline"`
	Type              string                  `json:"type"`
	VoteAverage       float64                 `json:"vote_average"`
	VoteCount         int                     `json:"vote_count"`
}

type TMDBCreatedBy struct {
	ID          int    `json:"id"`
	CreditID    string `json:"credit_id"`
	Name        string `json:"name"`
	Gender      int    `json:"gender"`
	ProfilePath string `json:"profile_path"`
}

type TMDBGenre struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type TMDBNetwork struct {
	ID            int    `json:"id"`
	LogoPath      string `json:"logo_path"`
	Name          string `json:"name"`
	OriginCountry string `json:"origin_country"`
}

type TMDBProductionCompany struct {
	ID            int    `json:"id"`
	LogoPath      string `json:"logo_path"`
	Name          string `json:"name"`
	OriginCountry string `json:"origin_country"`
}

type TMDBProductionCountry struct {
	Iso3166_1 string `json:"iso_3166_1"`
	Name      string `json:"name"`
}

type TMDBSpokenLanguage struct {
	EnglishName string `json:"english_name"`
	Iso639_1    string `json:"iso_639_1"`
	Name        string `json:"name"`
}

type TMDBSeason struct {
	AirDate      string        `json:"air_date"`
	EpisodeCount int           `json:"episode_count,omitempty"`
	Episodes     []TMDBEpisode `json:"episodes,omitempty"`
	ID           int           `json:"id"`
	Name         string        `json:"name"`
	Overview     string        `json:"overview"`
	PosterPath   string        `json:"poster_path"`
	SeasonNumber int           `json:"season_number"`
	VoteAverage  float64       `json:"vote_average"`
}

type TMDBShowCrew struct {
	Department string `json:"department"`
	Job string `json:"job"`
	CreditID string `json:"credit_id"`
	Adult bool `json:"adult"`
	Gender int `json:"gender"`
	ID int `json:"id"`
	KnownForDepartment string `json:"known_for_department"`
	Name string `json:"name"`
	OriginalName string `json:"original_name"`
	Popularity float64 `json:"popularity"`
	ProfilePath string `json:"profile_path"`
}

type TMDBGuestStar struct {
	Character string `json:"character"`
	CreditID    string `json:"credit_id"`
	Order       int    `json:"order"`
	Adult       bool   `json:"adult"`
	Gender      int    `json:"gender"`
	ID          int    `json:"id"`
	KnownForDepartment string `json:"known_for_department"`
	Name        string `json:"name"`
	OriginalName string `json:"original_name"`
	Popularity float64 `json:"popularity"`
	ProfilePath string `json:"profile_path"`
}

type TMDBEpisode struct {
	ID             int     `json:"id"`
	Name           string  `json:"name"`
	Overview       string  `json:"overview"`
	VoteAverage    float64 `json:"vote_average"`
	VoteCount      int     `json:"vote_count"`
	AirDate        string  `json:"air_date"`
	EpisodeNumber  int     `json:"episode_number"`
	ProductionCode string  `json:"production_code"`
	Runtime        int     `json:"runtime"`
	SeasonNumber   int     `json:"season_number"`
	ShowID         int     `json:"show_id"`
	StillPath      string  `json:"still_path"`
	EpisodeType string `json:"episode_type"`
	Crew []TMDBShowCrew `json:"crew"`
	GuestStars []TMDBGuestStar `json:"guest_stars"`
}
	
// TV Season Details
type TMDBSeasonDetails struct {
	AirDate      string        `json:"air_date"`
	Episodes     []TMDBEpisode `json:"episodes"`
	Name         string        `json:"name"`
	Networks     []TMDBNetwork `json:"networks"`
	Overview     string        `json:"overview"`
	ID           int           `json:"id"`
	PosterPath   string        `json:"poster_path"`
	SeasonNumber int           `json:"season_number"`
	VoteAverage  float64       `json:"vote_average"`
}