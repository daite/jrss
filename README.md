# Japanese Podcast MP3 Downloader

`jrss` is a small command‑line tool written in Go for downloading the latest
episodes of supported Japanese podcasts. It saves episodes as MP3 files and
displays a progress bar for each download.

## Prerequisites

- [Go](https://go.dev/) **1.22+**
- Internet connection

## Installation

Clone the repository and build the binary:

```bash
git clone https://github.com/<your-username>/jrss.git
cd jrss
go build -o bin/jrss ./cmd/jrss
# or using the provided Makefile
make build
```

## Usage

Run the downloader, specifying how many recent episodes to fetch and which RSS
feed to use:

```bash
# Download the latest episode from the doctor feed
./bin/jrss -n 1 -rss doctor

# Download the five most recent episodes from the Cozy Up! feed
./bin/jrss -n 5 -rss cozy
```

**Flags**

- `-n` &mdash; Number of latest episodes to download (default `1`)
- `-rss` &mdash; Select which RSS feed to use: `doctor` or `cozy` (default
  `doctor`)
- `-version` &mdash; Print the current version and exit

Each episode is saved to the current directory with a progress bar indicating
download progress.

## Features

- Fetches and parses RSS feeds
- Concurrently downloads podcast episodes
- Displays a progress bar for each download
- Handles XML namespaces for media content
- Customizable number of episodes via command‑line flags

## Background

"モーニングライフアップ 今日の早起きドクター" is a podcast providing valuable health and
wellness insights. This tool simplifies the process of downloading multiple
episodes at once. Using Go's concurrency features makes the tool efficient at
managing several downloads simultaneously.

## License

Licensed under the [MIT License](LICENSE).

