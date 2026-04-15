[English](README.md) | 中文

# pixivgo

Pixiv App-API 的 Go 客户端。[pixivpy](https://github.com/upbit/pixivpy) 的 Golang 复写。

## 特性

- 零外部依赖 — 仅使用标准库
- 完整的 API 覆盖 — 插画、用户、小说、搜索、收藏等（约 40 个方法）
- 类型安全的请求参数与响应模型
- 内置图片下载，支持流式写入
- SNI 绕过，适用于受限网络（DNS-over-HTTPS）
- 分页辅助工具，解析 `next_url`
- 构建后线程安全

## 快速开始

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

    // 使用 refresh token 认证
    _, err := client.Auth(ctx, os.Getenv("PIXIV_REFRESH_TOKEN"))
    if err != nil {
        log.Fatal(err)
    }

    // 搜索插画
    result, err := client.SearchIllust(ctx, pixivgo.SearchIllustParams{
        Word: "風景",
        Sort: pixivgo.SortPopularDesc,
    })
    if err != nil {
        log.Fatal(err)
    }

    for _, illust := range result.Illusts {
        fmt.Printf("%s (浏览: %d)\n", illust.Title, illust.TotalView)
    }
}
```

## 认证

Pixiv 需要使用 **refresh token** 进行认证，不再支持密码登录。

获取 refresh token 的工具：

- [gppt](https://github.com/eggplants/get-pixivpy-token) — 基于 Selenium 的 token 提取工具
- [手动 OAuth 流程](https://gist.github.com/ZipFile/c9ebedb224406f4f11845ab700124362)

```go
// 认证并在内部存储 token
authResp, err := client.Auth(ctx, "YOUR_REFRESH_TOKEN")

// 或者直接设置已有的 token
client.SetAuth(accessToken, refreshToken)
```

## 用法

### 排行榜

```go
ranking, err := client.IllustRanking(ctx, pixivgo.IllustRankingParams{
    Mode: pixivgo.ModeDay,
})
for _, illust := range ranking.Illusts {
    fmt.Println(illust.Title, "by", illust.User.Name)
}
```

### 下载

```go
path, err := client.Download(ctx, illust.ImageUrls.Large, &pixivgo.DownloadOptions{
    Path: "./downloads",
})
// 或流式写入到任意 io.Writer
err = client.DownloadToWriter(ctx, illust.ImageUrls.Large, w)
```

### 分页

```go
result, _ := client.IllustRanking(ctx, pixivgo.IllustRankingParams{Mode: pixivgo.ModeDay})

// 获取下一页
nextParams := pixivgo.ParseNextURL(result.NextURL)
if nextParams != nil {
    fmt.Println("next offset:", nextParams.Get("offset"))
}
```

### SNI 绕过

适用于受限网络环境下访问 Pixiv：

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

### 错误处理

所有错误都封装为 `PixivError`，包含 HTTP 状态码、响应头和响应体：

```go
var pe *pixivgo.PixivError
if errors.As(err, &pe) {
    fmt.Println(pe.StatusCode, pe.Body)
}
```

## API 概览

| 分类 | 示例方法                                                  |
| ---- | --------------------------------------------------------- |
| 插画 | `IllustDetail`, `IllustRanking`, `IllustRecommended`, ... |
| 用户 | `UserDetail`, `UserIllusts`, `UserFollowAdd`, ...         |
| 小说 | `NovelDetail`, `NovelSeries`, `WebviewNovel`, ...         |
| 搜索 | `SearchIllust`, `SearchNovel`, `SearchUser`               |
| 收藏 | `UserBookmarksIllust`, `IllustBookmarkAdd`, ...           |
| 其他 | `TrendingTagsIllust`, `ShowcaseArticle`, `UgoiraMetadata` |

完整 API 文档：[GoDoc](https://pkg.go.dev/github.com/txperl/pixivgo)

## 客户端选项

| 选项                       | 说明                                 |
| -------------------------- | ------------------------------------ |
| `WithHTTPClient(hc)`       | 自定义 `http.Client`（代理、绕过等） |
| `WithBaseURL(url)`         | 覆盖 API 基础 URL                    |
| `WithAcceptLanguage(l)`    | 标签翻译语言（如 `"en-us"`）         |
| `WithAdditionalHeaders(h)` | 添加自定义 HTTP 头                   |

## 致谢

本项目是 [@upbit](https://github.com/upbit) 的 [pixivpy](https://github.com/upbit/pixivpy/tree/4f2e9ea7fff6247d9f5bfe5a862e92c5dfe3b6dd) 的 Golang 复写，感谢原作者在 Python 客户端上的出色工作。

## License

Feel free to use, reuse and abuse the code in this project.
