package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"log"
	"math/rand"
	"net/http"
	"time"
)

type User struct {
	gorm.Model
	Name string  `gorm:"type:varchar(20);not null"`
	Telephone string `gorm:"type:varchar(11);not null;unique"`
	Password string `gorm:"size:255;not null"`

}

func main(){
	db:=InitDb()
	defer db.Close()

	r:=gin.Default()
	r.POST("/api/auth/register", func(ctx *gin.Context) {
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
			name = RandomName(10)
		}
		//查询手机号是否存在
		if isTelephoneExist(db,telephone){
			ctx.JSON(http.StatusUnprocessableEntity,gin.H{"code":422,"msg":"用户已存在"})
			return
		}

		newUser:=User{
			Name:  name,
			Telephone:telephone,
			Password:password,
		}
		db.Create(&newUser)

		log.Println(name,telephone,password)
		//创建用户

		//返回结果
		ctx.JSON(200,gin.H{
			"msg":"注册成功",
		})
	})
	panic(r.Run())
}

func RandomName(n int) string{
	var letters  = []byte("abcdefghijklmnopqrstuvwxyz1234567890")
	result:=make([]byte,n)
	rand.Seed(time.Now().Unix())
	for i:=range result {
		result[i] = letters[rand.Intn(len(letters))]
	}
	return  string(result)
}

func InitDb()*gorm.DB{
	driverName := "mysql"
	host :="localhost"
	port:= "3306"
	database:= "ginessential"
	username:="root"
	password:= "root"
	charset:="utf8"
	args:= fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true",
		username,password,host,port,database,charset)
	db,err := gorm.Open(driverName,args)
	if err != nil {
		panic("failed to connnect database,err:"+err.Error())
	}
	db.AutoMigrate(&User{})
	return db
}

func isTelephoneExist(db *gorm.DB,telephone string) bool {
	var user User
	db.Where("telephone=?",telephone).First(&user)
	if user.ID != 0 {
		return true
	}
	return false
}