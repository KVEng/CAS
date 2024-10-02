package main

import (
	"github.com/KVEng/CAS/auth"
	"github.com/KVEng/CAS/shared"
	"github.com/KVEng/CAS/token"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
)

func handleLogin(c *gin.Context) {
	if isLogin(c) {
		return
	}
	session := sessions.Default(c)
	username := c.PostForm("username")
	password := c.PostForm("password")
	if !auth.Verify(username, password, "admin") {
		c.HTML(http.StatusBadRequest, "login.html", gin.H{"error": "Invalid credentials"})
		c.Abort()
	}

	tk := token.TokenGenerator()
	token.ActiveToken(tk, username)

	session.Set("token", tk)
	session.Save()

	if c.Query(shared.REDIRECT_FLAG) != "" {
		c.Redirect(http.StatusFound, c.Query(shared.REDIRECT_FLAG))
		c.Abort()
	}
}

func mustLogin(c *gin.Context) {
	if !isLogin(c) {
		c.Redirect(http.StatusFound, "/cas/login?"+shared.REDIRECT_FLAG+"="+c.Request.URL.String())
		c.Abort()
	}
}

func isLogin(c *gin.Context) bool {
	session := sessions.Default(c)
	tk := session.Get("token")
	if tk == nil {
		return false
	}
	return token.IsTokenValid(tk.(string))
}

func loginPage(c *gin.Context) {
	if isLogin(c) {
		return
	}
	redir := c.Query(shared.REDIRECT_FLAG)
	if redir != "" {
		redir = "?" + shared.REDIRECT_FLAG + "=" + redir
	}
	c.HTML(http.StatusOK, "login.html", gin.H{
		"action": "/cas/login" + redir,
	})
	c.Abort()
}

func logout(c *gin.Context) {
	if !isLogin(c) {
		return
	}
	if c.Query("KEVINZONDA_CAS_IGNORE") == "true" {
		return
	}
	session := sessions.Default(c)
	tk := session.Get("token")
	if tk != nil {
		token.RemoveToken(tk.(string))
	}
	session.Delete("token")
	session.Clear()
	session.Save()
	c.HTML(http.StatusOK, "logout.html", gin.H{})
	c.Abort()
	return

}

func main() {
	shared.InitGlobalCfg()
	shared.InitGlobalRdb()

	store, _ := redis.NewStore(10, "tcp", shared.Config.RedisAddr, "", []byte(shared.REDIS_KEY))

	if shared.Config.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.Default()
	engine.LoadHTMLGlob("html/*")
	engine.Use(sessions.Sessions(shared.COOKIE_NAME, store))

	engine.GET("/cas/logout", logout, proxy)
	engine.GET("/cas/login", loginPage, proxy)
	engine.POST("/cas/login", handleLogin, proxy)
	engine.NoRoute(mustLogin, proxy)

	engine.Run(shared.Config.ListenAddr)
}
