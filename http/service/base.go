package service

import (
	"fmt"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhs "github.com/go-playground/validator/v10/translations/zh"
	"github.com/golang-jwt/jwt"
	"time"
)

func Validate(param interface{}) error {
	validate := validator.New()
	chinese := zh.New()
	uni := ut.New(chinese, chinese)
	trans, _ := uni.GetTranslator("zh")
	_ = zhs.RegisterDefaultTranslations(validate, trans)
	if err := validate.Struct(param); err != nil {
		if errs, ok := err.(validator.ValidationErrors); ok {
			errStr := ""
			for k, v := range errs.Translate(trans) {
				errStr += fmt.Sprintf("%s: %s；", k, v)
			}
			return fmt.Errorf("%s", errStr)
		}
	}
	return nil
}

// 生成 jwt-token
func GenerateJwtToken(encryptKey string, expireHours int, uid string) (string, error) {
	expire := time.Duration(expireHours)
	exp := time.Now().Add(time.Hour * expire).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uid": uid,
		"exp": exp,
	})
	tokenString, err := token.SignedString([]byte(encryptKey))
	return tokenString, err
}

// 解析 jwt-token
func ParseJwtToken(encryptKey, tokenString string) (string, bool) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(encryptKey), nil
	})
	if err != nil {
		return "", false
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if uid, ok := claims["uid"].(string); ok {
			return uid, true
		} else {
			return "", false
		}
	} else {
		return "", false
	}
}
