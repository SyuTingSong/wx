package main

import (
	"github.com/gin-gonic/gin"
	"wx/login"
)

func main() {
	r := gin.Default()
	rLogin := r.Group("/login")
	{
		rLogin.GET("/l2", login.L2)
		rLogin.GET("/user_info", login.UserInfo)
	}
}
