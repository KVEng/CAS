package main

import (
	"github.com/KVEng/CAS/auth"
	"github.com/KVEng/CAS/model"
	"github.com/KVEng/CAS/shared"
	"github.com/KVEng/CAS/token"
	"github.com/KevinZonda/GoX/pkg/panicx"
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
	if !auth.Verify(username, password, "") {
		c.HTML(http.StatusBadRequest, "login.html", gin.H{"error": "Invalid credentials"})
		c.Abort()
		return
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
	if isLogin(c) {
		return
	}
	c.Redirect(http.StatusFound, "/cas/login?"+shared.REDIRECT_FLAG+"="+c.Request.URL.String())
	c.Abort()
}

func isLogin(c *gin.Context) bool {
	session := sessions.Default(c)
	tk := session.Get("token")
	if tk == nil {
		return false
	}
	return token.IsTokenValid(tk.(string))
}

func verifyGroupByToken(c *gin.Context) bool {
	session := sessions.Default(c)
	tkStr := session.Get("token")
	if tkStr == nil {
		return false
	}
	tk := tkStr.(string)
	username := token.GetTokenUsername(tk)
	if username == "" {
		return false
	}
	return verifyGroup(c, username)
}

func verifyGroup(c *gin.Context, username string) bool {
	group := c.GetHeader(shared.GROUP_HEADER)
	if group == "" {
		group = "admin"
	}
	u, ok := shared.UserDb[username]
	if !ok {
		return false
	}
	return u.IsInGroup(group)
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
	if c.Query("KEVINZONDA_CAS_IGNORE") == "true" {
		return
	}
	if !isLogin(c) {
		c.Redirect(http.StatusFound, "/cas/login")
		c.Abort()
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

func changePasswordPage(c *gin.Context) {
	if c.Query("KEVINZONDA_CAS_IGNORE") == "true" {
		return
	}
	c.HTML(http.StatusOK, "change-password.html", gin.H{})
	c.Abort()
}

func handleChangePassword(c *gin.Context) {
	if c.Query("KEVINZONDA_CAS_IGNORE") == "true" {
		return
	}
	username := c.PostForm("username")
	oldPassword := c.PostForm("old_password")
	newPassword := c.PostForm("new_password")
	confirmNewPassword := c.PostForm("confirm_new_password")

	if newPassword != confirmNewPassword {
		c.HTML(http.StatusBadRequest, "change-password.html", gin.H{"error": "Password not match"})
		c.Abort()
		return
	}

	if !auth.Verify(username, oldPassword, "") {
		c.HTML(http.StatusBadRequest, "change-password.html", gin.H{"error": "Invalid credentials"})
		c.Abort()
		return
	}

	err := shared.ModifyUserDb(func(db map[string]model.User) {
		u, ok := db[username]
		if ok {
			u.Password = token.HashPasswd(newPassword)
			db[username] = u
		}
	})

	if err != nil {
		c.HTML(http.StatusBadRequest, "change-password.html", gin.H{"error": "Storage unit failure"})
		c.Abort()
		return
	}
	c.HTML(http.StatusOK, "change-password.html", gin.H{"error": "Password changed successfully"})
	c.Abort()
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

	engine.GET("/cas/password", changePasswordPage, mustLogin, proxy)
	engine.POST("/cas/password", handleChangePassword, mustLogin, proxy)

	engine.NoRoute(mustLogin, proxy)

	err := engine.Run(shared.Config.ListenAddr)
	panicx.NotNilErr(err)
}
