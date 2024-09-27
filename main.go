package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
)

const PROXY_PARAM = "/proxyPath"
const PROXY_REQ_HEADER = "KevinZonda-CAS-Proxy"

func proxy(c *gin.Context) {
	remoteUrl := c.GetHeader(PROXY_REQ_HEADER)
	remote, err := url.Parse(remoteUrl)
	if err != nil || remote.Scheme == "" || remote.Host == "" {
		c.String(http.StatusBadRequest, "KevinZonda CAS Error: %s", "PARSER_FAILURE")
		return
	}

	c.String(http.StatusOK, "KevinZonda CAS Proxy: %s", remoteUrl)
	return

	c.Request.Header.Del(PROXY_REQ_HEADER)

	proxy := httputil.NewSingleHostReverseProxy(remote)

	proxy.Director = func(req *http.Request) {
		req.Header = c.Request.Header
		req.Host = remote.Host
		req.URL.Scheme = remote.Scheme
		req.URL.Host = remote.Host
		req.URL.Path = remote.Path

	}

	proxy.ServeHTTP(c.Writer, c.Request)
}

func main() {
	engine := gin.Default()

	engine.Any("/*"+PROXY_PARAM, proxy)

	engine.Run("localhost:11392")
}
