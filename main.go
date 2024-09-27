package main

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
)

const PROXY_PARAM = "proxyPath"
const PROXY_REQ_HEADER = "KevinZonda-CAS-Proxy"
const COOKIE_NAME = "KEVINZONDA_CAS_SESSION"

func login(c *gin.Context) {
	session := sessions.Default(c)
	session.Set("username", "KevinZonda")
	session.Save()
}

func handleLogin(c *gin.Context) {
	session := sessions.Default(c)
	username := c.PostForm("username")
	password := c.PostForm("password")
	if username == "admin" && password == "password" {
		session.Set("user", username)
		session.Save()
		// c.Redirect(http.StatusFound, "/dashboard")
	} else {
		c.HTML(http.StatusBadRequest, "login.html", gin.H{"error": "Invalid credentials"})
	}
}

func mustLogin(c *gin.Context) {
	session := sessions.Default(c)

	if session.Get("username") == nil {
		c.String(http.StatusUnauthorized, "KevinZonda CAS Error: %s", "UNAUTHORIZED")
		c.Abort()
	}
}

func loginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{})
}

func main() {
	store, _ := redis.NewStore(10, "tcp", "localhost:6379", "", []byte("CAS_SESSION"))

	engine := gin.Default()

	engine.Use(sessions.Sessions(COOKIE_NAME, store))

	engine.GET("/login", loginPage)
	engine.POST("/login", handleLogin)

	engine.Any("/px/*"+PROXY_PARAM, mustLogin, proxy)

	engine.Run("localhost:11392")
}
