package model

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/andybalholm/cascadia"
)

// Config は scraperbot 全体の実行設定を表すルート構造体。
type Config struct {
	// Request は HTTP 取得に関する設定。
	Request RequestConfig `yaml:"request"`
	// Content は本文抽出・出力フォーマットに関する設定。
	Content ContentConfig `yaml:"content"`
	// PDF は PDF 取得・解析に関する設定。
	PDF PDFConfig `yaml:"pdf"`
	// Crawl はサイト横断クロールに関する設定。
	Crawl CrawlConfig `yaml:"crawl"`
	// Plugins は使用するプラグイン名の選択。
	Plugins PluginSelection `yaml:"plugins"`
	// Targets は処理対象 URL の一覧。
	Targets []string `yaml:"targets"`
	// Output は結果ファイルの出力先設定。
	Output OutputConfig `yaml:"output"`
}

// RequestConfig は HTTP リクエストのタイムアウト・リトライ・ヘッダを保持する。
type RequestConfig struct {
	// Headers は追加するリクエストヘッダ（キーはそのまま送信される）。
	Headers map[string]string `yaml:"headers"`
	// Timeout は 1 リクエストあたりの最大待ち時間。
	Timeout time.Duration `yaml:"timeout"`
	// RetryCount は失敗時の再試行回数（0 は再試行なし）。
	RetryCount int `yaml:"retry_count"`
	// RetryInterval は再試行の間隔。
	RetryInterval time.Duration `yaml:"retry_interval"`
}

// ContentConfig は HTML 本文の抽出方針と出力フォーマットを保持する。
type ContentConfig struct {
	// Formats は書き出す出力フォーマットの一覧。
	Formats []OutputFormat `yaml:"formats"`
	// OnlyMainContent はメインコンテンツ領域のみを抽出するか。
	OnlyMainContent bool `yaml:"only_main_content"`
	// IncludeTags は抽出対象に含める HTML タグ名。
	IncludeTags []string `yaml:"include_tags"`
	// ExcludeTags は抽出から除外する HTML タグ名。
	ExcludeTags []string `yaml:"exclude_tags"`
	// Selector は本文を絞り込む CSS セレクタ（空なら全体）。
	Selector string `yaml:"selector"`
	// ExtractLinks は結果にリンク一覧を含めるか。
	ExtractLinks bool `yaml:"extract_links"`
	// ExtractMetadata はメタデータ抽出を行うか。
	ExtractMetadata bool `yaml:"extract_metadata"`
}

// PDFConfig は PDF 処理の有効化と解析モードを保持する。
type PDFConfig struct {
	// Enabled は PDF の取得・解析を許可するか。
	Enabled bool `yaml:"enabled"`
	// Mode は PDF 解析モード（PDFParseMode を参照）。
	Mode PDFParseMode `yaml:"mode"`
	// MaxPages は解析する最大ページ数（0 は無制限）。
	MaxPages int `yaml:"max_pages"`
	// Output は PDF からの出力形式（PDFOutput を参照）。
	Output PDFOutput `yaml:"output"`
}

// CrawlConfig は BFS クロールの深度・件数・フィルタを保持する。
type CrawlConfig struct {
	// Enabled は複数 URL のクロールを有効にするか。
	Enabled bool `yaml:"enabled"`
	// MaxDepth はシードからの最大リンク深度。
	MaxDepth int `yaml:"max_depth"`
	// MaxPages は訪問する最大ページ数。
	MaxPages int `yaml:"max_pages"`
	// IncludePaths は許可する URL パスの正規表現（空なら制限なし）。
	IncludePaths []string `yaml:"include_paths"`
	// ExcludePaths は除外する URL パスの正規表現。
	ExcludePaths []string `yaml:"exclude_paths"`
	// AllowExternal は登録ドメイン外へのリンク追跡を許可するか。
	AllowExternal bool `yaml:"allow_external_links"`
	// AllowSubdomains はサブドメインへの追跡を許可するか。
	AllowSubdomains bool `yaml:"allow_subdomains"`
	// RequestDelay は連続リクエスト間の待機時間（>0 のとき並行度は 1 に制限される）。
	RequestDelay time.Duration `yaml:"request_delay"`
	// MaxConcurrency は同時に走るワーカー数。
	MaxConcurrency int `yaml:"max_concurrency"`
	// RespectRobotsTxt は robots.txt に従うか。
	RespectRobotsTxt bool `yaml:"respect_robots_txt"`
}

// PluginSelection はパイプライン各段で使うプラグイン名を保持する。
type PluginSelection struct {
	// PreProcessors は P2 で実行する PreProcessor 名の順序付き一覧。
	PreProcessors []string `yaml:"preprocessors"`
	// Parsers は P5 で登録される Parser 名の一覧。
	Parsers []string `yaml:"parsers"`
	// Transformer は P6 で使う Transformer 名（1 件）。
	Transformer string `yaml:"transformer"`
	// Filters は P7 で実行する Filter 名の順序付き一覧。
	Filters []string `yaml:"filters"`
	// LinkExtractor は P8 で使う LinkExtractor 名（1 件）。
	LinkExtractor string `yaml:"link_extractor"`
}

// OutputConfig は結果ファイルの保存先と命名規則を保持する。
type OutputConfig struct {
	// Dir は出力ディレクトリのパス。
	Dir string `yaml:"dir"`
	// FilePattern はファイル名テンプレート（{seq},{host},{path},{ext} が使える）。
	FilePattern string `yaml:"file_pattern"`
}

// Default は設計書で確定したデフォルト値を適用した Config を返す。
func Default() Config {
	return Config{
		Request: RequestConfig{
			Headers:       map[string]string{},
			Timeout:       60 * time.Second,
			RetryCount:    2,
			RetryInterval: 1 * time.Second,
		},
		Content: ContentConfig{
			Formats:         []OutputFormat{FormatMarkdown},
			OnlyMainContent: true,
			IncludeTags:     []string{},
			ExcludeTags:     []string{"script", "style", "noscript"},
			Selector:        "",
			ExtractLinks:    true,
			ExtractMetadata: true,
		},
		PDF: PDFConfig{
			Enabled:  true,
			Mode:     PDFModeAuto,
			MaxPages: 0,
			Output:   PDFOutputText,
		},
		Crawl: CrawlConfig{
			Enabled:          false,
			MaxDepth:         2,
			MaxPages:         100,
			IncludePaths:     nil,
			ExcludePaths:     nil,
			AllowExternal:    false,
			AllowSubdomains:  false,
			RequestDelay:     0,
			MaxConcurrency:   4,
			RespectRobotsTxt: true,
		},
		Plugins: PluginSelection{
			PreProcessors: nil,
			Parsers:       []string{"html", "pdf"},
			Transformer:   "markdown",
			Filters:       []string{"maincontent"},
			LinkExtractor: "default",
		},
		Output: OutputConfig{
			Dir:         "./out",
			FilePattern: "{seq}-{host}-{path}.{ext}",
		},
	}
}

// Validate は設計書の検証ルールを集中して評価する。
// 違反は errors.Join で集約して返す。
func (c *Config) Validate() error {
	var errs []error

	errs = append(errs, c.validateTargets()...)
	errs = append(errs, c.validateRequest()...)
	errs = append(errs, c.validateContent()...)
	errs = append(errs, c.validatePDF()...)
	errs = append(errs, c.validateCrawl()...)
	errs = append(errs, c.validateOutput()...)

	if c.Crawl.RequestDelay > 0 && c.Crawl.MaxConcurrency != 1 {
		c.Crawl.MaxConcurrency = 1
	}

	return errors.Join(errs...)
}

// validateTargets は targets の件数と URL 形式を検証する。
func (c *Config) validateTargets() []error {
	if len(c.Targets) == 0 {
		return []error{errors.New("targets: 少なくとも1件のURLが必要です")}
	}
	var errs []error
	for i, t := range c.Targets {
		if !strings.HasPrefix(t, "http://") && !strings.HasPrefix(t, "https://") {
			errs = append(errs, fmt.Errorf("targets[%d]: http:// または https:// で始まる必要があります: %q", i, t))
			continue
		}
		if _, err := url.Parse(t); err != nil {
			errs = append(errs, fmt.Errorf("targets[%d]: URLとしてパースできません: %w", i, err))
		}
	}
	return errs
}

// validateRequest は request セクションの数値・ヘッダを検証する。
func (c *Config) validateRequest() []error {
	var errs []error
	if c.Request.Timeout < time.Second || c.Request.Timeout > 300*time.Second {
		errs = append(errs, fmt.Errorf("request.timeout: 1s 以上 300s 以下 (現在: %s)", c.Request.Timeout))
	}
	if c.Request.RetryCount < 0 || c.Request.RetryCount > 10 {
		errs = append(errs, fmt.Errorf("request.retry_count: 0 以上 10 以下 (現在: %d)", c.Request.RetryCount))
	}
	if c.Request.RetryInterval < 100*time.Millisecond || c.Request.RetryInterval > 60*time.Second {
		errs = append(errs, fmt.Errorf("request.retry_interval: 100ms 以上 60s 以下 (現在: %s)", c.Request.RetryInterval))
	}
	for k, v := range c.Request.Headers {
		if strings.TrimSpace(k) == "" || strings.TrimSpace(v) == "" {
			errs = append(errs, fmt.Errorf("request.headers: 空のキーまたは値は許可されません (key=%q)", k))
		}
		if strings.ContainsAny(k, "\r\n") || strings.ContainsAny(v, "\r\n") {
			errs = append(errs, fmt.Errorf("request.headers: 改行を含むヘッダは許可されません (key=%q)", k))
		}
	}
	return errs
}

// validateContent は content セクションのフォーマット・タグ・セレクタを検証する。
func (c *Config) validateContent() []error {
	var errs []error
	seen := map[OutputFormat]bool{}
	for _, f := range c.Content.Formats {
		if !f.Valid() {
			errs = append(errs, fmt.Errorf("content.formats: 不正なフォーマット %q", f))
			continue
		}
		if seen[f] {
			errs = append(errs, fmt.Errorf("content.formats: 重複したフォーマット %q", f))
		}
		seen[f] = true
	}
	incSet := map[string]bool{}
	for _, t := range c.Content.IncludeTags {
		incSet[t] = true
	}
	for _, t := range c.Content.ExcludeTags {
		if incSet[t] {
			errs = append(errs, fmt.Errorf("content.exclude_tags: include_tags と同名タグは指定できません: %q", t))
		}
	}
	if s := strings.TrimSpace(c.Content.Selector); s != "" {
		if _, err := cascadia.Compile(s); err != nil {
			errs = append(errs, fmt.Errorf("content.selector: CSSセレクタとしてパースできません: %w", err))
		}
	}
	return errs
}

// validatePDF は pdf セクションの mode・output・max_pages を検証する。
func (c *Config) validatePDF() []error {
	var errs []error
	if !c.PDF.Mode.Valid() {
		errs = append(errs, fmt.Errorf("pdf.mode: 不正な値 %q", c.PDF.Mode))
	}
	if !c.PDF.Output.Valid() {
		errs = append(errs, fmt.Errorf("pdf.output: 不正な値 %q", c.PDF.Output))
	}
	if c.PDF.MaxPages < 0 || c.PDF.MaxPages > 10000 {
		errs = append(errs, fmt.Errorf("pdf.max_pages: 0 以上 10000 以下 (現在: %d)", c.PDF.MaxPages))
	}
	return errs
}

// validateCrawl は crawl セクションの深度・件数・正規表現を検証する。
func (c *Config) validateCrawl() []error {
	var errs []error
	if c.Crawl.MaxDepth < 0 || c.Crawl.MaxDepth > 10 {
		errs = append(errs, fmt.Errorf("crawl.max_depth: 0 以上 10 以下 (現在: %d)", c.Crawl.MaxDepth))
	}
	if c.Crawl.MaxPages < 1 || c.Crawl.MaxPages > 100000 {
		errs = append(errs, fmt.Errorf("crawl.max_pages: 1 以上 100000 以下 (現在: %d)", c.Crawl.MaxPages))
	}
	if c.Crawl.MaxConcurrency < 1 || c.Crawl.MaxConcurrency > 64 {
		errs = append(errs, fmt.Errorf("crawl.max_concurrency: 1 以上 64 以下 (現在: %d)", c.Crawl.MaxConcurrency))
	}
	if c.Crawl.RequestDelay < 0 || c.Crawl.RequestDelay > 60*time.Second {
		errs = append(errs, fmt.Errorf("crawl.request_delay: 0s 以上 60s 以下 (現在: %s)", c.Crawl.RequestDelay))
	}
	for i, p := range c.Crawl.IncludePaths {
		if _, err := regexp.Compile(p); err != nil {
			errs = append(errs, fmt.Errorf("crawl.include_paths[%d]: 不正な正規表現 %q: %w", i, p, err))
		}
	}
	for i, p := range c.Crawl.ExcludePaths {
		if _, err := regexp.Compile(p); err != nil {
			errs = append(errs, fmt.Errorf("crawl.exclude_paths[%d]: 不正な正規表現 %q: %w", i, p, err))
		}
	}
	return errs
}

var placeholderRe = regexp.MustCompile(`\{([a-zA-Z0-9_]+)\}`)

// validateOutput は output.file_pattern のプレースホルダを検証する。
func (c *Config) validateOutput() []error {
	allowed := map[string]bool{"seq": true, "host": true, "path": true, "ext": true}
	var errs []error
	for _, m := range placeholderRe.FindAllStringSubmatch(c.Output.FilePattern, -1) {
		if !allowed[m[1]] {
			errs = append(errs, fmt.Errorf("output.file_pattern: 未知のプレースホルダ {%s}", m[1]))
		}
	}
	return errs
}
