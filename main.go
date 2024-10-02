package main

import (
	"fmt"
	"github.com/KVEng/CAS/auth"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
)

const PROXY_REQ_HEADER = "KevinZonda-CAS-Proxy"
const COOKIE_NAME = "KEVINZONDA_CAS_SESSION"
const REDIS_KEY = "CAS_SESSION"
const REDIS_ADDR = "localhost:6379"

func handleLogin(c *gin.Context) {
	if isLogin(c) {
		return
	}
	session := sessions.Default(c)
	username := c.PostForm("username")
	password := c.PostForm("password")
	if auth.Verify(username, password, "admin") {
		session.Set("username", username)
		session.Save()
	} else {
		c.HTML(http.StatusBadRequest, "login.html", gin.H{"error": "Invalid credentials"})
	}
	c.Abort()
}

func mustLogin(c *gin.Context) {
	if !isLogin(c) {
		c.Redirect(http.StatusFound, "/login?redirect="+c.Request.URL.String())
		c.Abort()
	}
}

func isLogin(c *gin.Context) bool {
	session := sessions.Default(c)
	fmt.Println(session.Get("username"))
	fmt.Println(session)
	return session.Get("username") != nil
}

func loginPage(c *gin.Context) {
	if !isLogin(c) {
		c.HTML(http.StatusOK, "login.html", gin.H{})
		c.Abort()
	}
}

func main() {
	store, _ := redis.NewStore(10, "tcp", REDIS_ADDR, "", []byte(REDIS_KEY))

	engine := gin.Default()

	engine.LoadHTMLGlob("html/*")

	engine.Use(sessions.Sessions(COOKIE_NAME, store))

	engine.NoRoute(mustLogin)

	engine.GET("/cas/logout", func(c *gin.Context) {
		if isLogin(c) {
			session := sessions.Default(c)
			session.Clear()
			c.String(http.StatusOK, "CAS Logged out")
			c.Abort()
			return
		}
	}, proxy)
	engine.GET("/login", loginPage, proxy)
	engine.POST("/login", handleLogin, proxy)
	engine.GET("/incr", func(c *gin.Context) {
		session := sessions.Default(c)
		var count int
		v := session.Get("count")
		if v == nil {
			count = 0
		} else {
			count = v.(int)
			count++
		}
		session.Set("count", count)
		session.Save()
		c.JSON(200, gin.H{"count": count})
	})
	engine.NoRoute(mustLogin, proxy)

	engine.Run("localhost:11392")
}
