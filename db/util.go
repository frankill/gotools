package db

import "reflect"

// GetTags 获取结构体的 json 标签
//
// obj - 结构体
//
// 返回 json 标签列表
func GetJsonTags(obj any) []string {
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

// GetTags 获取结构体的 sql 标签
//
// obj - 结构体
//
// 返回 sql 标签列表
func GetSqlTags(obj any) []string {
	var tags []string
	typ := reflect.TypeOf(obj)

	if typ.Kind() != reflect.Struct {
		return tags
	}

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		tag := field.Tag.Get("sql")
		if tag != "" {
			tags = append(tags, tag)
		}
	}

	return tags
}

// GetTags 获取结构体的 tag 标签
//
// obj - 结构体
// tag - tag 名称
//
// 返回 tag 标签列表
func GetTags(obj any, tag string) []string {
	var tags []string
	typ := reflect.TypeOf(obj)

	if typ.Kind() != reflect.Struct {
		return tags
	}

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		tag := field.Tag.Get(tag)
		if tag != "" {
			tags = append(tags, tag)
		}
	}

	return tags
}

// GetFieldValues 获取结构体的字段
//
// obj - 结构体
//
// 返回  []T
func GetFieldValues[T any](obj any) []T {

	val := reflect.ValueOf(obj)
	typ := reflect.TypeOf(obj)

	values := make([]T, val.NumField())

	// 确保输入的 obj 是一个结构体
	if typ.Kind() != reflect.Struct {
		return values
	}

	for i := 0; i < val.NumField(); i++ {
		fieldValue := val.Field(i)

		// 确保字段是可导出的
		if !fieldValue.CanInterface() {
			continue
		}

		// 转换字段值为指定类型 T
		if convertedValue, ok := fieldValue.Interface().(T); ok {
			values[i] = convertedValue
		}
	}

	return values
}

// GetFieldMapValues 获取结构体的字段
//
// obj - 结构体
//
// 返回 key-value 字典
func GetFieldMapValues[T any](obj any) map[string]T {

	val := reflect.ValueOf(obj)
	typ := reflect.TypeOf(obj)

	values := make(map[string]T, val.NumField())

	// 确保输入的 obj 是一个结构体
	if typ.Kind() != reflect.Struct {
		return values
	}

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		fieldValue := val.Field(i)

		// 确保字段是可导出的
		if !fieldValue.CanInterface() {
			continue
		}

		// 获取字段名称
		fieldName := field.Name

		// 转换字段值为指定类型 T
		if convertedValue, ok := fieldValue.Interface().(T); ok {
			values[fieldName] = convertedValue
		}
	}

	return values
}
