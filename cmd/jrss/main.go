package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strings"
	"sync"

	"github.com/schollz/progressbar/v3"
)

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

// Function to fetch and download MP3 files
func fetchAndDownload(item Item, wg *sync.WaitGroup) {
	defer wg.Done()

	// Find the media content with type="audio/mpeg"
	var audioURL string
	for _, media := range item.MediaContent {
		if media.Type == "audio/mpeg" {
			audioURL = media.URL
			break
		}
	}

	if audioURL == "" {
		fmt.Printf("No audio URL found for item: %s\n", item.Title)
		return
	}

	// Clean the title and add ".mp3" as the filename
	filename := fmt.Sprintf("%s.mp3", strings.ReplaceAll(item.Title, "/", "_"))

	// Request the MP3 file
	resp, err := http.Get(audioURL)
	if err != nil {
		fmt.Printf("Failed to download %s: %v\n", filename, err)
		return
	}
	defer resp.Body.Close()

	// Create the file
	file, err := os.Create(filename)
	if err != nil {
		fmt.Printf("Failed to create file %s: %v\n", filename, err)
		return
	}
	defer file.Close()

	// Set up a progress bar
	bar := progressbar.DefaultBytes(
		resp.ContentLength,
		fmt.Sprintf("Downloading %s", filename),
	)

	// Write the content to the file with a progress bar
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
	// Command-line flags for limiting the number of episodes and selecting the RSS feed
	numEpisodes := flag.Int("n", 1, "Number of latest episodes to download")
	rssOption := flag.String("rss", "doctor", "Select which RSS feed to use: 'doctor' or 'cozy'")
	flag.Parse()

	// Determine the RSS feed URL based on the selected option
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

	// Fetch the RSS feed
	resp, err := http.Get(rssURL)
	if err != nil {
		fmt.Printf("Failed to fetch RSS feed: %v\n", err)
		return
	}
	defer resp.Body.Close()

	// Parse the RSS feed
	var rss RSS
	err = xml.NewDecoder(resp.Body).Decode(&rss)
	if err != nil {
		fmt.Printf("Failed to parse RSS feed: %v\n", err)
		return
	}

	sort.SliceStable(rss.Channel.Items, func(i, j int) bool {
		return i < j
	})

	// Create a wait group to manage concurrency
	var wg sync.WaitGroup

	// Iterate over the latest `n` items in the feed
	for i := 0; i < *numEpisodes && i < len(rss.Channel.Items); i++ {
		item := rss.Channel.Items[i]

		// Check if the item has any media content with type="audio/mpeg"
		hasAudio := false
		for _, media := range item.MediaContent {
			if media.Type == "audio/mpeg" {
				hasAudio = true
				break
			}
		}

		if hasAudio {
			wg.Add(1)
			// Download the MP3 file concurrently
			go fetchAndDownload(item, &wg)
		}
	}

	// Wait for all downloads to finish
	wg.Wait()
}
