package model

import (
	"github.com/go-playground/validator/v10"
	"reflect"
	"strings"
)

// Кастомная валидация OrderInfo
func validateNoSpecialChars(sl validator.StructLevel) {
	order := sl.Current().Interface().(OrderInfo)
	orderType := reflect.TypeOf(order)
	orderValue := reflect.ValueOf(order)

	//Запрет спецсимволов на пример точки с запятой и двойного дефиса
	for i := 0; i < orderValue.NumField(); i++ {
		fValue := orderValue.Field(i)
		if fValue.Kind() == reflect.String {
			fName := orderType.Field(i).Name
			fStr := fValue.String()

			for _, char := range []string{";", "--"} {
				if strings.Contains(fStr, char) {
					sl.ReportError(fValue, fName, "noSpecialChars", "", "")
				}

			}
		}
	}
}

// Регистрация кастомного валидатора для структуры OrderInfo
func init() {
	validate.RegisterStructValidation(validateNoSpecialChars, OrderInfo{})
}
