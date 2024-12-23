package musicbrainz

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/xaviercrochet/turbo-octo-adventure/pkg/net"
)

// FeedXml represents the main struct of the xml response we get from the API
type FeedXml struct {
	XMLName xml.Name `xml:"feed" json:"-"` // This can be ommited, it contains metadata bout the xml document
	Title   string   `xml:"title" json:"title"`
	ID      string   `xml:"id" json:"id"`
	Updated string   `xml:"updated" json:"updated"`
	Author  Author   `xml:"author" json:"author"`
	Entries []Entry  `xml:"entry" json:"entries"`
}

// Author represents the feed author...
type Author struct {
	Name string `xml:"name" json:"name"`
}

// Entry represents a single listen entry...
type Entry struct {
	ID        string    `xml:"id" json:"id"`
	Title     string    `xml:"title" json:"title"`
	Published time.Time `xml:"published" json:"published"`
	Updated   time.Time `xml:"updated" json:"updated"`
	Content   Content   `xml:"content" json:"content"`
}

// Content represents the listen details...
type Content struct {
	Type string `xml:"type,attr" json:"type"`
	Text string `xml:",chardata" json:"text"`
}

// Integrate the feed api from musicbrainz
func GetFeed(username string) (*Feed, error) {
	// 5000 is the maximum time range the API allows
	reqUrl := fmt.Sprintf("https://listenbrainz.org/syndication-feed/user/%s/listens?minutes=5000", username)

	resp, err := http.Get(reqUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to query feed api: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return &Feed{
			Username: username,
			Songs:    []*Song{},
		}, nil
	}

	if err := net.HttpStatusCodeToErr(resp); err != nil {
		return nil, fmt.Errorf("failed to query musicbrainz api: %v", err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response body: %w", err)
	}

	var feed FeedXml
	err = xml.Unmarshal(body, &feed)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize xml: %w", err)
	}

	// MAP the deserialized XML response to Feed, the struct that will be used later
	feedResp := &Feed{
		Username: username,
	}

	for _, entry := range feed.Entries {
		song := &Song{
			Title:      entry.Title,
			ListenedAt: entry.Updated,
		}

		feedResp.Songs = append(feedResp.Songs, song)
	}

	return feedResp, nil
}

// Feed and Song contain the relevant data for the feed API
type Feed struct {
	Username string  `json:"username"`
	Songs    []*Song `json:"songs"`
}

type Song struct {
	Title      string    `json:"title"`
	ListenedAt time.Time `json:"listened_at"`
}
