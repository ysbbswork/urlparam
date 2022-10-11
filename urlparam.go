// Package urlparam 提供url参数自动解析到结构体，实现url.Value to Struct和 Struct to url.Value
// 支持结构体标签，利用结构体的反射标签可以自定义参数名称的映射关系，
// 默认利用json标签，可以使此库可直接用pb.go生成的结构体，便于不同服务之间使用pb管理url params 的协议
// 支持修改标签类型，解析标签名可修改导出变量URLParamTag来改变
// 没有定义标签则默认使用结构体的字段名称作为url参数名称，并且 - 标签作为不解析标记，和json标签语法兼容
package urlparam

import (
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"strings"

	"github.com/spf13/cast"
)

// URLParamTag url 参数解编码时的标签名
var URLParamTag = "json"

// Marshal 将结构体编码成url参数的形式
func Marshal(s interface{}) (string, error) {
	v, err := Encode(s)
	if err != nil {
		return "", err
	}
	return v.Encode(), nil
}

// Unmarshal 将url参数解析到结构体中
func Unmarshal(uri string, s interface{}) error {
	values, err := url.ParseQuery(uri)
	if err != nil {
		return err
	}
	return Decode(values, s)
}

// Encode encode struct to url.Values
func Encode(s interface{}) (url.Values, error) {
	sv := reflectValue(s)
	st := reflectType(s)
	if err := checkType(st); err != nil {
		return nil, err
	}
	values := url.Values{}
	for i := 0; i < st.NumField(); i++ {
		valueField := sv.Field(i)
		if !valueField.CanInterface() { // 是否可导出
			continue
		}
		field := st.Field(i)
		structFieldValue := sv.FieldByName(field.Name)
		if structFieldValue.IsValid() {
			key := getValueKey(field)
			if key == "" {
				continue
			}
			values.Set(key, cast.ToString(structFieldValue.Interface()))
		}
	}
	return values, nil
}

// Decode decode url.Values to struct
func Decode(urlValues url.Values, s interface{}) error {
	sv := reflectValue(s)
	st := reflectType(s)
	if err := checkType(st); err != nil {
		return err
	}
	for i := 0; i < sv.NumField(); i++ {
		valueField := sv.Field(i)
		if !valueField.CanInterface() { // 是否可导出
			continue
		}
		structFeild := st.Field(i)
		key := getValueKey(structFeild)
		if key == "" {
			continue
		}
		v := urlValues.Get(key)
		switch valueField.Kind() {
		case reflect.String:
			valueField.SetString(v)
		case reflect.Bool:
			valueField.SetBool(cast.ToBool(v))
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			valueField.SetUint(cast.ToUint64(v))
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			valueField.SetInt(cast.ToInt64(v))
		case reflect.Float32, reflect.Float64:
			valueField.SetFloat(cast.ToFloat64(v))
		default:
			return fmt.Errorf("unsupported type: %v ,val: %v ,query key: %v", valueField.Type(), v, key)
		}
	}
	return nil
}

func getValueKey(field reflect.StructField) string {
	tag := field.Tag.Get(URLParamTag)
	if tag == "-" {
		return ""
	}
	key, _ := parseTag(tag)
	if key == "" {
		key = field.Name
	}
	return key
}

func checkType(typ reflect.Type) error {
	if typ.Kind() != reflect.Struct {
		return errors.New("input must be a struct or struct pointer")
	}
	return nil
}

func reflectValue(s interface{}) (val reflect.Value) {
	if reflect.TypeOf(s).Kind() != reflect.Struct {
		val = reflect.ValueOf(s).Elem()
	} else if reflect.TypeOf(s).Kind() != reflect.Ptr {
		val = reflect.ValueOf(s)
	}
	return
}

func reflectType(s interface{}) (typ reflect.Type) {
	return reflectValue(s).Type()
}

func parseTag(tag string) (tagName string, tagOptions []string) {
	s := strings.Split(tag, ",")
	return s[0], s[1:]
}