package pixivgo

import (
	"context"
	"net/http"
	"testing"
)

// TestIllustRanking is migrated from Python test_app_api.py::TestAppPixivAPI::test_illust_ranking.
// It uses the illust_ranking.json fixture (real API response) to verify response parsing.
func TestIllustRanking(t *testing.T) {
	fixture := loadFixture(t, "illust_ranking.json")

	client, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		// Verify request path
		if r.URL.Path != "/v1/illust/ranking" {
			t.Errorf("path = %s, want /v1/illust/ranking", r.URL.Path)
		}
		// Verify request parameters
		if got := r.URL.Query().Get("mode"); got != "day" {
			t.Errorf("mode = %q, want \"day\"", got)
		}
		if got := r.URL.Query().Get("date"); got != "2025-02-04" {
			t.Errorf("date = %q, want \"2025-02-04\"", got)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(fixture)
	})

	result, err := client.IllustRanking(context.Background(), IllustRankingParams{
		Mode:   ModeDay,
		Date:   String("2025-02-04"),
		NoAuth: false,
	})
	if err != nil {
		t.Fatalf("IllustRanking error: %v", err)
	}

	// Check fields (same assertions as Python test)
	if len(result.Illusts) != 30 {
		t.Fatalf("len(illusts) = %d, want 30", len(result.Illusts))
	}

	illust := result.Illusts[0]
	if illust.ID.Int() != 126839080 {
		t.Errorf("illust.ID = %d, want 126839080", illust.ID.Int())
	}
	if illust.Title != "奏鳴の宙" {
		t.Errorf("illust.Title = %q, want %q", illust.Title, "奏鳴の宙")
	}
	if illust.Caption != "「光る美少女展2024」出展作品" {
		t.Errorf("illust.Caption = %q, want %q", illust.Caption, "「光る美少女展2024」出展作品")
	}
	if illust.ImageUrls.Medium == "" {
		t.Error("illust.ImageUrls.Medium should not be empty")
	}
	if illust.ImageUrls.Large == "" {
		t.Error("illust.ImageUrls.Large should not be empty")
	}
	if illust.MetaSinglePage.OriginalImageURL == nil || *illust.MetaSinglePage.OriginalImageURL == "" {
		t.Error("illust.MetaSinglePage.OriginalImageURL should not be nil/empty")
	}
	if illust.TotalView != 79169 {
		t.Errorf("illust.TotalView = %d, want 79169", illust.TotalView)
	}

	// Check next url (same assertions as Python test using parse_qs)
	qs := ParseNextURL(result.NextURL)
	if qs == nil {
		t.Fatal("ParseNextURL returned nil")
	}
	if got := qs.Get("date"); got != "2025-02-04" {
		t.Errorf("next_url date = %q, want \"2025-02-04\"", got)
	}
	if got := qs.Get("offset"); got != "30" {
		t.Errorf("next_url offset = %q, want \"30\"", got)
	}
}
