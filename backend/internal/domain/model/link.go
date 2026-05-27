package model

import "net/url"

// Link はクロール時の追跡対象URLを表す。
type Link struct {
	// URL は追跡対象の絶対 URL。
	URL *url.URL
	// Depth はシードからのリンク深度。
	Depth int
}
