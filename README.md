# モーニングライフアップ 今日の早起きドクター Podcast MP3 Downloader

This Go program downloads the latest episodes in MP3 format from the RSS feed of the "モーニングライフアップ 今日の早起きドクター" podcast. 
It supports concurrent downloads and features a progress bar for each file being downloaded.


## Features

- Fetches and parses RSS feeds.
- Concurrently downloads the latest podcast episodes.
- Displays a progress bar for each download.
- Handles XML namespaces for media content.
- Customizable number of episodes to download via command-line flags.

## Background

"モーニングライフアップ 今日の早起きドクター" is a podcast providing valuable health and wellness insights. This tool simplifies the process of downloading multiple episodes at once, ensuring you never miss out on the latest content. The use of Go’s concurrency features makes the tool efficient for managing multiple downloads simultaneously.
