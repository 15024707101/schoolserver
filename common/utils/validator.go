package utils

import (
	"fmt"
	"github.com/go-playground/locales/zh"
	"github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
	"reflect"
)

//func ParamsValidator (s interface{}) error{
//	validate := validator.New()
//	errs := validate.Struct(s)
//	return errs
//}


func ParamsValidator (s interface{} ) string {
	uni := ut.New(zh.New())
	trans, _ := uni.GetTranslator("zh")
	validate := validator.New()

	//注册一个函数，获取struct tag里自定义的label作为字段名
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name:=fld.Tag.Get("label")
		return name
	})
	//注册翻译器
	err := zh_translations.RegisterDefaultTranslations(validate, trans)
	if err!=nil {
		fmt.Println(err)
	}

	err = validate.Struct(s)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			fmt.Println(err.Translate(trans))
			return err.Translate(trans)
		}
	}

	return ""

}