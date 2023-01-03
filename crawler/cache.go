package crawler

import (
	"context"
	"net/http"
	"time"
)

func CacheAllPages(ctx context.Context, options CachePagesOptions, delay time.Duration) []CachingResult {
	cachingResults := CachePages(ctx, options)

	for {
		anySkipped := false
		for _, cachingResult := range cachingResults {
			if cachingResult.Skipped {
				anySkipped = true
				break
			}
		}

		if anySkipped {
			time.Sleep(delay)
			cachingResults = CachePages(ctx, options)
		} else {
			break
		}
	}

	return cachingResults
}

func CachePages(ctx context.Context, options CachePagesOptions) []CachingResult {
	pagesCount := len(options.Pages)
	resultCh := make(chan CachingResult, pagesCount)
	result := make([]CachingResult, pagesCount)
	requestsCount := 0

	for url, cacheKey := range options.Pages {
		exists, err := options.CacheStorage.exists(cacheKey)

		if err != nil {
			resultCh <- CachingResult{Url: url, Key: cacheKey, Err: err}
			requestsCount++
			continue
		}

		if exists {
			resultCh <- CachingResult{Url: url, Key: cacheKey, Hit: true}
			continue
		}

		if requestsCount >= options.RequestsLimit {
			resultCh <- CachingResult{Url: url, Key: cacheKey, Skipped: true}
			continue
		}

		requestsCount++
		go func(url string, cacheKey string) {
			resultCh <- WritePage(ctx, options.HttpClient, url, cacheKey, options.CacheStorage)
		}(url, cacheKey)
	}

	for i := 0; i < pagesCount; i++ {
		result[i] = <-resultCh
	}
	return result
}

type CachePagesOptions struct {
	HttpClient    *http.Client
	Pages         PageUrlToCacheKey
	RequestsLimit int
	CacheStorage  CacheStorage
}

func WritePage(ctx context.Context, httpClient *http.Client, url string, key string, cacheStorage CacheStorage) CachingResult {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return CachingResult{Url: url, Key: key, Err: err}
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return CachingResult{Url: url, Key: key, Err: err}
	}
	defer resp.Body.Close()

	err = cacheStorage.write(key, resp.Body)
	if err != nil {
		return CachingResult{Url: url, Key: key, Err: err}
	}

	return CachingResult{Url: url, Key: key}
}

type CachingResult struct {
	Url     string
	Key     string
	Hit     bool
	Skipped bool
	Err     error
}

type PageUrlToCacheKey = map[string]string

func BuildCachingSummary(cachingResult []CachingResult) CachingSummary {
	summary := CachingSummary{}
	for _, pageResult := range cachingResult {
		if pageResult.Err != nil {
			summary.Err = append(summary.Err, pageResult)
			continue
		}
		if pageResult.Hit {
			summary.Hit = append(summary.Hit, pageResult)
			continue
		}
		if pageResult.Skipped {
			summary.Skipped = append(summary.Skipped, pageResult)
			continue
		}

		summary.Written = append(summary.Written, pageResult)

	}
	return summary
}

type CachingSummary struct {
	Written []CachingResult
	Hit     []CachingResult
	Skipped []CachingResult
	Err     []CachingResult
}
