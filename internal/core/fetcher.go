package core

import (
	"fmt"

	"scraperbot/internal/domain/model"
	"scraperbot/internal/domain/plugin"
)

// NewFetcherFromConfig は plugins.fetcher 設定に従い登録済み Fetcher を生成する。
func NewFetcherFromConfig(cfg *model.Config) (Fetcher, error) {
	name := string(cfg.Plugins.Fetcher)
	if name == "" {
		name = string(model.FetcherHTTP)
	}
	return Default().NewFetcher(name, cfg)
}

// ResolveHostHTTP はプラグイン Host 向け HTTP クライアントを返す。
// ページ Fetcher が plugin.HTTPClient を実装すればそれを使い、そうでなければ http Fetcher を別途生成する。
func ResolveHostHTTP(cfg *model.Config, pageFetcher Fetcher) (plugin.HTTPClient, error) {
	if hc, ok := pageFetcher.(plugin.HTTPClient); ok {
		return hc, nil
	}
	f, err := Default().NewFetcher(string(model.FetcherHTTP), cfg)
	if err != nil {
		return nil, err
	}
	hc, ok := f.(plugin.HTTPClient)
	if !ok {
		return nil, fmt.Errorf("fetcher %q does not implement plugin.HTTPClient", model.FetcherHTTP)
	}
	return hc, nil
}
