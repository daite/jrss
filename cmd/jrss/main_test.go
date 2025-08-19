package main

import (
	"testing"
	"time"
)

func TestParsePubDate(t *testing.T) {
	rfc1123zTime, _ := time.Parse(time.RFC1123Z, "Mon, 02 Jan 2006 15:04:05 -0700")
	rfc1123Time, _ := time.Parse(time.RFC1123, "Mon, 02 Jan 2006 15:04:05 MST")
	tests := []struct {
		name     string
		input    string
		expected time.Time
	}{
		{"rfc1123z", "Mon, 02 Jan 2006 15:04:05 -0700", rfc1123zTime},
		{"rfc1123", "Mon, 02 Jan 2006 15:04:05 MST", rfc1123Time},
		{"invalid", "not a date", time.Time{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parsePubDate(tt.input)
			if !got.Equal(tt.expected) {
				t.Fatalf("expected %v, got %v", tt.expected, got)
			}
		})
	}
}

func TestGetAudioURL(t *testing.T) {
	item := Item{MediaContent: []MediaContent{{URL: "http://example.com/a.mp3", Type: "audio/mpeg"}}}
	url, ok := getAudioURL(item)
	if !ok {
		t.Fatalf("expected to find audio URL")
	}
	if url != "http://example.com/a.mp3" {
		t.Fatalf("unexpected URL: %s", url)
	}

	item = Item{MediaContent: []MediaContent{{URL: "http://example.com/video.mp4", Type: "video/mp4"}}}
	if url, ok := getAudioURL(item); ok {
		t.Fatalf("unexpected audio URL found: %s", url)
	}
}

func FuzzParsePubDate(f *testing.F) {
	f.Add("Mon, 02 Jan 2006 15:04:05 -0700")
	f.Add("Mon, 02 Jan 2006 15:04:05 MST")
	f.Add("not a date")

	f.Fuzz(func(t *testing.T, date string) {
		got := parsePubDate(date)
		if want, err := time.Parse(time.RFC1123Z, date); err == nil {
			if !got.Equal(want) {
				t.Fatalf("expected %v, got %v", want, got)
			}
		} else if want, err := time.Parse(time.RFC1123, date); err == nil {
			if !got.Equal(want) {
				t.Fatalf("expected %v, got %v", want, got)
			}
		} else {
			if !got.IsZero() {
				t.Fatalf("expected zero time for %q", date)
			}
		}
	})
}

func FuzzGetAudioURL(f *testing.F) {
	f.Add("http://example.com/a.mp3", "audio/mpeg")
	f.Add("", "audio/mpeg")
	f.Add("http://example.com/a.mp3", "video/mp4")

	f.Fuzz(func(t *testing.T, url, mediaType string) {
		item := Item{MediaContent: []MediaContent{{URL: url, Type: mediaType}}}
		gotURL, ok := getAudioURL(item)
		if mediaType == "audio/mpeg" {
			if !ok || gotURL != url {
				t.Fatalf("expected %q, got %q (ok=%v)", url, gotURL, ok)
			}
		} else if ok {
			t.Fatalf("unexpected URL %q", gotURL)
		}
	})
}
