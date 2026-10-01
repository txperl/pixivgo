package pixivgo

import (
	"context"
	"net/http"
	"reflect"
	"testing"
)

func TestIllustBookmarkAddTags(t *testing.T) {
	for _, tc := range []struct {
		name    string
		tags    []string
		present bool
		encoded string
	}{
		{"omitted", nil, false, ""},
		{"cleared", []string{}, true, ""},
		{"unicode", []string{"風景", "猫&犬", "C++"}, true, "風景 猫&犬 C++"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/v2/illust/bookmark/add" || r.Method != http.MethodPost {
					t.Errorf("unexpected request %s %s", r.Method, r.URL.Path)
				}
				if err := r.ParseForm(); err != nil {
					t.Fatal(err)
				}
				_, present := r.PostForm["tags[]"]
				if present != tc.present || r.PostForm.Get("tags[]") != tc.encoded {
					t.Errorf("tags = %v, want present=%v value=%q", r.PostForm["tags[]"], tc.present, tc.encoded)
				}
				if !reflect.DeepEqual(r.PostForm["restrict"], []string{"private"}) {
					t.Errorf("restrict = %v", r.PostForm["restrict"])
				}
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{}`))
			})
			if err := c.IllustBookmarkAdd(context.Background(), IllustBookmarkAddParams{IllustID: 42, Restrict: RestrictPrivate, Tags: tc.tags}); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestIllustBookmarkDetailRegisteredTags(t *testing.T) {
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("illust_id") != "42" {
			t.Errorf("query = %v", r.URL.Query())
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"bookmark_detail":{"is_bookmarked":true,"restrict":"private","tags":[{"name":"猫","is_registered":true},{"name":"suggestion","is_registered":false}]}}`))
	})
	got, err := c.IllustBookmarkDetail(context.Background(), IllustBookmarkDetailParams{IllustID: 42})
	if err != nil {
		t.Fatal(err)
	}
	if len(got.BookmarkDetail.Tags) != 2 || !got.BookmarkDetail.Tags[0].IsRegistered || got.BookmarkDetail.Tags[1].IsRegistered {
		t.Fatalf("detail = %+v", got.BookmarkDetail)
	}
}
