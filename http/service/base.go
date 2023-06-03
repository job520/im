package service

import (
	"fmt"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhs "github.com/go-playground/validator/v10/translations/zh"
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
