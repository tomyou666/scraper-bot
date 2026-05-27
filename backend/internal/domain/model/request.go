package model

import "net/url"

// Request はパイプラインで扱うスクレイピングリクエストを表す。
// P2 PreProcessor はこの構造体に対して副作用を加えてよい。
type Request struct {
	// URL は取得対象の絶対 URL。
	URL *url.URL
	// Method は HTTP メソッド（通常は GET）。
	Method string
	// Headers はこのリクエストに付与するヘッダ。
	Headers map[string]string
	// Depth はクロール時のシードからの深度（単一 URL では 0）。
	Depth int
	// Meta はプラグイン間で共有する任意のメタデータ。
	Meta map[string]any
}

// NewRequest は GET メソッドのリクエストを構築する。
//
// u は取得対象 URL。
// depth はクロール深度（シード URL では 0）。
func NewRequest(u *url.URL, depth int) *Request {
	return &Request{
		URL:     u,
		Method:  "GET",
		Headers: map[string]string{},
		Depth:   depth,
		Meta:    map[string]any{},
	}
}
