package chromefetcher

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/chromedp/chromedp"

	"scraperbot/internal/core"
	"scraperbot/internal/domain/model"
)

// Client は chromedp による URL フェッチを提供する。
type Client struct {
	// reqCfg はタイムアウト・リトライ設定。
	reqCfg model.RequestConfig
	// fetcherCfg はブラウザ実行・待機に関する設定。
	fetcherCfg model.FetcherConfig
	// browserPath は解決済みのブラウザ実行ファイルパス。
	browserPath string
	// userAgent はリクエストごとに上書き可能な既定 UA（ResolveUserAgent の結果をキャッシュしない）。
}

// New は設定を元に chromedp フェッチャを構築する。
func New(reqCfg model.RequestConfig, fetcherCfg model.FetcherConfig) (*Client, error) {
	path, err := ResolveBrowserPath(fetcherCfg.BrowserPath)
	if err != nil {
		return nil, err
	}
	return &Client{
		reqCfg:      reqCfg,
		fetcherCfg:  fetcherCfg,
		browserPath: path,
	}, nil
}

// Get は指定 URL をヘッドレス Chromium で取得し model.Response を返す。
func (c *Client) Get(ctx context.Context, u *url.URL, headers map[string]string) (*model.Response, error) {
	var lastErr error
	attempts := c.reqCfg.RetryCount + 1
	for i := 0; i < attempts; i++ {
		reqCtx, cancel := context.WithTimeout(ctx, c.reqCfg.Timeout)
		res, err := c.fetchOnce(reqCtx, u, headers)
		cancel()

		if err == nil {
			return res, nil
		}
		lastErr = err
		if !isRetryableFetchError(err) {
			break
		}
		if i+1 < attempts {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(c.reqCfg.RetryInterval):
			}
		}
	}
	if lastErr == nil {
		lastErr = errors.New("unknown chromium fetch error")
	}
	return nil, fmt.Errorf("Chromium取得失敗 (url=%s): %w", u.String(), lastErr)
}

func (c *Client) fetchOnce(ctx context.Context, u *url.URL, headers map[string]string) (*model.Response, error) {
	ua := ResolveUserAgent(c.fetcherCfg, headers)

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.ExecPath(c.browserPath),
		chromedp.UserAgent(ua),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-sandbox", true),
	)
	if c.fetcherCfg.Headless {
		opts = append(opts, chromedp.Flag("headless", true))
	} else {
		opts = append(opts, chromedp.Flag("headless", false))
	}

	allocCtx, allocCancel := chromedp.NewExecAllocator(ctx, opts...)
	defer allocCancel()

	browserCtx, browserCancel := chromedp.NewContext(allocCtx)
	defer browserCancel()

	var html string
	tasks := []chromedp.Action{
		chromedp.Navigate(u.String()),
	}
	if sel := strings.TrimSpace(c.fetcherCfg.WaitVisibleSelector); sel != "" {
		tasks = append(tasks, chromedp.WaitVisible(sel, chromedp.ByQuery))
	}
	tasks = append(tasks, chromedp.OuterHTML("html", &html, chromedp.ByQuery))

	if err := chromedp.Run(browserCtx, tasks...); err != nil {
		return nil, err
	}

	return &model.Response{
		URL:         u,
		StatusCode:  200,
		Headers:     map[string]string{"Content-Type": "text/html; charset=utf-8"},
		ContentType: "text/html; charset=utf-8",
		Body:        []byte(html),
		FetchedAt:   time.Now(),
	}, nil
}

func isRetryableFetchError(err error) bool {
	if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
		return false
	}
	msg := strings.ToLower(err.Error())
	if strings.Contains(msg, "executable") && strings.Contains(msg, "not found") {
		return false
	}
	return true
}

// 静的型チェック: core.Fetcher を満たす。
var _ core.Fetcher = (*Client)(nil)
