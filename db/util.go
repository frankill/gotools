package db

import "reflect"

// GetTags 获取结构体的 json 标签
//
// obj - 结构体
//
// 返回 json 标签列表
func GetTags(obj interface{}) []string {
	var tags []string
	typ := reflect.TypeOf(obj)

	if typ.Kind() != reflect.Struct {
		return tags
	}

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		tag := field.Tag.Get("json")
		if tag != "" {
			tags = append(tags, tag)
		}
	}

	return tags
}
