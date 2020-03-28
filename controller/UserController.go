package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"leoyzx/vue+golang/common"
	"leoyzx/vue+golang/dto"
	"leoyzx/vue+golang/model"
	"leoyzx/vue+golang/response"
	"leoyzx/vue+golang/util"
	"log"
	"net/http"
)

func Register(ctx *gin.Context) {

	DB := common.GetDB()
	//获取参数

	name := ctx.PostForm("name")
	telephone := ctx.PostForm("telephone")
	password := ctx.PostForm("password")

	//数据验证
	if len(telephone) != 11 {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "手机号长度不符")

		//ctx.JSON(http.StatusUnprocessableEntity,gin.H{"code":422,"msg":"手机号长度不符"})
		//log.Println(len(telephone))
		return
	}
	if len(password) < 6 {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "密码长度不得少于6位")

		//ctx.JSON(http.StatusUnprocessableEntity,gin.H{"code":422,"msg":"密码长度不得少于6位"})
		return
	}
	//验证用户的名字是否为空,则随机生成一个名字
	if len(name) == 0 {
		name = util.RandomName(10)
	}
	//查询手机号是否存在
	if isTelephoneExist(DB, telephone) {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "用户已存在")

		//ctx.JSON(http.StatusUnprocessableEntity,gin.H{"code":422,"msg":"用户已存在"})
		return
	}

	//创建用户
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		response.Response(ctx, http.StatusUnprocessableEntity, 500, nil, "系统加密错误")

		//ctx.JSON(http.StatusInternalServerError,gin.H{"code":500,"msg":"系统加密错误"})
		return
	}
	newUser := model.User{
		Name:      name,
		Telephone: telephone,
		Password:  string(hashedPassword),
	}
	DB.Create(&newUser)

	//返回结果
	//ctx.JSON(200,gin.H{
	//	"code":200,
	//	"msg":"注册成功",
	//})
	response.Success(ctx, nil, "注册成功")

}

func isTelephoneExist(db *gorm.DB, telephone string) bool {
	var user model.User
	db.Where("telephone=?", telephone).First(&user)
	if user.ID != 0 {
		return true
	}
	return false
}

func Login(ctx *gin.Context) {
	DB := common.GetDB()
	//获取参数
	telephone := ctx.PostForm("telephone")
	password := ctx.PostForm("password")
	//数据验证
	if len(telephone) != 11 {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "手机号长度不符"})
		//log.Println(len(telephone))
		return
	}
	if len(password) < 6 {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "密码长度不得少于6位"})
		return
	}
	//判断手机号是否存在
	var user model.User
	DB.Where("telephone=?", telephone).First(&user)
	if user.ID == 0 {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "用户不存在"})
		return
	}
	//判断密码是否正确
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "密码错误"})
		return
	}

	//发放token
	token, err := common.ReleaseToken(user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "系统异常"})
		log.Printf("Token generate error : %v", err)
		return
	}
	//返回结果
	//ctx.JSON(200,gin.H{
	//	"code":200,
	//	"data":gin.H{"token":token},
	//	"msg":"登陆成功",
	//})
	response.Success(ctx, gin.H{"token": token}, "注册成功")
}

func Info(ctx *gin.Context) {
	user, _ := ctx.Get("user")
	ctx.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{"user": dto.ToUserDto(user.(model.User))}}) //类型断言 避免user内数据过多识别类型错误
}
