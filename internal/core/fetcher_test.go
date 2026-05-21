package core_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"scraperbot/internal/core"
	"scraperbot/internal/domain/model"
	"scraperbot/internal/infrastructure/httpclient"

	_ "scraperbot/plugins/fetcher-chromium"
	_ "scraperbot/plugins/fetcher-http"
)

func TestNewFetcherFromConfig_httpDefault(t *testing.T) {
	t.Parallel()
	cfg := model.Default()
	cfg.Targets = []string{"https://example.com/"}

	f, err := core.NewFetcherFromConfig(&cfg)
	require.NoError(t, err)
	_, ok := f.(*httpclient.Client)
	assert.True(t, ok, "デフォルト fetcher は httpclient")
}

func TestNewFetcherFromConfig_chromiumMissingBrowser(t *testing.T) {
	t.Parallel()
	cfg := model.Default()
	cfg.Targets = []string{"https://example.com/"}
	cfg.Plugins.Fetcher = model.FetcherChromium
	cfg.Plugins.FetcherConfig.BrowserPath = "/nonexistent/chromium-binary"

	_, err := core.NewFetcherFromConfig(&cfg)
	require.Error(t, err)
}

func TestNewFetcherFromConfig_unknownFetcher(t *testing.T) {
	t.Parallel()
	cfg := model.Default()
	cfg.Targets = []string{"https://example.com/"}
	cfg.Plugins.Fetcher = "selenium"

	_, err := core.NewFetcherFromConfig(&cfg)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "fetcher not found")
}

func TestResolveHostHTTP_httpUsesSameFetcher(t *testing.T) {
	t.Parallel()
	cfg := model.Default()
	cfg.Targets = []string{"https://example.com/"}

	pageFetcher, err := core.NewFetcherFromConfig(&cfg)
	require.NoError(t, err)
	hc, err := core.ResolveHostHTTP(&cfg, pageFetcher)
	require.NoError(t, err)
	assert.Same(t, pageFetcher, hc)
}
