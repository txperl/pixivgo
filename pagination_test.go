package pixivgo

import (
	"testing"
)

func TestParseNextURL_Nil(t *testing.T) {
	result := ParseNextURL(nil)
	if result != nil {
		t.Errorf("got %v, want nil", result)
	}
}

func TestParseNextURL_Empty(t *testing.T) {
	result := ParseNextURL(String(""))
	if result != nil {
		t.Errorf("got %v, want nil", result)
	}
}

func TestParseNextURL_Simple(t *testing.T) {
	nextURL := "https://app-api.pixiv.net/v1/illust/ranking?offset=30&filter=for_ios"
	result := ParseNextURL(&nextURL)

	if result == nil {
		t.Fatal("got nil, want url.Values")
	}
	if got := result.Get("offset"); got != "30" {
		t.Errorf("offset = %q, want \"30\"", got)
	}
	if got := result.Get("filter"); got != "for_ios" {
		t.Errorf("filter = %q, want \"for_ios\"", got)
	}
}

func TestParseNextURL_PHPArrayParams(t *testing.T) {
	nextURL := "https://app-api.pixiv.net/v2/illust/related?seed_illust_ids[]=1&seed_illust_ids[]=2&seed_illust_ids[]=3"
	result := ParseNextURL(&nextURL)

	if result == nil {
		t.Fatal("got nil, want url.Values")
	}

	ids := result["seed_illust_ids"]
	if len(ids) != 3 {
		t.Fatalf("seed_illust_ids has %d values, want 3", len(ids))
	}
	expected := []string{"1", "2", "3"}
	for i, want := range expected {
		if ids[i] != want {
			t.Errorf("seed_illust_ids[%d] = %q, want %q", i, ids[i], want)
		}
	}
}

func TestParseNextURL_MixedParams(t *testing.T) {
	nextURL := "https://app-api.pixiv.net/v2/illust/related?illust_id=12345&filter=for_ios&seed_illust_ids[]=100&seed_illust_ids[]=200"
	result := ParseNextURL(&nextURL)

	if result == nil {
		t.Fatal("got nil, want url.Values")
	}
	if got := result.Get("illust_id"); got != "12345" {
		t.Errorf("illust_id = %q, want \"12345\"", got)
	}
	if got := result.Get("filter"); got != "for_ios" {
		t.Errorf("filter = %q, want \"for_ios\"", got)
	}
	ids := result["seed_illust_ids"]
	if len(ids) != 2 {
		t.Fatalf("seed_illust_ids has %d values, want 2", len(ids))
	}
}

func TestParseNextURL_InvalidURL(t *testing.T) {
	bad := "://bad"
	result := ParseNextURL(&bad)
	if result != nil {
		t.Errorf("got %v, want nil for invalid URL", result)
	}
}

func TestParseNextURL_RealPixivURL(t *testing.T) {
	// Real next_url from illust_ranking API response
	nextURL := "https://app-api.pixiv.net/v1/illust/ranking?filter=for_ios&mode=day&date=2025-02-04&offset=30"
	result := ParseNextURL(&nextURL)

	if result == nil {
		t.Fatal("got nil, want url.Values")
	}
	if got := result.Get("date"); got != "2025-02-04" {
		t.Errorf("date = %q, want \"2025-02-04\"", got)
	}
	if got := result.Get("offset"); got != "30" {
		t.Errorf("offset = %q, want \"30\"", got)
	}
	if got := result.Get("mode"); got != "day" {
		t.Errorf("mode = %q, want \"day\"", got)
	}
}
