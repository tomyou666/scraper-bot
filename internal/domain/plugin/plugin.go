package plugin

import (
	"context"
	"net/url"
)

// Kind はプラグインのパイプライン上の役割を表す。
type Kind string

const (
	// KindPreProcessor は P2 リクエスト前処理。
	KindPreProcessor Kind = "preprocessor"
	// KindParser は P5 レスポンス解析。
	KindParser Kind = "parser"
	// KindTransformer は P6 結果変換。
	KindTransformer Kind = "transformer"
	// KindFilter は P7 コンテンツ絞り込み。
	KindFilter Kind = "filter"
	// KindLinkExtractor は P8 リンク抽出。
	KindLinkExtractor Kind = "link_extractor"
	// KindFetcher は URL 取得（ページ・robots.txt 等）。
	KindFetcher Kind = "fetcher"
)

// Metadata はプラグインの識別情報を保持する。
type Metadata struct {
	// Name はレジストリ登録名。
	Name string
	// Version はプラグインのセマンティックバージョン。
	Version string
	// Kind はパイプライン上の種別。
	Kind Kind
	// Description は人間向けの短い説明。
	Description string
}

// Plugin は全プラグイン種別が実装するライフサイクルインタフェース。
type Plugin interface {
	Metadata() Metadata
	Init(ctx context.Context, host Host) error
	Close(ctx context.Context) error
}

// Host はプラグインに渡される最小の依存集合。
// グローバル変数を介さずにここから取得する。
type Host interface {
	Logger() Logger
	Config(key string) (string, bool)
	HTTP() HTTPClient
}

// Logger は構造化ログ出力の最小インタフェース。
type Logger interface {
	Debug(msg string, kv ...any)
	Info(msg string, kv ...any)
	Warn(msg string, kv ...any)
	Error(msg string, kv ...any)
}

// HTTPClient はプラグインから利用する HTTP 実行の抽象。
type HTTPClient interface {
	Do(ctx context.Context, req *HTTPRequest) (*HTTPResponse, error)
}

// HTTPRequest は HTTPClient.Do に渡すリクエスト表現。
type HTTPRequest struct {
	// Method は HTTP メソッド。
	Method string
	// URL はリクエスト先の絶対 URL。
	URL *url.URL
	// Headers は付与するリクエストヘッダ。
	Headers map[string]string
	// Body はリクエストボディ（GET では通常 nil）。
	Body []byte
}

// HTTPResponse は HTTPClient.Do の戻り値。
type HTTPResponse struct {
	// StatusCode は HTTP ステータスコード。
	StatusCode int
	// Headers はレスポンスヘッダ。
	Headers map[string]string
	// Body はレスポンス本文。
	Body []byte
}
