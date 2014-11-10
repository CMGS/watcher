package utils

import (
	"io/ioutil"
	"os"
	"reflect"
	"strconv"

	"../logs"
)

func WritePid(path string) {
	if err := ioutil.WriteFile(path, []byte(strconv.Itoa(os.Getpid())), 0755); err != nil {
		logs.Assert(err, "Save pid file failed")
	}
}

func MakeWrapper(fptr interface{}) {
	var maker = func(in []reflect.Value) []reflect.Value {
		wrapper := in[0].Elem()
		client := in[1]
		wrapperType := wrapper.Type()
		for i := 1; i < wrapperType.NumField(); i++ {
			field := wrapper.Field(i)
			f := client.MethodByName(wrapperType.Field(i).Name)
			field.Set(f)
		}
		return []reflect.Value{in[0]}
	}
	fn := reflect.ValueOf(fptr).Elem()
	v := reflect.MakeFunc(fn.Type(), maker)
	fn.Set(v)
}
