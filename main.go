package main

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
)

const PROXY_REQ_HEADER = "KevinZonda-CAS-Proxy"
const COOKIE_NAME = "KEVINZONDA_CAS_SESSION"

func handleLogin(c *gin.Context) {
	if isLogin(c) {
		return
	}
	session := sessions.Default(c)
	username := c.PostForm("username")
	password := c.PostForm("password")
	if username == "admin" && password == "password" {
		session.Set("user", username)
		session.Save()
		// redir := c.Query("redirect")
		// if redir != "" {
		// 	c.Redirect(http.StatusFound, redir)
		// 	return
		// }
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
	return session.Get("username") != nil
}

func loginPage(c *gin.Context) {
	if isLogin(c) {
		return
	}
	c.HTML(http.StatusOK, "login.html", gin.H{})
	c.Abort()
}

func main() {
	store, _ := redis.NewStore(10, "tcp", "localhost:6379", "", []byte("CAS_SESSION"))

	engine := gin.Default()

	engine.LoadHTMLGlob("html/*")

	engine.Use(sessions.Sessions(COOKIE_NAME, store))

	engine.NoRoute(mustLogin)

	engine.GET("/login", loginPage, proxy)
	engine.POST("/login", handleLogin, proxy)
	engine.NoRoute(mustLogin, proxy)

	engine.Run("localhost:11392")
}
