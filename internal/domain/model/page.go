package model

import "net/url"

// Page は将来の出力ストレージ層が使うエンティティ。MVPでは Result を主に使う。
type Page struct {
	// URL はページの正規 URL。
	URL *url.URL
	// Title はページタイトル。
	Title string
	// Metadata は抽出したメタデータの key-value。
	Metadata map[string]string
}
