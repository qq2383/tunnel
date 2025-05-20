package config

import (
	"os"
	"reflect"

	"gopkg.in/yaml.v3"
)

var (
	filePath string
	cnf      any
)

func Add(c any) {
	cnf = c
}

func Get(field string) any {
	rv := reflect.ValueOf(cnf)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}

	if field != "" {
		rv = rv.FieldByName(field)
	}
	return rv.Addr().Interface()
}

func Getpath() string {
	return filePath
}

func Read(fp string) error {
	filePath = fp
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	dec := yaml.NewDecoder(f)
	err = dec.Decode(cnf)
	if err != nil {
		return err
	}
	return nil
}

func Setpath(fp string) {
	filePath = fp
}

func Write() error {
	buf, err := yaml.Marshal(&cnf)
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, buf, os.ModePerm)
}
