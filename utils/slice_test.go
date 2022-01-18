package utils

import (
	"reflect"
	"testing"
)

func TestToSlice(t *testing.T) {
	arr := []int{1, 2, 3}

	s := ToSlice(arr)
	if reflect.ValueOf(s).Type().Kind() != reflect.Slice {
		t.Errorf("should return a slice, but got %s", reflect.ValueOf(s).Type().Kind())
	}
}
