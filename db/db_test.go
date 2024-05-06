package db

import (
	"reflect"
	"testing"

	"gorm.io/gorm"
)

func TestInit(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name string
		args args
		want *gorm.DB
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Init(tt.args.url); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Init() = %v, want %v", got, tt.want)
			}
		})
	}
}
