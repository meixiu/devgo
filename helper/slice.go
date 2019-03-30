package helper

import "reflect"

type (
	GetIder interface {
		GetId() int
	}
)

// SliceIdIndex 获取[]GetIder切片里所有ID => index的对应关系
func SliceIdIndex(slice interface{}) map[int]int {
	value := reflect.ValueOf(slice)
	if value.Kind() != reflect.Slice {
		panic("parameters can only be slices")
	}
	m := make(map[int]int)
	for i := 0; i < value.Len(); i++ {
		v, ok := value.Index(i).Interface().(GetIder)
		if !ok {
			panic("The slice element must implement the GetIder interface")
		}
		m[v.GetId()] = i
	}
	return m
}

// SliceIds 获取[]GetIder切片里所有ID的集合
func SliceIds(slice interface{}) []int {
	value := reflect.ValueOf(slice)
	if value.Kind() != reflect.Slice {
		panic("parameters can only be slices")
	}
	m := make([]int, 0)
	for i := 0; i < value.Len(); i++ {
		v, ok := value.Index(i).Interface().(GetIder)
		if !ok {
			panic("The slice element must implement the GetIder interface")
		}
		m = append(m, v.GetId())
	}
	return m
}
