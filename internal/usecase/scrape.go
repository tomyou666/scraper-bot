// Package usecase はプレゼンテーション層からカーネルへのシナリオを束ねる。
package usecase

import (
	"context"
	"fmt"
	"net/url"

	"scraperbot/internal/core"
	"scraperbot/internal/domain/model"
)

// Scrape は単一URLの取得→出力までを実行するユースケース。
type Scrape struct {
	// Kernel は初期化済みプラグインを束ねるカーネル。
	Kernel *core.Kernel
	// Fetcher は HTTP 取得実装。
	Fetcher core.Fetcher
}

// NewScrape は単一 URL スクレイプ用ユースケースを構築する。
func NewScrape(k *core.Kernel, f core.Fetcher) *Scrape {
	return &Scrape{Kernel: k, Fetcher: f}
}

// Run は与えられた target URL に対してパイプラインを1回走らせる。
func (s *Scrape) Run(ctx context.Context, target string) (*model.Result, error) {
	u, err := url.Parse(target)
	if err != nil {
		return nil, fmt.Errorf("invalid target url: %w", err)
	}
	req := model.NewRequest(u, 0)

	pipe := core.NewPipeline(s.Kernel, s.Fetcher)
	out, err := pipe.Run(ctx, req)
	if err != nil {
		return nil, err
	}
	return out.Result, nil
}
