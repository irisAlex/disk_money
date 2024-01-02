// Author liyan
// Date 2020-03-27 5:24 下午
// Mail liyana@hualala.com
// org base cloud platform

package mongodb

import (
	"fmt"
	"reflect"
	"strings"
)

type collection interface {
	CollectionName() string
}

//type indexer interface {
//	Meta() []mgo.Index
//}

func colName(model interface{}) string {
	name := getColName(model)
	if name != "" {
		return name
	}
	//AppToken -> app_token
	tmp := fmt.Sprintf("%T", model)
	tmp = strings.Replace(tmp, "*", "", -1)
	tmp = strings.Replace(tmp, "]", "", -1)
	tmp = strings.Replace(tmp, "[", "", -1)
	ts := strings.Split(tmp, ".")
	if len(ts) < 2 {
		return tmp
	}
	return ts[1]
}

func getColName(model interface{}) string {
	s := reflect.TypeOf(model)
	for {
		if s.Kind() == reflect.Ptr || s.Kind() == reflect.Slice {
			s = s.Elem()
		} else {
			break
		}
	}

	v := reflect.New(s)
	vv := v.Interface()
	c, ok := vv.(collection)
	if !ok {
		return ""
	}
	return c.CollectionName()
}
