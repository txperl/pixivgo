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

	// Create client (optionally with language for translated tags)
	client := pixivgo.NewClient(
		pixivgo.WithAcceptLanguage("en-us"),
	)

	// Authenticate
	authResp, err := client.Auth(ctx, refreshToken)
	if err != nil {
		log.Fatalf("Auth failed: %v", err)
	}
	fmt.Printf("Logged in as: %s (ID: %d)\n", authResp.User.Name, authResp.User.ID)

	// Search illustrations
	result, err := client.SearchIllust(ctx, pixivgo.SearchIllustParams{
		Word: "風景",
		Sort: pixivgo.SortPopularDesc,
	})
	if err != nil {
		log.Fatalf("Search failed: %v", err)
	}

	fmt.Printf("\nFound %d illustrations:\n", len(result.Illusts))
	for i, illust := range result.Illusts {
		if i >= 5 {
			break
		}
		fmt.Printf("  [%d] %s (views: %d, bookmarks: %d)\n",
			illust.ID, illust.Title, illust.TotalView, illust.TotalBookmarks)
	}

	// Pagination — get next page
	nextParams := pixivgo.ParseNextURL(result.NextURL)
	if nextParams != nil {
		fmt.Printf("\nNext page available (offset=%s)\n", nextParams.Get("offset"))
	}

	// Get daily ranking
	ranking, err := client.IllustRanking(ctx, pixivgo.IllustRankingParams{
		Mode: pixivgo.ModeDay,
	})
	if err != nil {
		log.Fatalf("Ranking failed: %v", err)
	}

	fmt.Printf("\nDaily ranking top 3:\n")
	for i, illust := range ranking.Illusts {
		if i >= 3 {
			break
		}
		fmt.Printf("  #%d: %s by %s\n", i+1, illust.Title, illust.User.Name)
	}
}
