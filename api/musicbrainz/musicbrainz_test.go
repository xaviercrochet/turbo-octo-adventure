package musicbrainz

import (
	"reflect"
	"testing"
	"time"
)

func TestFeedXmlToFeed(t *testing.T) {
	tests := []struct {
		name     string
		username string
		feedXml  FeedXml
		expected *Feed
	}{
		{
			name:     "empty feed",
			username: "Test Username",
			feedXml:  FeedXml{Entries: []Entry{}},
			expected: &Feed{
				Username: "Test Username",
				Songs:    []*Song{},
			},
		},
		{
			name:     "single song",
			username: "Test Username",
			feedXml: FeedXml{
				Entries: []Entry{
					{
						Title:   "Test Song",
						Updated: parseTime("2024-01-01T12:00:00Z"),
					},
				},
			},
			expected: &Feed{
				Username: "Test Username",
				Songs: []*Song{
					{
						Title:      "Test Song",
						ListenedAt: parseTime("2024-01-01T12:00:00Z"),
					},
				},
			},
		},
		{
			name:     "multiple songs",
			username: "Test Username",
			feedXml: FeedXml{
				Entries: []Entry{
					{
						Title:   "Song 1",
						Updated: parseTime("2024-01-01T12:00:00Z"),
					},
					{
						Title:   "Song 2",
						Updated: parseTime("2024-01-02T12:00:00Z"),
					},
					{
						Title:   "Song 3",
						Updated: parseTime("2024-01-03T12:00:00Z"),
					},
				},
			},
			expected: &Feed{
				Username: "Test Username",
				Songs: []*Song{
					{
						Title:      "Song 1",
						ListenedAt: parseTime("2024-01-01T12:00:00Z"),
					},
					{
						Title:      "Song 2",
						ListenedAt: parseTime("2024-01-02T12:00:00Z"),
					},
					{
						Title:      "Song 3",
						ListenedAt: parseTime("2024-01-03T12:00:00Z"),
					},
				},
			},
		},
		{
			name:     "empty username",
			username: "",
			feedXml: FeedXml{
				Entries: []Entry{
					{
						Title:   "Test Song",
						Updated: parseTime("2024-01-01T12:00:00Z"),
					},
				},
			},
			expected: &Feed{
				Username: "",
				Songs: []*Song{
					{
						Title:      "Test Song",
						ListenedAt: parseTime("2024-01-01T12:00:00Z"),
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := feedXmlToFeed(tt.username, tt.feedXml)

			if result.Username != tt.expected.Username {
				t.Errorf("feedXmlToFeed: username = %v, expected %v", result.Username, tt.expected.Username)
			}

			if len(result.Songs) != len(tt.expected.Songs) {
				t.Errorf("feedXmlToFeed: result %v songs, expected %v songs", len(result.Songs), len(tt.expected.Songs))
				return
			}

			for i := range result.Songs {
				// recursively compares Song's attributes.
				if !reflect.DeepEqual(result.Songs[i], tt.expected.Songs[i]) {
					t.Errorf("feedXmlToFeed: song[%d] = %+v, expected %+v", i, result.Songs[i], tt.expected.Songs[i])
				}
			}
		})
	}
}

// generate timestamp out of 2024-01-01T12:00:00Z
func parseTime(date string) time.Time {
	t, _ := time.Parse(time.RFC3339, date)
	return t
}
