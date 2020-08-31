package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go_util/api/v1/login"
	"time"
)

var router *gin.Engine

func init() {
	InitRunParam()
	router = gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOriginFunc:  func(origin string) bool { return true },
		AllowMethods:     []string{"OPTIONS", "GET", "POST", "PUT", "DELETE", "PATCH"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// 模块路由地址
	loginRoute("login")

	router.Run(Run.Port)
}

// 路由地址
func loginRoute(name string) {
	route := router.Group(name)
	route.Any("login", login.UserLogin)
}
