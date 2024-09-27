package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
)

const PROXY_PARAM = "proxyPath"
const PROXY_REQ_HEADER = "KevinZonda-CAS-Proxy"

func proxy(c *gin.Context) {
	session := sessions.Default(c)

	if session.Get("username") == nil {
		c.String(http.StatusUnauthorized, "KevinZonda CAS Error: %s", "UNAUTHORIZED")
		return
	}

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
	store, _ := redis.NewStore(10, "tcp", "localhost:6379", "", []byte("CAS_SESSION"))

	engine := gin.Default()

	engine.Use(sessions.Sessions("KEVINZONDA_CAS_SESSION", store))

	engine.GET("/login", func(c *gin.Context) {
		session := sessions.Default(c)
		session.Set("username", "KevinZonda")
		session.Save()
	})

	engine.Any("/c/*"+PROXY_PARAM, proxy)

	engine.Run("localhost:11392")
}
