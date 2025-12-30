package scrapers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

const baseflareSolverrURL = "http://localhost:8191/v1"

type flareSolveRequest struct {
	Cmd        string `json:"cmd"`
	URL        string `json:"url"`
	MaxTimeout int    `json:"maxTimeout"`
}

type flareSolveResponse struct {
	Status   string `json:"status"`
	Message  string `json:"message"`
	Solution struct {
		UserAgent string `json:"userAgent"`
		Cookies   []struct {
			Name   string `json:"name"`
			Value  string `json:"value"`
			Domain string `json:"domain"`
			Path   string `json:"path"`
		} `json:"cookies"`
	} `json:"solution"`
}

type SolvedSession struct {
	UserAgent string
	Cookies   map[string]string
}

func solve(urlStr string) (*SolvedSession, error) {
	reqBody := flareSolveRequest{
		Cmd:        "request.get",
		URL:        urlStr,
		MaxTimeout: 60000,
	}

	b, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	client := &http.Client{Timeout: 90 * time.Second}
	resp, err := client.Post(baseflareSolverrURL, "application/json", bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var fsResp flareSolveResponse
	if err := json.NewDecoder(resp.Body).Decode(&fsResp); err != nil {
		return nil, err
	}

	if fsResp.Status != "ok" {
		return nil, errors.New(fsResp.Message)
	}

	cookies := make(map[string]string)
	for _, c := range fsResp.Solution.Cookies {
		cookies[c.Name] = c.Value
	}

	return &SolvedSession{
		UserAgent: fsResp.Solution.UserAgent,
		Cookies:   cookies,
	}, nil
}
