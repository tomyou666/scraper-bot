package model

import (
	"net/url"
	"time"
)

// Response は HTTP 取得結果を表す（P5 Parser の入力）。
type Response struct {
	// URL は取得したリソースの URL。
	URL *url.URL
	// StatusCode は HTTP ステータスコード。
	StatusCode int
	// Headers はレスポンスヘッダ（先頭値のみ保持）。
	Headers map[string]string
	// ContentType は Content-Type ヘッダの値。
	ContentType string
	// Body はレスポンス本文の生バイト列。
	Body []byte
	// FetchedAt は取得完了時刻。
	FetchedAt time.Time
}
