// Package chromiumplugin は chromedp による URL 取得 Fetcher を登録する。
package chromiumplugin

import (
	"scraperbot/internal/core"
	"scraperbot/internal/domain/model"
	"scraperbot/internal/infrastructure/chromefetcher"
)

func init() {
	core.RegisterFetcher(string(model.FetcherChromium), func(cfg *model.Config) (core.Fetcher, error) {
		return chromefetcher.New(cfg.Request, cfg.Plugins.FetcherConfig)
	})
}
