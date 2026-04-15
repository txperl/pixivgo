package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/txperl/pixivgo"
)

func main() {
	refreshToken := os.Getenv("PIXIV_REFRESH_TOKEN")
	if refreshToken == "" {
		log.Fatal("Set PIXIV_REFRESH_TOKEN environment variable")
	}

	ctx := context.Background()
	client := pixivgo.NewClient()

	// Authenticate
	if _, err := client.Auth(ctx, refreshToken); err != nil {
		log.Fatalf("Auth failed: %v", err)
	}

	// Get daily ranking
	result, err := client.IllustRanking(ctx, pixivgo.IllustRankingParams{
		Mode: pixivgo.ModeDay,
	})
	if err != nil {
		log.Fatalf("Ranking failed: %v", err)
	}

	// Download first 3 illustrations
	os.MkdirAll("./downloads", 0755)
	for i, illust := range result.Illusts {
		if i >= 3 {
			break
		}
		fmt.Printf("Downloading: %s...\n", illust.Title)
		path, err := client.Download(ctx, illust.ImageUrls.Large, &pixivgo.DownloadOptions{
			Path: "./downloads",
		})
		if err != nil {
			log.Printf("  Error: %v", err)
			continue
		}
		fmt.Printf("  Saved to: %s\n", path)
	}
}
