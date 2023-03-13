package gocsv

import (
	"bufio"
	"errors"
	"io"
	"os"
	"reflect"
	"strings"
)

var (
	//TagName       = "csv"
	SkipHeaderNum = 2 // 表头需要跳过的行数
	KeysIndex     = 1 // 索引字符所在的行索引
)

type ReadRule struct {
	File  string
	Items interface{}
}

func OpenRead(list ...*ReadRule) error {
	for _, rule := range list {
		is := reflect.TypeOf(rule.Items)
		var rfItem reflect.Type
		switch is.Kind() {
		case reflect.Array, reflect.Slice:
			rfItem = is.Elem().Elem()
			break
		default:
			return errors.New("err: items is not in [array,slice]")
		}

		f, err := os.OpenFile(rule.File, os.O_RDWR|os.O_CREATE, os.ModePerm)
		if err != nil {
			return err
		}
		buf := bufio.NewReader(f)
		indexMap := make(map[string]int)
		var resultRe = reflect.ValueOf(rule.Items)
		i := -1
		for {
			i++
			b, _, err := buf.ReadLine()
			if err != nil {
				if err == io.EOF {
					break
				}
				return err
			}

			rows := strings.Split(string(b), ",")
			if i < SkipHeaderNum {
				if i == 1 {
					for j, key := range rows {
						indexMap[key] = j
					}
				}
				continue
			}

			it := reflect.New(rfItem)
			err = SetStructFieldVal(it.Interface(), indexMap, rows)
			if err != nil {
				println(err.Error())
			}
			resultRe = reflect.Append(resultRe, it)
		}
		reflect.ValueOf(&rule.Items).Elem().Set(resultRe)
	}
	return nil
}

type ReadBackRule struct {
	File     string
	Item     interface{}
	Callback func(v interface{})
}

func LoadCSVRowCallback(fileName string, itemObj interface{}, callback func(v interface{})) error {
	f, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	buf := bufio.NewReader(f)
	indexMap := make(map[string]int)
	var itemType = reflect.TypeOf(itemObj).Elem()
	i := -1
	for {
		i++
		b, _, err := buf.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		rows := strings.Split(string(b), ",")
		if i < SkipHeaderNum {
			if i == KeysIndex {
				for j, key := range rows {
					indexMap[key] = j
				}
			}
			continue
		}

		it := reflect.New(itemType)
		err = SetStructFieldVal(it.Interface(), indexMap, rows)
		if err != nil {
			println(err.Error())
		}
		callback(it.Interface())
	}
	return nil
}

func SetStructFieldVal(ptr interface{}, idx map[string]int, data []string) (err error) {
	v := reflect.ValueOf(ptr).Elem()
	for i := 0; i < v.NumField(); i++ {
		fieldInfo := v.Type().Field(i)
		kd := fieldInfo.Type.Kind()
		switch kd {
		case reflect.Struct:
			val := reflect.New(fieldInfo.Type)
			if err = SetStructFieldVal(val.Interface(), idx, data); err != nil {
				return
			}
			v.Field(i).Set(val.Elem())
			continue
		case reflect.Ptr:
			val := reflect.New(fieldInfo.Type.Elem())
			if err = SetStructFieldVal(val.Interface(), idx, data); err != nil {
				return
			}
			v.Field(i).Set(val)
			continue
		}

		name := fieldInfo.Tag.Get("csv")
		if name == "" || name == "-" {
			continue
		}

		if index, ok := idx[name]; ok {
			err = setField(v.Field(i), data[index], true)
			if err != nil {
				return err
			}
		}
	}
	return
}
