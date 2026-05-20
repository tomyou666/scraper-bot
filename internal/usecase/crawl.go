package usecase

import (
	"context"
	"fmt"
	"net/url"

	"scraperbot/internal/core"
	"scraperbot/internal/domain/model"
)

// Crawl はクローリングシナリオを束ねるユースケース。
type Crawl struct {
	// Kernel は初期化済みプラグインを束ねるカーネル。
	Kernel *core.Kernel
	// Fetcher は HTTP 取得実装。
	Fetcher core.Fetcher
	// Robots は robots.txt 判定（nil 可）。
	Robots core.RobotsChecker
	// Sink は各ページの Result 受け取り先（nil 可）。
	Sink core.ResultSink
}

// NewCrawl はクロール用ユースケースを構築する。
//
// robots は nil の場合 robots 判定をスキップする。
// sink は nil の場合、収集した Result は戻り値のスライスにのみ格納される。
func NewCrawl(k *core.Kernel, f core.Fetcher, robots core.RobotsChecker, sink core.ResultSink) *Crawl {
	return &Crawl{Kernel: k, Fetcher: f, Robots: robots, Sink: sink}
}

// Run はシード URL 一覧からクロールを実行し、統計と収集結果を返す。
func (c *Crawl) Run(ctx context.Context, targets []string) (*core.CrawlStats, []*model.Result, error) {
	if len(targets) == 0 {
		return nil, nil, fmt.Errorf("no target URLs")
	}
	seeds := make([]*url.URL, 0, len(targets))
	for _, t := range targets {
		u, err := url.Parse(t)
		if err != nil {
			return nil, nil, fmt.Errorf("invalid seed url %q: %w", t, err)
		}
		seeds = append(seeds, u)
	}

	var collected []*model.Result
	sink := c.Sink
	if sink == nil {
		sink = func(r *model.Result) { collected = append(collected, r) }
	} else {
		original := sink
		sink = func(r *model.Result) {
			collected = append(collected, r)
			original(r)
		}
	}

	crawler := core.NewCrawler(c.Kernel, c.Fetcher, c.Robots, sink)
	stats, err := crawler.Run(ctx, seeds)
	return stats, collected, err
}
