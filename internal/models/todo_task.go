package models

import "gorm.io/gorm"

type TodoTask struct {
    gorm.Model  // adds ID, created_at etc.
    Title       string `json:"title"`
    Description string `json:"description"`
	State       bool   `json:"state"`
}
