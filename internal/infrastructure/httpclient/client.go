package httpclient

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"time"

	"scraperbot/internal/domain/model"
	"scraperbot/internal/domain/plugin"
)

// Client は domain.plugin.HTTPClient と core が必要とする Get の両方を提供する。
type Client struct {
	// cfg はタイムアウト・リトライ設定。
	cfg model.RequestConfig
	// hc は net/http のクライアント実体。
	hc *http.Client
}

// New は RequestConfig を元に *http.Client をラップした実装を返す。
func New(cfg model.RequestConfig) *Client {
	return &Client{
		cfg: cfg,
		hc: &http.Client{
			Timeout: cfg.Timeout,
		},
	}
}

// Do は plugin.HTTPClient の要件を満たすシンプルな実行 API。
func (c *Client) Do(ctx context.Context, req *plugin.HTTPRequest) (*plugin.HTTPResponse, error) {
	httpReq, err := http.NewRequestWithContext(ctx, req.Method, req.URL.String(), bytes.NewReader(req.Body))
	if err != nil {
		return nil, err
	}
	for k, v := range req.Headers {
		httpReq.Header.Set(k, v)
	}
	res, err := c.hc.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	hdr := flattenHeaders(res.Header)
	return &plugin.HTTPResponse{
		StatusCode: res.StatusCode,
		Headers:    hdr,
		Body:       body,
	}, nil
}

// Get はパイプラインで使う高レベル取得 API。リトライ/タイムアウトを内包する。
func (c *Client) Get(ctx context.Context, u *url.URL, headers map[string]string) (*model.Response, error) {
	var lastErr error
	attempts := c.cfg.RetryCount + 1
	for i := 0; i < attempts; i++ {
		reqCtx, cancel := context.WithTimeout(ctx, c.cfg.Timeout)
		res, err := c.doOnce(reqCtx, u, headers)
		cancel()

		if err == nil && res.StatusCode < 500 {
			return res, nil
		}
		if err != nil {
			lastErr = err
		} else {
			lastErr = fmt.Errorf("http %d", res.StatusCode)
		}
		if err != nil && !isRetryableHTTPError(err) {
			break
		}
		if i+1 < attempts {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(c.cfg.RetryInterval):
			}
		}
	}
	if lastErr == nil {
		lastErr = errors.New("unknown http error")
	}
	return nil, fmt.Errorf("HTTP取得失敗 (url=%s): %w", u.String(), lastErr)
}

// doOnce はリトライなしで 1 回 GET し model.Response を組み立てる。
func (c *Client) doOnce(ctx context.Context, u *url.URL, headers map[string]string) (*model.Response, error) {
	httpReq, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	for k, v := range headers {
		httpReq.Header.Set(k, v)
	}
	res, err := c.hc.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return &model.Response{
		URL:         u,
		StatusCode:  res.StatusCode,
		Headers:     flattenHeaders(res.Header),
		ContentType: res.Header.Get("Content-Type"),
		Body:        body,
		FetchedAt:   time.Now(),
	}, nil
}

// isRetryableHTTPError はタイムアウト・キャンセルなど、再試行しても改善しないエラーを判定する。
func isRetryableHTTPError(err error) bool {
	if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
		return false
	}
	var netErr net.Error
	if errors.As(err, &netErr) && netErr.Timeout() {
		return false
	}
	return true
}

// flattenHeaders は http.Header を先頭値のみの map に変換する。
func flattenHeaders(h http.Header) map[string]string {
	out := make(map[string]string, len(h))
	for k, v := range h {
		if len(v) > 0 {
			out[k] = v[0]
		}
	}
	return out
}

// 静的型チェック: plugin.HTTPClient を満たす。
var _ plugin.HTTPClient = (*Client)(nil)
