package pixivgo

import "strings"

func defaultFilter(f Filter) string {
	if f == "" {
		return string(FilterForIOS)
	}
	if f == FilterNone {
		return ""
	}
	return string(f)
}

func formatBool(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

func joinStrings(s []string) string {
	return strings.Join(s, ",")
}
