package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
	"leoyzx/vue+golang/common"
	"log"
	"os"
)

func main(){
	InitConfig()
	db:=common.InitDb()
	defer db.Close()

	r:=gin.Default()
	r=CollectRoute(r)

	port:=viper.GetString("port")
	if port!="" {
		panic(r.Run(":"+port))
	}
	panic(r.Run())
}

func InitConfig(){
	workDir,err:=os.Getwd()
	if err != nil {
		panic(err)
	}
	viper.New()
	viper.AddConfigPath(workDir+"/config")
	viper.SetConfigName("application")
	viper.SetConfigType("yml")
	//viper.AddConfigPath(".")
	err = viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}
}



