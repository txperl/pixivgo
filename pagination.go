package pixivgo

import (
	"net/url"
	"strings"
)

// ParseNextURL extracts query parameters from a Pixiv next_url pagination string.
// It handles PHP-style array parameters (e.g., "seed_illust_ids[]" → "seed_illust_ids").
// Returns nil if nextURL is nil or empty.
func ParseNextURL(nextURL *string) url.Values {
	if nextURL == nil || *nextURL == "" {
		return nil
	}

	parsed, err := url.Parse(*nextURL)
	if err != nil {
		return nil
	}

	result := url.Values{}
	for key, values := range parsed.Query() {
		// Merge PHP-style array params: seed_illust_ids[] → seed_illust_ids
		if strings.Contains(key, "[") && strings.HasSuffix(key, "]") {
			cleanKey := key[:strings.Index(key, "[")]
			for _, v := range values {
				result.Add(cleanKey, v)
			}
		} else {
			// For non-array params, take the last value (matching Python behavior)
			result.Set(key, values[len(values)-1])
		}
	}

	return result
}
