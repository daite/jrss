package main

import (
	"context"
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/schollz/progressbar/v3"
)

var Version = "0.1.1"

// Define RSS URLs as constants
const (
	doctorRSSURL = "https://www.omnycontent.com/d/playlist/67122501-9b17-4d77-84bd-a93d00dc791e/3c31cad9-230a-4a5f-b487-a9de001adcdd/1e498682-cfe8-4f7e-adb1-aa5b0019ae1d/podcast.rss"
	cozyRSSURL   = "https://www.omnycontent.com/d/playlist/67122501-9b17-4d77-84bd-a93d00dc791e/3c31cad9-230a-4a5f-b487-a9de001adcdd/39cee2d4-8502-4b84-b11b-a9de001ca4cc/podcast.rss"
)

// RSS structure
type RSS struct {
	Channel struct {
		Items []Item `xml:"item"`
	} `xml:"channel"`
}

// Item represents an individual episode or article in the RSS feed.
type Item struct {
	Title        string         `xml:"title"`
	MediaContent []MediaContent `xml:"http://search.yahoo.com/mrss/ content"`
	PubDate      string         `xml:"pubDate"`
}

// MediaContent represents each <media:content> element.
type MediaContent struct {
	URL  string `xml:"url,attr"`
	Type string `xml:"type,attr"`
}

// getAudioURL extracts the audio URL from an Item
func getAudioURL(item Item) (string, bool) {
	for _, media := range item.MediaContent {
		if media.Type == "audio/mpeg" {
			return media.URL, true
		}
	}
	return "", false
}

// fetchAndDownload downloads the MP3 file using context for timeout handling
func fetchAndDownload(title, audioURL string, wg *sync.WaitGroup) {
	defer wg.Done()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, audioURL, nil)
	if err != nil {
		fmt.Printf("Failed to create request: %v\n", err)
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("Failed to download %s: %v\n", title, err)
		return
	}
	defer resp.Body.Close()

	filename := fmt.Sprintf("%s.mp3", strings.ReplaceAll(title, "/", "_"))
	file, err := os.Create(filename)
	if err != nil {
		fmt.Printf("Failed to create file %s: %v\n", filename, err)
		return
	}
	defer file.Close()

	bar := progressbar.DefaultBytes(resp.ContentLength, fmt.Sprintf("Downloading %s", filename))
	_, err = io.Copy(io.MultiWriter(file, bar), resp.Body)
	if err != nil {
		fmt.Printf("Failed to save %s: %v\n", filename, err)
		return
	}

	// Ensure progress bar finishes
	bar.Finish()

	// Move cursor to the beginning of the line and clear the line
	fmt.Print("\033[1A") // Move cursor up one line
	fmt.Print("\033[K")  // Clear the line

	// Print the download complete message
	fmt.Printf("Download complete: %s\n", filename)
}

func main() {
	numEpisodes := flag.Int("n", 1, "Number of latest episodes to download")
	rssOption := flag.String("rss", "doctor", "Select which RSS feed to use: 'doctor' or 'cozy'")
	showVersion := flag.Bool("version", false, "Show the current version") // ✅ 버전 플래그 추가
	flag.Parse()

	if *showVersion {
		fmt.Printf("version: %s\n", Version)
		return
	}

	var rssURL string
	switch *rssOption {
	case "cozy":
		rssURL = cozyRSSURL
	case "doctor":
		rssURL = doctorRSSURL
	default:
		fmt.Println("Invalid RSS feed option. Please choose 'doctor' or 'cozy'.")
		return
	}

	resp, err := http.Get(rssURL)
	if err != nil {
		fmt.Printf("Failed to fetch RSS feed: %v\n", err)
		return
	}
	defer resp.Body.Close()

	var rss RSS
	if err := xml.NewDecoder(resp.Body).Decode(&rss); err != nil {
		fmt.Printf("Failed to parse RSS feed: %v\n", err)
		return
	}

	sort.SliceStable(rss.Channel.Items, func(i, j int) bool {
		return i < j
	})

	var wg sync.WaitGroup
	for i := 0; i < *numEpisodes && i < len(rss.Channel.Items); i++ {
		item := rss.Channel.Items[i]
		audioURL, found := getAudioURL(item)
		if !found {
			fmt.Printf("No audio found for: %s\n", item.Title)
			continue
		}

		wg.Add(1)
		go fetchAndDownload(item.Title, audioURL, &wg)
	}

	wg.Wait()
}
