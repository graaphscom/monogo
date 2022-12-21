package crawler

import (
	"context"
	"errors"
	"io"
	"io/fs"
	"log"
	"net/http"
	"net/http/cookiejar"
	"os"
	"path"
	"sync"
)

func CachePages(ctx context.Context, httpClient *http.Client, pages map[string]FilePath) <-chan error {
	var wg sync.WaitGroup
	wg.Add(len(pages))
	result := make(chan error)

	for url, dstFilePath := range pages {
		go func(url string, dstFilePath FilePath) {
			defer wg.Done()

			err := os.MkdirAll(dstFilePath.DirPath, 0750)
			if err != nil {
				result <- err
				return
			}

			err = CachePage(
				ctx,
				url,
				path.Join(dstFilePath.DirPath, dstFilePath.FileName),
				httpClient,
			)

			if err != nil {
				result <- err
				return
			}
		}(url, dstFilePath)
	}
	wg.Wait()

	return result
}

type FilePath struct {
	DirPath  string
	FileName string
}

func CachePage(ctx context.Context, url string, filePath string, httpClient *http.Client) error {
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0666)
	defer func() { err = file.Close() }()
	if errors.Is(err, fs.ErrExist) {
		// don't make http request when html page already saved
		return nil
	}
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	resp, err := httpClient.Do(req)
	_, err = io.Copy(file, resp.Body)
	defer resp.Body.Close()

	return err
}

func NewHttpClient(logger *log.Logger) (*http.Client, error) {
	httpClientCookieJar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}
	return &http.Client{
		// Without proper cookie handling it will fall into a redirect loop
		Jar:       httpClientCookieJar,
		Transport: loggingRoundTripper{proxied: http.DefaultTransport, logger: logger},
	}, nil
}

type loggingRoundTripper struct {
	proxied http.RoundTripper
	logger  *log.Logger
}

func (lrt loggingRoundTripper) RoundTrip(req *http.Request) (res *http.Response, e error) {
	lrt.logger.Printf("Sending request to %s\n", req.URL)
	return lrt.proxied.RoundTrip(req)
}
