package model

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type TodoTask struct {
	gorm.Model         // adds ID, created_at etc.
	Title       string `json:"title" validate:"required"`
	Description string `json:"description" validate:"required"`
	State       bool   `json:"state" validate:"required"`
}

type TodoTaskPayload struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description" validate:"required"`
	State       bool   `json:"state" validate:"required"`
}

// ValidateTodoTaskPayload validates the TodoTaskPayload fields
func (t *TodoTaskPayload) ValidateTodoTaskPayload() error {
	v := validator.New()
	err := v.Struct(t)
	if err != nil {
		return fmt.Errorf("validation fails: %v", err)
	}
	return nil
}

// Validator instance
var validate *validator.Validate

func init() {
	validate = validator.New()
}

// ValidateTodoTask validates the TodoTask struct
func ValidateTodoTask(todoTask TodoTask) error {
	err := validate.Struct(todoTask)
	if err != nil {
		return fmt.Errorf("validation fails: %v", err)
	}
	return nil
}
