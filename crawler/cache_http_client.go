package crawler

import (
	"log"
	"net/http"
	"net/http/cookiejar"
)

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
