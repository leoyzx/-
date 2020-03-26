package middleware

import (
	"github.com/gin-gonic/gin"
	"leoyzx/vue+golang/common"
	"leoyzx/vue+golang/model"
	"net/http"
	"strings"
)

func AuthMiddleware()gin.HandlerFunc{
	return func(context *gin.Context) {
		//获取authorization header
		tokenString := context.GetHeader("Authorization")
		//validate token formate 验证token格式
		if tokenString==""||!strings.HasPrefix(tokenString,"Bearer "){
			context.JSON(http.StatusUnauthorized,gin.H{"code":401,"msg":"权限不足"})
			context.Abort()
			return
		}
		tokenString = tokenString[7:]
		token ,claims,err :=common.ParseToken(tokenString)
		if err != nil||!token.Valid {
			context.JSON(http.StatusUnauthorized,gin.H{"code":401,"msg":"权限不足"})
			context.Abort()
			return
		}

		userId:=claims.Userid
		DB:=common.GetDB()
		var user model.User
		DB.First(&user,userId)
		if userId==0{
			context.JSON(http.StatusUnauthorized,gin.H{"code":401,"msg":"权限不足"})
			context.Abort()
			return
		}

		context.Set("user",user)
		context.Next()
	}
}