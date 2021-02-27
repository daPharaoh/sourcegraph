package query

import "strings"

func ParseVisibility(s string) repoPrivacy {
	switch strings.ToLower(s) {
	case "private":
		return Private
	case "public":
		return Public
	default:
		return AnyPrivacy
	}
}
