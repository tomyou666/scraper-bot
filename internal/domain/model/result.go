package model

import "net/url"

// Result はパイプライン最終出力。P6 Transformer が組み立てる。
type Result struct {
	// URL は結果の対象 URL。
	URL *url.URL
	// Markdown は変換後の Markdown 本文。
	Markdown string
	// HTML は整形済み HTML 本文。
	HTML string
	// RawHTML は取得した生 HTML。
	RawHTML string
	// JSON は構造化出力用の任意フィールド。
	JSON map[string]any
	// Links はページ内から抽出したリンク URL。
	Links []*url.URL
	// Metadata はメタデータの key-value。
	Metadata map[string]string
}
