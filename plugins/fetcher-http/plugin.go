// Package httpplugin は標準 HTTP による URL 取得 Fetcher を登録する。
package httpplugin

import (
	"scraperbot/internal/core"
	"scraperbot/internal/domain/model"
	"scraperbot/internal/infrastructure/httpclient"
)

func init() {
	core.RegisterFetcher(string(model.FetcherHTTP), func(cfg *model.Config) (core.Fetcher, error) {
		return httpclient.New(cfg.Request), nil
	})
}
