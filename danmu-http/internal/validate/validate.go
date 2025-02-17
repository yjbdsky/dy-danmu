package validate

import (
	"sync"

	"github.com/go-playground/validator/v10"
)

var (
	validate *validator.Validate
	once     sync.Once
)

// GetValidator returns the singleton validator instance
func GetValidator() *validator.Validate {
	once.Do(func() {
		validate = validator.New()

		// 在这里注册自定义验证规则
		// validate.RegisterValidation("customtag", customValidationFunc)
	})
	return validate
}

// Struct validates a struct
func Struct(s interface{}) error {
	return GetValidator().Struct(s)
}

// Var validates a single variable
func Var(field interface{}, tag string) error {
	return GetValidator().Var(field, tag)
}

type PageRequest struct {
	Page           int    `json:"page" form:"page" binding:"required,min=1"`
	PageSize       int    `json:"page_size" form:"page_size" binding:"required,min=1,max=500"`
	OrderBy        string `json:"order_by" form:"order_by" binding:"omitempty"`
	OrderDirection string `json:"order_direction" form:"order_direction" binding:"omitempty"`
}
