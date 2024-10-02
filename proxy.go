package main

import (
	"github.com/KVEng/CAS/shared"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
)

func proxy(c *gin.Context) {
	remoteUrl := c.GetHeader(shared.PROXY_REQ_HEADER)
	remote, err := url.Parse(remoteUrl)
	if err != nil || remote.Scheme == "" || remote.Host == "" {
		c.String(http.StatusBadRequest, "KevinZonda CAS Error: %s", "PARSER_FAILURE")
		return
	}

	c.Request.Header.Del(shared.PROXY_REQ_HEADER)

	px := httputil.NewSingleHostReverseProxy(remote)
	cks := c.Request.Cookies()

	px.Director = func(req *http.Request) {
		req.Header = c.Request.Header
		req.Host = remote.Host
		req.URL.Scheme = remote.Scheme
		req.URL.Host = remote.Host
		req.URL.Path = remote.Path
		for _, ck := range cks {
			if ck.Name == shared.COOKIE_NAME {
				continue
			}
			req.AddCookie(ck)
		}

	}

	px.ServeHTTP(c.Writer, c.Request)
}
