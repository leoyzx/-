package main

import (
	"github.com/gin-gonic/gin"
	"leoyzx/vue+golang/controller"
	"leoyzx/vue+golang/middleware"
)

func CollectRoute(r *gin.Engine) *gin.Engine {
	r.POST("/api/auth/register", controller.Register)
	r.POST("/api/auth/login",controller.Login)
	r.POST("api/auth/info",middleware.AuthMiddleware(),controller.Info)
	return r
}

