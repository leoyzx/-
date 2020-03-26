package common

import (
	"github.com/dgrijalva/jwt-go"
	"leoyzx/vue+golang/model"
	"time"
)

var jwtKey=[]byte("a_secret_crect")

type Claims struct {
	Userid uint
	jwt.StandardClaims
}

func ReleaseToken(user model.User)(string,error){
	expirationTime:=time.Now().Add(7*24*time.Hour)
	claims:=&Claims{
		Userid:         user.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:time.Now().Unix(),
			Issuer:"admin",
			Subject:"user token",
		},
	}

	token:=jwt.NewWithClaims(jwt.SigningMethodHS256,claims)
	tokenString,err :=token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	return tokenString,nil
}
//eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VyaWQiOjQsImV4cCI6MTU4NTgyNDgxMywiaWF0IjoxNTg1MjIwMDEzLCJpc3MiOiJhZG1pbiIsInN1YiI6InVzZXIgdG9rZW4ifQ.ls1tDdo87Ic3KiM5ah1BV8ZPtMWFSOM81wX8fdbmBEM
func ParseToken(tokenString string) (*jwt.Token,*Claims,error) {
	claims:=&Claims{}
	token,err:=jwt.ParseWithClaims(tokenString,claims,func(token *jwt.Token)(i interface{},err error){
		return jwtKey,nil
	})
	return token,claims,err
}