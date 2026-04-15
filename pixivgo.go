// Package pixivgo provides an unofficial Go client for the Pixiv App-API (6.x).
//
// It supports illustration/novel/user queries, search, bookmarks, rankings,
// image downloads, and optional SNI bypass for restricted networks.
//
// Basic usage:
//
//	client := pixivgo.NewClient()
//	resp, err := client.Auth(ctx, "your_refresh_token")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	result, err := client.SearchIllust(ctx, pixivgo.SearchIllustParams{
//	    Word: "landscape",
//	})
package pixivgo

// Int returns a pointer to the given int value.
// Useful for setting optional *int fields in parameter structs.
func Int(v int) *int { return &v }

// String returns a pointer to the given string value.
// Useful for setting optional *string fields in parameter structs.
func String(v string) *string { return &v }

// Bool returns a pointer to the given bool value.
func Bool(v bool) *bool { return &v }
