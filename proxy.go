package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
)

func proxy(c *gin.Context) {
	remoteUrl := c.GetHeader(PROXY_REQ_HEADER)
	remote, err := url.Parse(remoteUrl)
	if err != nil || remote.Scheme == "" || remote.Host == "" {
		c.String(http.StatusBadRequest, "KevinZonda CAS Error: %s", "PARSER_FAILURE")
		return
	}

	c.Request.Header.Del(PROXY_REQ_HEADER)

	proxy := httputil.NewSingleHostReverseProxy(remote)
	cks := c.Request.Cookies()

	proxy.Director = func(req *http.Request) {
		req.Header = c.Request.Header
		req.Host = remote.Host
		req.URL.Scheme = remote.Scheme
		req.URL.Host = remote.Host
		req.URL.Path = remote.Path
		for _, ck := range cks {
			if ck.Name == COOKIE_NAME {
				continue
			}
			req.AddCookie(ck)
		}

	}

	proxy.ServeHTTP(c.Writer, c.Request)
}
