package configParser

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
)

func SetValue(ptr interface{}, tag string) error {
	if reflect.TypeOf(ptr).Kind() != reflect.Ptr {
		return fmt.Errorf("not a pointer")
	}
	v := reflect.ValueOf(ptr).Elem()
	return setValueRecursion(v, tag)
}

func setValueRecursion(v reflect.Value, tag string) error {
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		//fmt.Printf("%d. %v (%v), tag: '%v'\n", i+1, f.Name, f.Type.Name(), f.Tag.Get(tag))
		if _type := f.Type.Kind(); _type != reflect.Struct {
			value := func() (string, bool) {
				switch tag {
				case "default":
					val := t.Field(i).Tag.Get(tag)
					return val, val != "-"
				case "env":
					val := os.Getenv(t.Field(i).Tag.Get("env"))
					return val, val != ""
				}
				panic(fmt.Sprintf("unknown type %s", tag))
			}
			if val, do := value(); do == true {
				if err := setField(v.Field(i), val); err != nil {
					continue
				}
			}
		} else {
			setValueRecursion(v.Field(i), tag)
		}
	}
	return nil
}

func setField(field reflect.Value, defaultVal string) error {

	if !field.CanSet() {
		return fmt.Errorf("Can't set value\n")
	}

	switch field.Kind() {

	case reflect.Int:
		if val, err := strconv.ParseInt(defaultVal, 10, 64); err == nil {
			field.Set(reflect.ValueOf(int(val)).Convert(field.Type()))
		}
	case reflect.String:
		field.Set(reflect.ValueOf(defaultVal).Convert(field.Type()))
	case reflect.Bool:
		if val, err := strconv.ParseBool(defaultVal); err == nil {
			field.Set(reflect.ValueOf(val).Convert(field.Type()))
		}
	}

	return nil
}
