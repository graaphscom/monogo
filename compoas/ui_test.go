package compoas

import (
	"log"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUIHandler(t *testing.T) {
	handler, _ := UIHandler(SwaggerUIBundleConfig{Url: "/openapi.json"}, "/swagger-ui", log.Fatal)
	handlerFunc := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		handler.ServeHTTP(writer, request)
	})
	require.HTTPBodyContains(t, handlerFunc, "get", "/swagger-ui", nil, "/openapi.json")
	require.HTTPBodyContains(t, handlerFunc, "get", "/swagger-ui/assets/swagger-ui-bundle.js", nil, "For license information please see swagger-ui-bundle.js.LICENSE.txt")
}
