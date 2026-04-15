English | [中文](README_zh.md)

# pixivgo

A Go client for the Pixiv App-API. Golang rewrite of [pixivpy](https://github.com/upbit/pixivpy).

## Features

- Zero external dependencies — standard library only
- Full API coverage — illustrations, users, novels, search, bookmarks, and more (~40 methods)
- Type-safe request parameters and response models
- Built-in image download with streaming support
- SNI bypass for restricted networks (DNS-over-HTTPS)
- Pagination helper for `next_url` parsing
- Thread-safe after construction

## Quick Start

```bash
go get github.com/txperl/pixivgo
```

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"

    "github.com/txperl/pixivgo"
)

func main() {
    ctx := context.Background()
    client := pixivgo.NewClient()

    // Authenticate with refresh token
    _, err := client.Auth(ctx, os.Getenv("PIXIV_REFRESH_TOKEN"))
    if err != nil {
        log.Fatal(err)
    }

    // Search illustrations
    result, err := client.SearchIllust(ctx, pixivgo.SearchIllustParams{
        Word: "風景",
        Sort: pixivgo.SortPopularDesc,
    })
    if err != nil {
        log.Fatal(err)
    }

    for _, illust := range result.Illusts {
        fmt.Printf("%s (views: %d)\n", illust.Title, illust.TotalView)
    }
}
```

## Authentication

Pixiv requires a **refresh token** for authentication. Password login is no longer supported.

To obtain a refresh token, use one of these tools:

- [gppt](https://github.com/eggplants/get-pixivpy-token) — Selenium-based token extractor
- [Manual OAuth flow](https://gist.github.com/ZipFile/c9ebedb224406f4f11845ab700124362)

```go
// Authenticate and store tokens internally
authResp, err := client.Auth(ctx, "YOUR_REFRESH_TOKEN")

// Or set tokens directly if you already have them
client.SetAuth(accessToken, refreshToken)
```

## Usage

### Ranking

```go
ranking, err := client.IllustRanking(ctx, pixivgo.IllustRankingParams{
    Mode: pixivgo.ModeDay,
})
for _, illust := range ranking.Illusts {
    fmt.Println(illust.Title, "by", illust.User.Name)
}
```

### Download

```go
path, err := client.Download(ctx, illust.ImageUrls.Large, &pixivgo.DownloadOptions{
    Path: "./downloads",
})
// or stream to any io.Writer
err = client.DownloadToWriter(ctx, illust.ImageUrls.Large, w)
```

### Pagination

```go
result, _ := client.IllustRanking(ctx, pixivgo.IllustRankingParams{Mode: pixivgo.ModeDay})

// Get next page
nextParams := pixivgo.ParseNextURL(result.NextURL)
if nextParams != nil {
    fmt.Println("next offset:", nextParams.Get("offset"))
}
```

### SNI Bypass

For accessing Pixiv from restricted networks:

```go
import "github.com/txperl/pixivgo/bypass"

httpClient, hosts, err := bypass.NewHTTPClient(ctx)
if err != nil {
    log.Fatal(err)
}
client := pixivgo.NewClient(
    pixivgo.WithHTTPClient(httpClient),
    pixivgo.WithBaseURL(hosts),
)
```

### Error Handling

All errors are wrapped in `PixivError`, which includes HTTP status, headers, and response body:

```go
var pe *pixivgo.PixivError
if errors.As(err, &pe) {
    fmt.Println(pe.StatusCode, pe.Body)
}
```

## API Overview

| Category | Example Methods                                           |
| -------- | --------------------------------------------------------- |
| Illust   | `IllustDetail`, `IllustRanking`, `IllustRecommended`, ... |
| User     | `UserDetail`, `UserIllusts`, `UserFollowAdd`, ...         |
| Novel    | `NovelDetail`, `NovelSeries`, `WebviewNovel`, ...         |
| Search   | `SearchIllust`, `SearchNovel`, `SearchUser`               |
| Bookmark | `UserBookmarksIllust`, `IllustBookmarkAdd`, ...           |
| Misc     | `TrendingTagsIllust`, `ShowcaseArticle`, `UgoiraMetadata` |

Full API reference: [GoDoc](https://pkg.go.dev/github.com/txperl/pixivgo)

## Client Options

| Option                     | Description                                      |
| -------------------------- | ------------------------------------------------ |
| `WithHTTPClient(hc)`       | Custom `http.Client` (for proxies, bypass, etc.) |
| `WithBaseURL(url)`         | Override API base URL                            |
| `WithAcceptLanguage(l)`    | Language for tag translations (e.g. `"en-us"`)   |
| `WithAdditionalHeaders(h)` | Add custom HTTP headers                          |

## Acknowledgments

This project is a Golang rewrite of [pixivpy](https://github.com/upbit/pixivpy/tree/4f2e9ea7fff6247d9f5bfe5a862e92c5dfe3b6dd) by [@upbit](https://github.com/upbit). Thank you for the excellent work on the original Python client.

## License

Feel free to use, reuse and abuse the code in this project.
