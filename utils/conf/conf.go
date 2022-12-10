package conf

import (
	"errors"
	"github.com/jinzhu/configor"
	"reflect"
)

func LoadConfig(path string, v interface{}) error {
	reflectValue := reflect.ValueOf(v)
	if reflectValue.Kind() != reflect.Ptr {
		return errors.New("resp not is prt")
	}
	if reflectValue.IsNil() {
		return errors.New("resp is nil prt")
	}

	err := configor.Load(v, path)
	if err != nil {
		return err
	}
	return nil
}
