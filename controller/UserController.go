package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"leoyzx/vue+golang/common"
	"leoyzx/vue+golang/model"
	"leoyzx/vue+golang/util"
	"log"
	"net/http"
)

func Register(ctx *gin.Context) {

	DB:=common.GetDB()
	//获取参数

	name := ctx.PostForm("name")
	telephone := ctx.PostForm("telephone")
	password := ctx.PostForm("password")

	//数据验证
	if len(telephone)!=11 {
		ctx.JSON(http.StatusUnprocessableEntity,gin.H{"code":422,"msg":"手机号长度不符"})
		//log.Println(len(telephone))
		return
	}
	if len(password)<6{
		ctx.JSON(http.StatusUnprocessableEntity,gin.H{"code":422,"msg":"密码长度不得少于6位"})
		return
	}
	//验证用户的名字是否为空,则随机生成一个名字
	if len(name)==0 {
		name = util.RandomName(10)
	}
	//查询手机号是否存在
	if isTelephoneExist(DB,telephone){
		ctx.JSON(http.StatusUnprocessableEntity,gin.H{"code":422,"msg":"用户已存在"})
		return
	}

	newUser:=model.User{
		Name:  name,
		Telephone: telephone,
		Password: password,
	}
	DB.Create(&newUser)

	log.Println(name,telephone,password)
	//创建用户

	//返回结果
	ctx.JSON(200,gin.H{
		"msg":"注册成功",
	})
}

func isTelephoneExist(db *gorm.DB,telephone string) bool {
	var user model.User
	db.Where("telephone=?",telephone).First(&user)
	if user.ID != 0 {
		return true
	}
	return false
}

